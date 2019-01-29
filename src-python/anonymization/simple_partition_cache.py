from typing import Dict
from uuid import UUID

from anonymization.partition import Partition


class SimplePartitionCache:
    def __init__(self):
        self.partitions: Dict[UUID, Partition] = {}

    def get_partition(self, partition_id: UUID):
        return self.partitions[partition_id]

    def add_partition(self, partition: Partition):
        self.partitions[partition.id] = partition

    def remove_partition(self, partition_id: UUID):
        del self.partitions[partition_id]
