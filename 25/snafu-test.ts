// Run with `deno test snafu-test.ts`

import { assertStrictEquals } from "https://deno.land/std@0.171.0/testing/asserts.ts";

import { decode, encode } from "./snafu.ts";

const cases: [number, string][] = [
  [1, "1"],
  [2, "2"],
  [3, "1="],
  [4, "1-"],
  [5, "10"],
  [6, "11"],
  [7, "12"],
  [8, "2="],
  [9, "2-"],
  [10, "20"],
  [11, "21"],
  [15, "1=0"],
  [20, "1-0"],
  [31, "111"],
  [32, "112"],
  [37, "122"],
  [107, "1-12"],
  [198, "2=0="],
  [201, "2=01"],
  [353, "1=-1="],
  [906, "12111"],
  [1257, "20012"],
  [1747, "1=-0-2"],
  [2022, "1=11-2"],
  [12345, "1-0---0"],
  [314159265, "1121-1110-1=0"],
];

cases.forEach(([decimal, snafu]) => {
  Deno.test(`decode ${snafu} to ${decimal}`, () => {
    assertStrictEquals(decode(snafu), decimal);
  });
});

cases.forEach(([decimal, snafu]) => {
  Deno.test(`encode ${decimal} to ${snafu}`, () => {
    assertStrictEquals(encode(decimal), snafu);
  });
});
