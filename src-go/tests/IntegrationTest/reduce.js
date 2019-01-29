function reduce(key, values) {
    const sumReducer = (accumulator, currentValue) => accumulator + currentValue;

    return {
        min: Math.min(...values.map(value => value.min)),
        max: Math.max(...values.map(value => value.max)),
        sum: values.map(value => value.sum).reduce(sumReducer),
        count: values.map(value => value.count).reduce(sumReducer)
    };
}