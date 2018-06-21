function finalize(key, reducedValue) {
    return {
        min: reducedValue.min,
        max: reducedValue.max,
        avg: reducedValue.sum / reducedValue.count
    };
}