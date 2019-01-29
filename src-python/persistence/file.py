import json
import os

from anonymization.mondrian import Mondrian
from data.model import Records, IncrementalUpdate
from persistence.metadata import Metadata
from persistence.util import json_default


class FilePersistedMondrian:
    def __init__(self, name: str, mondrian: Mondrian):
        self.name = name
        self.mondrian = mondrian
        self._initialize()

    def anonymize(self, records: Records):
        removed, added = self.mondrian.anonymize(records)
        self._persist_results(self.mondrian.root)
        return added

    def anonymize_incremental(self, update: IncrementalUpdate):
        removed_trees, removed_partitions, added = self.mondrian.anonymize_incremental(update)
        self._persist_results(self.mondrian.root)
        return added

    def _persist_results(self, tree):
        with self._mondrian_tree("w") as tree_storage:
            tree = json.dumps(tree, default=json_default)
            tree_storage.write(tree)

    def _initialize(self):
        if not os.path.exists(self.name):
            os.mkdir(self.name)
        self._write_metadata()

    def _write_metadata(self):
        with self._mondrian_metadata("w") as metadata_storage:
            metadata = json.dumps(Metadata(self.mondrian), default=json_default)
            metadata_storage.write(metadata)


    def _mondrian_file(self, mode: str):
        return open(os.path.join(self.name, "data.mo"), mode)

    def _mondrian_metadata(self, mode: str):
        return open(os.path.join(self.name, "metadata.mo"), mode)

    def _mondrian_tree(self, mode: str):
        return open(os.path.join(self.name, "tree.mo"), mode)
