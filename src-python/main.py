import data.adult as adult
import data.mappers as mappers
import quality.gcp as gcp
import quality.average_equivalence_size as average_equivalence_size
import quality.discernability_metric as discernability_metric

from anonymization.mondrian import Mondrian
from persistence.mongo import MongoPersistedMondrian
from persistence.file import FilePersistedMondrian
from util.profile import Profile

mappers = {attribute.name: mappers.numeric_occurrence for attribute in adult.attributes if attribute.is_categorical}


def print_table(values):
    for k, v in values.items():
        print(f"{k}\t{v[0]}\t{v[1]}")


def metric(fn, limit: int):
    model = adult.read("resources/adult.data", mappers=mappers, limit=limit)
    mondrian = Mondrian.from_model(model, k=10)
    mondrian.anonymize(model.records)
    results = mondrian.collect_results()
    return fn(mondrian, limit, results)


def metric_inc(fn, mondrian: Mondrian, records, total: int):
    mondrian.anonymize_incremental(records)
    results = mondrian.collect_results()
    return fn(mondrian, total, results)


def generate_metric(metric_name, metric_fn):
    print(f"----------{metric_name}----------")
    batch_size = 2000
    records_size = 30000
    records = {}

    print("Full:")
    for i in range(batch_size, records_size + batch_size, batch_size):
        metric_i = metric(metric_fn, i)
        records[i] = [metric_i]
        print(f"{i}\t{metric_i}")

    print("\nIncremental:")
    model, rest = adult.read_split("resources/adult.data", mappers=mappers, initial_batch=batch_size)
    mondrian = Mondrian.from_model(model, k=10)
    mondrian.anonymize(model.records)
    initial_batch_metric = metric_fn(mondrian, batch_size, mondrian.collect_results())
    print(f"{batch_size}\t{initial_batch_metric}")
    records[batch_size].append(initial_batch_metric)

    for i in range(0, records_size - batch_size, batch_size):
        total = i + batch_size*2
        metric_i = metric_inc(metric_fn, mondrian, model.produce_incremental_update(rest[i:i+batch_size]), total)
        records[total].append(metric_i)
        print(f"{total}\t{metric_i}")

    diffs = {k: (max(v[0], v[1]) - min(v[0], v[1])) / max(v[0], v[1]) for k, v in records.items()}

    print("\nResults:")
    print_table(records)
    average_diff = sum(diffs.values()) / len(diffs)
    print(f"Average diff: {average_diff }")


def measure_performance():
    batch_size = 2000
    records_size = 30000
    records = {}

    def update_record(i):
        def set_time(time):
            if i in records:
                records[i].append(time)
            else:
                records[i] = [time]
        return set_time

    print("\nRegular:")
    model, rest = adult.read_split("resources/adult.data", mappers=mappers, initial_batch=batch_size)
    mondrian = FilePersistedMondrian("adult", Mondrian.from_model(model, k=10))
    update_record(batch_size)(0)
    with Profile(f"{batch_size}", update_record(batch_size)):
        mondrian.anonymize(model.records)

    for i in range(0, records_size - batch_size, batch_size):
        total = i + batch_size*2
        update = model.produce_incremental_update(rest[i:i+batch_size])
        update_record(total)(0)
        with Profile(f"{total}", update_record(total)):
            mondrian.anonymize_incremental(update)

    print("\nPersisted:")
    model, rest = adult.read_split("resources/adult.data", mappers=mappers, initial_batch=batch_size)
    mondrian = MongoPersistedMondrian("adult", Mondrian.from_model(model, k=10))
    with Profile(f"{batch_size}", update_record(batch_size)):
        mondrian.anonymize(model.records)

    for i in range(0, records_size - batch_size, batch_size):
        total = i + batch_size*2
        update = model.produce_incremental_update(rest[i:i+batch_size])
        with Profile(f"{total}", update_record(total)):
            mondrian.anonymize_incremental(update)

    sums = {}
    previous_key = None
    for key, value in records.items():
        previous = [0, 0] if not previous_key else sums[previous_key]
        sums[key] = [value[0] + previous[0], value[1] + previous[1]]
        previous_key = key

    print("\nResults:")
    print_table(records)

    print("\nSums:")
    print_table(sums)


if __name__ == "__main__":
    generate_metric("Discernability penalty", lambda _, __, results: discernability_metric.evaluate(results))
    generate_metric("Average equivalence class size", lambda mondrian, total, results: average_equivalence_size.evaluate(mondrian.k, total, results))
    generate_metric("GCP", lambda mondrian, total, results: gcp.evaluate(total, mondrian, results))
    measure_performance()
