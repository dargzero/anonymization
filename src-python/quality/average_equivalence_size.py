from typing import List

from anonymization.partition import Partition
from data.model import Model


def evaluate(k: int, total: int, results: List[Partition]):
    return total / (len(results) * k)
