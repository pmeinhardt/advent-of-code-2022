const d2n: Readonly<Record<string, number>> = Object.freeze({
  "=": -2,
  "-": -1,
  "0":  0,
  "1":  1,
  "2":  2,
});

export function decode(value: string): number {
  let result = 0;

  const len = value.length;

  for (let i = 0; i < len; i++) {
    const digit = value[len - 1 - i];

    if (digit in d2n) {
      result += d2n[digit] * Math.pow(5, i);
    } else {
      throw new Error(`Invalid: ${value}`);
    }
  }

  return result;
}

const n2d = Object.freeze(["0", "1", "2", "=", "-"] as const);

function divmod(n: number, d: number): [number, number] {
  const mod = n % d;
  const div = (n - mod) / d;
  return [div, mod];
}

export function encode(value: number): string {
  let result = "";

  while (value !== 0) {
    const [div, mod] = divmod(value, 5);

    const digit = n2d[mod < 0 ? mod + 5 : mod];
    result = digit + result;

    const rem = Math.abs(mod) > 2 ? Math.sign(mod) : 0;
    value = div + rem;
  }

  if (result.length === 0 && value === 0) result += "0";

  return result;
}
