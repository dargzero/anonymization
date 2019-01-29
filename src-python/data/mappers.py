from typing import Dict

from data.model import Value


def numeric_occurrence(mapping: Dict[Value, Value], value: Value) -> Value:
    return len(mapping)
