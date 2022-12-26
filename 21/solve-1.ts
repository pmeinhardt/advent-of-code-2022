import { readLines } from "https://deno.land/std@0.167.0/io/mod.ts";

const operators = Object.freeze(["+", "-", "*", "/"] as const);

type Operator = typeof operators[number];
type Value = number;
type Variable = string;

type Expression = Readonly<[Operator, Variable, Variable] | [Value]>;

type Rules = Record<Variable, Expression>;

type Stack = Array<Operator | Value | Variable>;

function isval(value: unknown): value is Value {
  return typeof value === "number";
}

function apply(op: Operator, a: Value, b: Value): Value {
  if (op === "+") return a + b;
  if (op === "-") return a - b;
  if (op === "*") return a * b;
  return a / b;
}

function evaluate(rules: Rules, root: Variable): Stack {
  const stack: Stack = [...rules[root]];

  for (let i = 0; i < stack.length; i++) {
    const symbol = stack[i];
    const expr = rules[symbol];
    if (expr) stack.splice(i, 1, ...expr);
  }

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

    const state = stack.join(" ");
    const rls = JSON.stringify(rules, null, 2);
    throw Error(`Evaluation failed: ${state}\nRules:\n${rls}`);
  }

  return stack;
}

function isop(value: string): value is Operator {
  return (operators as readonly string[]).includes(value);
}

function isnum(value: string): boolean {
  return /^\d+$/.test(value);
}

function num(value: string): number {
  return Number.parseInt(value, 10);
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

async function main() {
  const rules: Rules = {};

  for await (const line of readLines(Deno.stdin)) {
    const [key, value] = line.split(": ");
    rules[key] = parse(value);
  }

  console.log(...evaluate(rules, "root"));
}

if (import.meta.main) await main();
