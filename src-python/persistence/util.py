from uuid import UUID


def json_default(o):
    if isinstance(o, set):
        return list(o)
    elif isinstance(o, UUID):
        return o.hex
    elif callable(o):
        return None
    return o.__dict__

