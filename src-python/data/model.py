from dataclasses import dataclass
from typing import Dict, Callable, List, Iterable, Optional, Union


Value = Union[str, float, int]
ValueMapper = Callable[[Dict[Value, Value], Value], Value]
ValueMappers = Dict[str, ValueMapper]
Record = List[Value]
Records = List[Record]


@dataclass
class Attribute:
    name: str
    is_qi: bool = False
    is_categorical: bool = False
    is_sensitive: bool = False

    def to_serializable(self):
        return self.__dict__


class Mapping:
    def __init__(self, mapper: ValueMapper):
        self.mapper = mapper
        self.map: Dict[Value, Value] = {}
        self.reverse_map: Dict[Value, Value] = {}


class ValueTransform:
    def __init__(self, mappers: ValueMappers):
        self.mapping: Dict[str, Mapping] = {}
        for attribute_name, mapper in mappers.items():
            self.mapping[attribute_name] = Mapping(mapper)

    def get_mapped(self, attribute: Attribute, value: Value) -> Value:
        if attribute.name not in self.mapping:
            return value

        attribute_map = self.mapping[attribute.name]
        if value in attribute_map.map:
            return attribute_map.map[value]

        # If the value is not found, we need to run the mapper
        mapped_value = attribute_map.mapper(attribute_map.map, value)
        attribute_map.map[value] = mapped_value
        attribute_map.reverse_map[mapped_value] = value
        return mapped_value

    def get_unmapped(self, attribute: Attribute, value: Value) -> Value:
        if attribute.name not in self.mapping:
            return value
        return self.mapping[attribute.name].reverse_map[value]


class IncrementalUpdate:
    def __init__(self, attributes: List[Attribute], records: Records):
        self.attributes = attributes
        self.records = records


class Model:
    def __init__(self, attributes: List[Attribute], mappers: Optional[ValueMappers]):
        self.records = []
        self.attributes = attributes
        self.transform = ValueTransform(mappers) if mappers else None

    def add_record(self, record: list):
        self.records.append(self._transform_record(record))

    def add_records(self, records: Iterable[Record]):
        self.records.extend((self._transform_record(record) for record in records))

    def produce_incremental_update(self, records: Iterable[Record]) -> IncrementalUpdate:
        transformed_records = [self._transform_record(record) for record in records]
        self.records.extend(transformed_records)
        return IncrementalUpdate(self.attributes, transformed_records)

    def original_records(self):
        for record in self.records:
            yield self._reverse_transform_record(record)

    def _transform_record(self, record: Record):
        if not self.transform:
            return record
        return [self.transform.get_mapped(self.attributes[index], value) for index, value in enumerate(record)]

    def _reverse_transform_record(self, record: Record):
        if not self.transform:
            return record
        return [self.transform.get_unmapped(self.attributes[index], value) for index, value in enumerate(record)]
