"""
Reads and parses the Adult dataset (https://archive.ics.uci.edu/ml/datasets/adult) into a format that is consumable
for the anonymization algorithm
"""
from typing import Iterable, List, Optional

from data.model import Attribute, Model, ValueMappers

attributes: List[Attribute] = [
    Attribute("age",                                   is_qi=True),
    Attribute("workclass",       is_categorical=True,  is_qi=True),
    Attribute("final-weight"),
    Attribute("education",       is_categorical=True),
    Attribute("education-num",   is_categorical=True,  is_qi=True),
    Attribute("marital-status",  is_categorical=True,  is_qi=True),
    Attribute("occupation",      is_categorical=True,  is_qi=True),
    Attribute("relationship",    is_categorical=True),
    Attribute("race",            is_categorical=True,  is_qi=True),
    Attribute("sex",             is_categorical=True,  is_qi=True),
    Attribute("capital-gain"),
    Attribute("capital-loss"),
    Attribute("hours-per-week"),
    Attribute("native-country",  is_categorical=True,  is_qi=True, is_sensitive=True),
    Attribute("salary-class"),
]


def _try_parse(value):
    try:
        return float(value)
    except ValueError:
        return value


def _meaningful_lines(lines: Iterable[str]):
    for line in lines:
        line = line.strip()
        if line and "?" not in line:
            yield line


def _generate_record(line: str) -> List[str]:
    line = line.replace(" ", "")
    split_line = line.split(",")
    return [_try_parse(value) for value in split_line]


def read_split(file_path: str, mappers: Optional[ValueMappers], initial_batch: int):
    model = Model(attributes, mappers)
    with open(file_path, "rt") as file:
        records = [_generate_record(line) for line in _meaningful_lines(file)]
        model.add_records(records[0:initial_batch])
        return model, records[initial_batch:]


def read(file_path: str, mappers: Optional[ValueMappers], limit = None) -> Model:
    model = Model(attributes, mappers)
    with open(file_path, "rt") as file:
        records = [_generate_record(line) for line in _meaningful_lines(file)]
        if limit and limit != 0:
            records = records[:limit]
        model.add_records(records)

    return model
