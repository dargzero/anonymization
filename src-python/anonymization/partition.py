import uuid

from dataclasses import dataclass
from typing import Dict, List, Optional, Union
from statistics import median, median_high

from data.model import Value
from anonymization.qi_range import QiRanges


@dataclass
class Range:
    low: Union[float, int]
    high: Union[float, int]

    def to_serializable(self):
        return {"low": self.low, "high": self.high}


class Partition:
    def __init__(self, records: List[List[Value]], qi_indices: List[int], partition_range: Dict[int, Range],
                 partition_id=None):
        self.id = partition_id or uuid.uuid4()
        self.records = records
        self.qi_indices = qi_indices
        self.range = dict(partition_range)
        self.disallowed_cuts = set()

    def make_low(self, qi: int, upper: float):
        records = [record for record in self.records if record[qi] <= upper]
        new_range = dict(self.range)
        new_range[qi] = Range(min((record[qi] for record in self.records)), upper)
        return Partition(records, self.qi_indices, new_range)

    def make_high(self, qi: int, lower: float):
        records = [record for record in self.records if record[qi] > lower]
        new_range = dict(self.range)
        new_range[qi] = Range(lower, max((record[qi] for record in self.records)))
        return Partition(records, self.qi_indices, new_range)

    def add_record(self, record: List[Value]):
        self.records.append(record)

    def transfer_record(self, other: 'Partition', index: int):
        record = self.records.pop(index)
        other.records.append(record)

    def find_median(self, qi: int, is_categorical: bool) -> Optional[float]:
        if is_categorical:
            return median_high((record[qi] for record in self.records))
        else:
            return median((record[qi] for record in self.records))

    def choose_dimension(self, ranges: QiRanges) -> int:
        if len(self.range) == 0:
            max_range = max(ranges.items(),
                            key=lambda dimension: dimension[1].range)
            return max_range[0]
        return max((qi for qi
                    in self.qi_indices
                    if qi not in self.disallowed_cuts),
                   key=lambda qi: self.normalized_range(qi, ranges))

    def normalized_range(self, qi: int, ranges: QiRanges) -> float:
        qi_range = ranges[qi]
        try:
            partition_range = self.range[qi]
        except KeyError:
            partition_range = self._minmax_dimension(qi)
        return (partition_range.high - partition_range.low) / qi_range.range

    def _minmax_dimension(self, qi: int):
        min_element = self.records[0][qi]
        max_element = self.records[0][qi]
        for i in range(1, len(self.records)):
            value = self.records[i][qi]
            if value < min_element:
                min_element = value
            if value > max_element:
                max_element = value
        return Range(min_element, max_element)

    def __len__(self):
        return len(self.records)
