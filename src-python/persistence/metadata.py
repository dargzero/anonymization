from anonymization.mondrian import Mondrian


class Metadata:
    def __init__(self, mondrian: Mondrian):
        self.k = mondrian.k
        self.l = mondrian.l
        self.qi_ranges = mondrian.qi_ranges
        self.attributes = mondrian.attributes
        self.mapping = mondrian.transform.mapping

    def to_serializable(self):
        return {
            "k": self.k,
            "l": self.l,
            "qi_ranges": [{
                "dimension": qi_range[0],
                "range": qi_range[1].to_serializable()
            } for qi_range in self.qi_ranges.items()],
            "attributes": [att.to_serializable()
                           for att in self.attributes],
            "mapping": [{
                "dimension": mapping[0],
                "map": [{
                    "from": m[0],
                    "to": m[1]
                } for m in mapping[1].map.items()],
            } for mapping in self.mapping.items()]
        }
