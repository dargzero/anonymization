import time


class Profile:
    def __init__(self, name, callback):
        self.start = time.time()
        self.name = name
        self.callback = callback

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        end = time.time()
        runtime = end - self.start
        self.callback(runtime)
        print(f"{self.name} took {runtime} seconds to complete")
