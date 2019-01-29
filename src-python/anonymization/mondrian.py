import uuid
from dataclasses import dataclass
from typing import List, Optional

from anonymization.simple_partition_cache import SimplePartitionCache
from data.model import IncrementalUpdate, Model, Record, Records, Attribute
from anonymization.partition import Partition
from anonymization.qi_range import QiRange, QiRanges
from scheduler import thread


def _generate_qi_ranges(records: Records, qi_attributes: List[int]) -> QiRanges:
    ranges: QiRanges = {}
    for qi in qi_attributes:
        values = sorted(list(set((record[qi] for record in records))))
        qi_range = values[-1] - values[0]
        ranges[qi] = QiRange(qi_range, values)
    return ranges


def _update_qi_ranges(qi_ranges: QiRanges, records: Records):
    for qi, r in qi_ranges.items():
        values = sorted(list(set(record[qi] for record in records).union(r.values)))
        qi_range = values[-1] - values[0]
        qi_ranges[qi] = QiRange(qi_range, values)


class MondrianTree:
    def __init__(self):
        self.id = uuid.uuid4()
        self.dimension: Optional[int] = None
        self.value: Optional[float, int] = None
        self.children: List['MondrianTree'] = []
        self.partition_id: Optional[uuid.UUID] = None

    def collect_results(self, mondrian: 'Mondrian'):
        def collect_results_recursive(tree: 'MondrianTree', r: List[Partition]):
            if tree.partition_id:
                results.append(mondrian.partition_cache.get_partition(tree.partition_id))
            for child in tree.children:
                collect_results_recursive(child, r)

        results = []
        collect_results_recursive(self, results)
        return results


@dataclass
class TreeBuilder:
    tree: MondrianTree
    partition: Partition


class Mondrian:
    root: Optional[MondrianTree]

    def __init__(self, attributes: List[Attribute], transform, qi_attributes: List[int], qi_ranges: QiRanges,
                 partition_cache=SimplePartitionCache(), k: Optional[int]=None, l: Optional[int]=None):
        if l:
            self.k = max(k or 0, l)
            self.l = l
        else:
            self.k = k
            self.l = None
        self.attributes = attributes
        self.transform = transform
        self.qi_attributes = qi_attributes
        self.qi_ranges = qi_ranges
        self.partition_cache = partition_cache

    @staticmethod
    def from_model(model: Model, k: Optional[int]=None, l: Optional[int]=None):
        qi_attributes = [i for i, attribute in enumerate(model.attributes) if attribute.is_qi]
        qi_ranges = _generate_qi_ranges(model.records, qi_attributes)
        return Mondrian(model.attributes, model.transform, qi_attributes, qi_ranges, k=k, l=l)

    def collect_results(self):
        return self.root.collect_results(self)

    def anonymize(self, records: Records):
        self.root = MondrianTree()
        builders = [TreeBuilder(self.root, self._make_initial_partition(records))]
        results = self._anonymize_iterative(builders)
        for result in results:
            self.partition_cache.add_partition(result)
        return [], results

    def anonymize_incremental(self, update: IncrementalUpdate):
        if not self.root:
            raise ValueError("Cannot run incremental updates against a non-anonymized dataset")
        _update_qi_ranges(self.qi_ranges, update.records)
        builders, removed_partitions = self._distribute_updates(update)

        removed_trees = [b.tree.id for b in builders]
        results = self._anonymize_iterative(builders)
        for result in results:
            self.partition_cache.add_partition(result)
        return removed_trees, removed_partitions, results

    def _distribute_updates(self, update: IncrementalUpdate):
        modified_builders = {}
        removed_partitions = []
        for record in update.records:
            builder = self._add_record_to_partition(record)
            modified_builders[builder.partition] = builder

        size_threshold = 2 * self.k
        builders = list(builder for builder
                        in modified_builders.values()
                        if len(builder.partition) >= size_threshold)
        for builder in builders:
            partition_id = builder.tree.partition_id
            if partition_id:
                removed_partitions.append(partition_id)
                self.partition_cache.remove_partition(partition_id)
            builder.partition.disallowed_cuts = set()
            builder.tree.partition_id = None
        return builders, removed_partitions

    def _add_record_to_partition(self, record: Record):
        node = self.root
        while not node.partition_id:
            value = record[node.dimension]
            if value <= node.value:
                node = node.children[0]
            else:
                node = node.children[1]

        node_partition = self.partition_cache.get_partition(node.partition_id)
        node_partition.records.append(record)
        return TreeBuilder(node, node_partition)

    def _anonymize_iterative(self, builders: List[TreeBuilder]):
        results = []
        wavefront = builders

        while wavefront:
            current = wavefront.pop()
            next_partitions = self._anonymize_partition(current, results)
            wavefront.extend(next_partitions)

        return results

    def _anonymize_partition(self, builder: TreeBuilder, results: List[Partition]):
        partition = builder.partition

        # If a partition has no more allowable cuts
        # then we know it is done -> add it to the result set
        new_partitions = []
        if len(partition.disallowed_cuts) == len(self.qi_attributes):
            results.append(partition)
            builder.tree.partition_id = partition.id
            return new_partitions

        # Otherwise choose a dimension and calculate a median
        dimension = partition.choose_dimension(self.qi_ranges)
        attribute = self.attributes[dimension]
        median = partition.find_median(dimension,
                                       attribute.is_categorical)

        # And split the original partition into two
        low = partition.make_low(dimension, median)
        high = partition.make_high(dimension, median)

        if not self._satisfies_anonymity_constraints(low) or \
           not self._satisfies_anonymity_constraints(high):
            partition.disallowed_cuts.add(dimension)
            return [builder]

        low_builder = TreeBuilder(MondrianTree(), low)
        high_builder = TreeBuilder(MondrianTree(), high)

        builder.tree.dimension = dimension
        builder.tree.value = median
        builder.tree.children.append(low_builder.tree)
        builder.tree.children.append(high_builder.tree)

        new_partitions.append(low_builder)
        new_partitions.append(high_builder)
        return new_partitions

    def _make_initial_partition(self, records: Records):
        return Partition(records, self.qi_attributes, {})

    def _satisfies_anonymity_constraints(self, partition: Partition):
        return self._is_k_anonym(partition) and self._is_l_diverse(partition)

    def _is_k_anonym(self, partition: Partition):
        return len(partition) >= self.k

    def _is_l_diverse(self, partition: Partition):
        if not self.l:
            return True
        sensitive_attribute_values = {i: set() for i, att
                                      in enumerate(self.attributes)
                                      if att.is_sensitive}

        for record in partition.records:
            for index, values in sensitive_attribute_values.items():
                values.add(record[index])

        for values in sensitive_attribute_values.values():
            if len(values) < self.l:
                    return False
        return True
