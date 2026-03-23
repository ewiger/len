# calculator example

This directory turns the original calculator sketch into a concrete example with
three len layers plus a runnable Go implementation.

Contents:

- `calculator.l0` captures the plain-language intent and examples.
- `calculator.l1` captures the semantic model, relations, and quasi algorithms.
- `calculator.l2` captures target-specific Go implementation guidance.
- `calculator.go` is a hand-written Go implementation derived from the L1/L2
  design.
- `calculator_test.go` is a table-driven test suite derived from the L0
  examples.

The Go package is runnable today. The len files are design artifacts: they use
target-surface constructs such as `contract` and `struct`, which are ahead of
the current MVP parser surface in this repository.

## Run

From this directory:

```sh
make test
```

Or from the repository root:

```sh
go test ./examples/advanced/calculator
```

## Example cases

- `1 + 2 * 3` -> `7`
- `(1 + 2) * 3` -> `9`
- `8 / 2 + 1` -> `5`
- `8 / (2 + 2)` -> `2`
- `1 / 0` -> error
- `1 + )` -> parse error