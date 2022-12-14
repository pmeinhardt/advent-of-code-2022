def cmp(a, b):
    if type(a) == type(b) == int:
        return a - b

    if type(a) == type(b) == list:
        for x, y in zip(a, b):
            rel = cmp(x, y)
            if rel != 0:
                return rel
        return len(a) - len(b)

    if type(a) == list:
        return cmp(a, [b])

    if type(b) == list:
        return cmp([a], b)

    raise "unknown types"

def read(file):
    while True:
        try:
            left, right = map(eval, [next(file), next(file)])
            yield left, right
            next(file)
        except StopIteration:
            break
