import { readLines } from "https://deno.land/std@0.167.0/io/mod.ts";

type Monkey = {
  items: number[];
  op: (n: number) => number;
  check: (n: number) => boolean;
  yes: number;
  no: number;
};

// Helpers for converting input

function int(s: string) {
  return Number.parseInt(s, 10);
}

const c = {
  list<T>(value: string, sep: RegExp, mapping: (s: string) => T): T[] {
    return value.split(sep).map(mapping);
  },
  operation(value: string): (n: number) => number {
    const [x, op, y] = value.replace(/^new = /, "").split(/\s+/);

    const a = (n: number) => x === "old" ? n : int(x);
    const b = (n: number) => y === "old" ? n : int(y);

    if (op === "+") {
      return (n: number) => a(n) + b(n);
    } else if (op === "*") {
      return (n: number) => a(n) * b(n);
    }

    throw new Error(`Unrecognized operation: ${value}`);
  },
  check(value: string): (n: number) => boolean {
    const divisor = int(value.replace(/^divisible by /, ""));
    return (n: number) => n % divisor === 0;
  },
  branch(value: string): number {
    return int(value.replace(/^throw to monkey /, ""));
  },
};

// Read input

type Lines = AsyncIterableIterator<string>;

async function scan(lines: Lines, pattern: RegExp): Promise<string> {
  while (true) {
    const { value, done } = await lines.next();
    if (done) throw new Error(`Not found: ${pattern}`);
    if (pattern.test(value)) return value;
  }
}

async function from(lines: Lines): Promise<string[]> {
  await scan(lines, /Monkey \d+:/);

  const spec = [
    await scan(lines, /Starting items:/),
    await scan(lines, /Operation:/),
    await scan(lines, /Test:/),
    await scan(lines, /If true:/),
    await scan(lines, /If false:/),
  ].map((line) => line.trim().split(/:\s*/)[1]);

  return spec;
}

function build(spec: string[]): Monkey {
  const items = c.list(spec[0], /,\s*/, int);
  const op = c.operation(spec[1]);
  const check = c.check(spec[2]);
  const yes = c.branch(spec[3]);
  const no = c.branch(spec[4]);

  return { items, op, check, yes, no };
}

// Main

async function main() {
  const lines = readLines(Deno.stdin);
  const monkeys: Monkey[] = [];

  while (true) {
    try {
      const spec = await from(lines);
      monkeys.push(build(spec));
    } catch {
      break;
    }
  }

  const rounds = 20;
  const inspections = monkeys.map(() => 0);

  for (let i = 0; i < rounds; i++) {
    for (let j = 0; j < monkeys.length; j++) {
      const monkey = monkeys[j];

      while (monkey.items.length > 0) {
        const item = Math.floor(monkey.op(monkey.items.shift()!) / 3);
        const k = monkey.check(item) ? monkey.yes : monkey.no;
        const other = monkeys[k];
        other.items.push(item);

        inspections[j] += 1;
      }
    }
  }

  const sorted = inspections.slice().sort((a, b) => b - a);
  const business = sorted.slice(0, 2).reduce((x, t) => t * x);

  console.log(business);
}

if (import.meta.main) await main();
