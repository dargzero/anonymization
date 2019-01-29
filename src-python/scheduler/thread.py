from concurrent.futures import ThreadPoolExecutor


def run_tasks(tasks):
    results = []
    with ThreadPoolExecutor(max_workers=4) as executor:
        futures = [executor.submit(task) for task in tasks]
        for future in futures:
            result = future.result()
            results.extend(result)
    return results
