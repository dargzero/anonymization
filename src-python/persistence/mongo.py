import functools
from typing import List, Dict, Any, Set
from uuid import UUID

import pymongo
from pymongo import UpdateOne, InsertOne

from anonymization.mondrian import Mondrian, MondrianTree
from anonymization.partition import Partition, Range
from data.model import Records, IncrementalUpdate
from persistence.metadata import Metadata


def named_record(name, record):
    return {"_id": name, "value": record}


def diff_tree(tree: MondrianTree, cut: Set[str],
              added_nodes: List[MondrianTree]):
    if tree.id in cut:
        added_nodes.append(tree)
        return
    for child in tree.children:
        diff_tree(child, cut, added_nodes)


def serialize_partition(partition: Partition):
    return {
        "id": partition.id,
        "records": partition.records,
        "qi_indices": partition.qi_indices,
        "range": [{"dimension": qi, "range": qi_range.to_serializable()} for qi, qi_range in partition.range.items()]
    }


def deserialize_partition(partition: Dict[str, Any]):
    partition = Partition(partition["records"],
                          partition["qi_indices"],
                          {r["dimension"]: Range(r["range"]["low"], r["range"]["high"]) for r in partition["range"]},
                          partition_id=partition["id"])
    return partition


def serialize_tree_node(tree: MondrianTree, partitions: Dict[UUID, Partition]):
    return {
        "_id": tree.id,
        "dimension": tree.dimension,
        "partition": serialize_partition(partitions[tree.partition_id]) if tree.partition_id else None,
        "children": [child.id for child in tree.children],
    }


def serialize_tree(tree: MondrianTree, partitions: Dict[UUID, Partition]):
    nodes = []
    queue = [tree]
    while queue:
        current = queue.pop()
        nodes.append(serialize_tree_node(current, partitions))
        queue.extend(current.children)
    return nodes


def tree_children(node: MondrianTree):
    return {"children": [child.id for child in node.children]}


def produce_updates(nodes: List[MondrianTree], is_incremental, parts):
    operations = []
    if is_incremental:
        operations.extend(UpdateOne({"_id": node.id},
                                    {"$unset": {"partition": {}},
                                     "$set": tree_children(node)})
                          for node in nodes)
        nodes = [child for node in nodes for child in node.children]

    operations.extend(InsertOne(serialized) for node in nodes
                      for serialized in serialize_tree(node, parts))

    return operations


def update_tree(collection: pymongo.collection.Collection,
                tree: MondrianTree, previous_leaves: Set[str],
                partitions: Dict[UUID, Partition]):
    new_nodes = []
    is_incremental = True
    if not previous_leaves:
        new_nodes = [tree]
        is_incremental = False
    else:
        diff_tree(tree, previous_leaves, new_nodes)
    collection.bulk_write(produce_updates(new_nodes, is_incremental,
                                          partitions), ordered=False, bypass_document_validation=True)


class MongoPartitionCache:
    def __init__(self, db):
        self.db = db

    @functools.lru_cache(maxsize=128)
    def get_partition(self, partition_id: UUID):
        trees: pymongo.collection.Collection = self.db.trees
        tree_with_partition = trees.find_one({"partition.id": partition_id})
        return deserialize_partition(tree_with_partition["partition"])

    def add_partition(self, partition: Partition):
        pass

    def remove_partition(self, partition_id: UUID):
        pass


class MongoPersistedMondrian:
    def __init__(self, name: str, mondrian: Mondrian):
        self.name = name
        self.mondrian = mondrian
        self.client = pymongo.MongoClient("localhost", 27017)
        self.mondrian_db = self.client.mondrian

        # self.mondrian.partition_cache = MongoPartitionCache(self.mondrian_db)
        self._initialize()

    def anonymize(self, records: Records):
        removed, added = self.mondrian.anonymize(records)
        update_tree(self.mondrian_db.trees, self.mondrian.root, removed, {part.id: part for part in added})
        return added

    def anonymize_incremental(self, update: IncrementalUpdate):
        removed_trees, removed_partitions, added = self.mondrian.anonymize_incremental(update)
        removed_trees_set = set(removed_trees)
        update_tree(self.mondrian_db.trees, self.mondrian.root, removed_trees_set, {part.id: part for part in added})
        return added

    def _initialize(self):
        self._write_metadata()

    def _write_metadata(self):
        metadata = Metadata(self.mondrian)
        self.mondrian_db.metadata.update_one({"_id": self.name},
                                             {"$set": named_record(self.name, metadata.to_serializable())},
                                             upsert=True)
