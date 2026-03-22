# quasi examples

This directory contains embedded `quasi` examples written for the accepted L1
function shape.

- `fn` owns the `quasi` block
- `implements` links a function to its abstract relation
- `ensures` is used for direct postconditions when the example benefits from
  spelling them out

The sorting corpus is split into a shared vocabulary file plus per-algorithm
files:

- `sorting/sorting.l1` defines the common sequence vocabulary and generic
  sorting laws
- `sorting/*.l1` defines one algorithm relation, one correctness spec, and one
  `fn` with a `quasi using style ProceduralAlgorithm:` block

Validate the sorting corpus as a directory:

```bash
go run ./cmd/len-cli validate examples/quasi/sorting
```