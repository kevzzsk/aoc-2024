# [Advent of Code 2024](https://adventofcode.com/2024)

This repo contains my attempt at Advent of Code 2024 with Golang. Inside there will be code solutions separated by the advent days folder.


## Learnings

- Day 2 p2
    - `append()` slice function modifies the original slice
- Day 3
    - Golang regex are not able to do lookaheads `(?=)` or lookbehind `(?<=)` due to it using [RE2](https://stackoverflow.com/questions/30305542/using-positive-lookahead-regex-with-re2)
      - Can convert such expresion to using non-capturing group instead `(?:)`