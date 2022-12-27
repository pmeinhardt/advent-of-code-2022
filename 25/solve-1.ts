import { readLines } from "https://deno.land/std@0.167.0/io/mod.ts";

import * as snafu from "./snafu.ts";

async function main() {
  const values: number[] = [];

  for await (const line of readLines(Deno.stdin)) {
    values.push(snafu.decode(line));
  }

  const sum = values.reduce((acc, n) => acc + n, 0);

  console.log(snafu.encode(sum));
}

if (import.meta.main) await main();
