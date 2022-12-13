def lt(a, b):
    if type(a) == type(b) == int:
        return a < b

    if type(a) == type(b) == list:
        for x, y in zip(a, b):
            if x == y:
                continue
            return lt(x, y)
        return len(a) < len(b)

    if type(a) == list:
        return lt(a, [b])

    if type(b) == list:
        return lt([a], b)

    raise "unknown types"

def read(file):
    while True:
        try:
            left, right = map(eval, [next(file), next(file)])
            yield left, right
            next(file)
        except StopIteration:
            break
