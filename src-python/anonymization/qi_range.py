from typing import Dict, List, Union
from dataclasses import dataclass


from data.model import Value


@dataclass
class QiRange:
    range: Union[float, int]
    values: List[Value]

    def to_serializable(self):
        return self.__dict__


QiRanges = Dict[int, QiRange]
