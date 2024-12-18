# [Advent of Code 2024](https://adventofcode.com/2024)

This repo contains my attempt at Advent of Code 2024 with Golang. Inside there will be code solutions separated by the advent days folder.


## Learnings

- Day 2 p2
    - `append()` slice function modifies the original slice
- Day 3
    - Golang regex are not able to do lookaheads `(?=)` or lookbehind `(?<=)` due to it using [RE2](https://stackoverflow.com/questions/30305542/using-positive-lookahead-regex-with-re2)
      - Can convert such expresion to using non-capturing group instead `(?:)`
- Day 9
    - Slices is always pass by reference while arrays are pass by value. Use `copy()` to copy slices into new variables.
    - `fmt.Printf("%+v\n", v)` use `%+v` to print the struct key as well
- Day 14
    - When iterating over a slice, the `v` value from `i,v := range ...` is copied by value. So any update wont be affect the underlying array