function map() {{
    const getValue = number => {{
        return {{
            min: number,
            max: number,
            sum: number,
            count: 1
        }};
    }};

    const number = Number.parseFloat(this.{0});
    if (!Number.isNaN(number)) {{
        emit(null, getValue(0));
        return;
    }}

    const regex = /^(\[|\])(\d+\.?\d*), (\d+\.?\d*)(\[|\])$/;
    const match = regex.exec(this.{0});
    if (match == null)
        return;
    
    const min = Number.parseFloat(match[2])
    const max = Number.parseFloat(match[3])
    if (!Number.isNaN(min) && !Number.isNaN(max))
        emit(null, getValue(max - min))
}}