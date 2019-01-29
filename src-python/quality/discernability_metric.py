from typing import List

from anonymization.partition import Partition


def evaluate(results: List[Partition]):
    return sum((len(partition)**2 for partition in results))
