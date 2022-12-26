import { readLines } from "https://deno.land/std@0.167.0/io/mod.ts";

const operators = Object.freeze(["+", "-", "*", "/"] as const);

type Operator = typeof operators[number];
type Value = number;
type Variable = string;

type Expression = Readonly<[Operator, Variable, Variable] | [Value]>;

type Rules = Record<Variable, Expression>;

type Stack = Array<Operator | Value | Variable>;

function isvar(value: unknown): value is Variable {
  return typeof value === "string" && !isop(value);
}

function isval(value: unknown): value is Value {
  return typeof value === "number";
}

function isop(value: unknown): value is Operator {
  return (operators as readonly unknown[]).includes(value);
}

function isnum(value: string): boolean {
  return /^\d+$/.test(value);
}

function num(value: string): number {
  return Number.parseInt(value, 10);
}

function expand(rules: Rules, expr: Expression): Stack {
  const stack: Stack = [...expr];

  for (let i = 0; i < stack.length; i++) {
    const symbol = stack[i];
    const expr = rules[symbol];
    if (expr) stack.splice(i, 1, ...expr);
  }

  return stack;
}

function apply(op: Operator, a: Value, b: Value): Value {
  if (op === "+") return a + b;
  if (op === "-") return a - b;
  if (op === "*") return a * b;
  return a / b;
}

function evaluate(stack: Stack): Stack {
  for (let i = stack.length - 1; i >= 0; i--) {
    const x = stack[i];

    if (isval(x)) {
      continue;
    }

    const a = stack[i + 1];
    const b = stack[i + 2];

    if (isop(x) && isval(a) && isval(b)) {
      stack.splice(i, 3, apply(x, a, b));
      continue;
    }
  }

  return stack;
}

function parse(value: string): Expression {
  if (isnum(value)) {
    return [num(value)];
  }

  const match = value.match(/^(?<a>[a-z]+) (?<op>[-+/*]) (?<b>[a-z]+)$/);

  if (match && match.groups && isop(match.groups.op)) {
    const { a, b, op } = match.groups;
    return [op, a, b];
  }

  throw new Error(`Invalid expression: ${value}`);
}

function getexpr(stack: Stack, index: number): Stack {
  let lv = 1;
  let i = index;

  while (lv > 0) {
    const e = stack[i];

    if (isop(e)) {
      lv++;
    } else {
      lv--;
    }

    i += 1;
  }

  return stack.slice(index, i);
}

function isdone(stack: Stack): boolean {
  return stack.length === 1;
}

function isdet(stack: Stack): stack is [Value] {
  return isdone(stack) && isval(stack[0]);
}

const inverse: Record<Operator, Operator> = {
  "+": "-",
  "-": "+",
  "*": "/",
  "/": "*",
};

function solve(rules: Rules, expr: Expression, value: number) {
  let result = value;
  let stack = expand(rules, expr);

  while (!isdone(stack)) {
    const op = stack[0] as Operator;

    let l = getexpr(stack, 1);
    let r = getexpr(stack, 1 + l.length);

    l = evaluate(l);
    r = evaluate(r);

    if (isdet(l)) {
      if (op === "/") { // (result = l / r) <==> (r = l / result)
        result = apply("/", l[0], result);
      } else if (op === "-") { // (result = l - r) <==> (r = l - result)
        result = apply("-", l[0], result);
      } else {
        result = apply(inverse[op], result, l[0]);
      }
      stack = r;
    } else if (isdet(r)) {
      result = apply(inverse[op], result, r[0]);
      stack = l;
    } else {
      throw Error(`Variables in more than one branch:\n${l}\n${r}`);
    }
  }

  if (isvar(stack[0])) return result;

  throw new Error(`Invalid resolution: ${result} = ${stack}`);
}

async function main() {
  const rules: Rules = {};

  for await (const line of readLines(Deno.stdin)) {
    const [key, value] = line.split(": ");
    rules[key] = parse(value);
  }

  delete rules["humn"];

  const [, left, right] = rules.root as readonly [Operator, Variable, Variable];

  console.log(solve(rules, ["-", left, right], 0));
}

if (import.meta.main) await main();
