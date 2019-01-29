from typing import List

from anonymization.mondrian import Mondrian
from anonymization.partition import Partition
from data.model import Model


def evaluate(total: int, mondrian: Mondrian, results: List[Partition]):
    gcp = 0.0
    for partition in results:
        ncp = len(partition) * sum((partition.normalized_range(qi, mondrian.qi_ranges) for qi in mondrian.qi_attributes))
        gcp += ncp

    gcp /= (total * len(mondrian.qi_attributes))
    return gcp
