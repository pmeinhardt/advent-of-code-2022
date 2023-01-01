# Advent of Code 2022 🎄

Solutions to puzzles are organized into directories, one directory for each day — `01`, `02`, `03`, …, `24`, `25`:

```
01:
  solve-1           # implementation for part 1, expects puzzle input on stdin
  solve-2           # implementation for part 2, expects puzzle input on stdin
  test-input.txt    # example input
  input.txt         # puzzle input (.gitignore-d)
```

To invoke puzzle implementations from the project root, you can run:

```text
./solve <day> [part [--test|-t]]
```

For instance, `./solve 5` and `./solve 5 1` will invoke day 5 part 1.
Running `./solve 5 2` will invoke part 2.

Adding `--test` or `-t` after the part number, will use input from `test-input.txt` instead of `input.txt`.
