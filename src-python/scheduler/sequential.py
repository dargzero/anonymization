def run_tasks(tasks):
    results = []
    for task in tasks:
        results.extend(task())
    return results
