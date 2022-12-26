import { readLines } from "https://deno.land/std@0.167.0/io/mod.ts";

function evaluate(rules: Record<string, string>, start: string) {
  let expr = start;

  while (true) {
    const match = expr.match(/[a-z]+/);
    if (match === null) break;
    const key = match[0];
    expr = expr.replace(key, `(${rules[key]})`);
  }

  return eval(expr);
}

async function main() {
  const rules: Record<string, string> = {};

  for await (const line of readLines(Deno.stdin)) {
    const [key, value] = line.split(": ");
    rules[key] = value;
  }

  console.log(evaluate(rules, rules.root));
}

if (import.meta.main) await main();
