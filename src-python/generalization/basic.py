from typing import List

from anonymization.partition import Partition
from data.model import ValueTransform, Records, Attribute


def generalize(partitions: List[Partition], attributes: List[Attribute], transform: ValueTransform) -> Records:
    result = []
    for partition in partitions:
        generalized_values = {}
        for qi, qi_range in partition.range.items():
            attribute = attributes[qi]
            low_unmapped = transform.get_unmapped(attribute, qi_range.low)
            high_unmapped = transform.get_unmapped(attribute, qi_range.high)
            if low_unmapped == high_unmapped:
                generalized_values[qi] = str(low_unmapped)
            else:
                generalized_values[qi] = f"{low_unmapped}-{high_unmapped}"

        for record in partition.records:
            result.append([generalized_values[i] if i in generalized_values else v for i, v in enumerate(record)])

    return result
