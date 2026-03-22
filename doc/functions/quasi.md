# quasi

`quasi` is a required embedded block on an `fn` declaration.

In the accepted level 1 direction, the host parser does not parse the inside of a quasi block as a dedicated quasi grammar. Instead, it captures the block as raw indented lines and records any explicit style named in the header. Later, validation resolves the style profile and performs style-specific surface validation.

This keeps the core `len.l1` grammar stable while still allowing different pseudocode styles to be checked in a controlled way.

## Accepted MVP Model

The accepted proposal treats quasi as:

- a required clause owned by `fn`
- a header of the form `quasi:` or `quasi using style Name:`
- an indentation-sensitive raw block captured by the host parser
- a body validated later by a style profile

The accepted MVP does not define a general-purpose inner quasi statement grammar in the host parser. It also does not require the parser to build a dedicated quasi statement AST.

The source of truth for that behavior is:

- [doc/proposals/accepted/lip-0001-cli-parser-n-validation/README.md](../proposals/accepted/lip-0001-cli-parser-n-validation/README.md)
- [doc/proposals/accepted/lip-0001-cli-parser-n-validation/grammar.md](../proposals/accepted/lip-0001-cli-parser-n-validation/grammar.md)
- [doc/proposals/accepted/lip-0001-cli-parser-n-validation/quasi-styles.md](../proposals/accepted/lip-0001-cli-parser-n-validation/quasi-styles.md)

## Host Embedding

`quasi` is attached only to an `fn` declaration. It is not a preferred top-level declaration form.

Preferred shape:

```text
fn bubble_sort(input: Seq) -> output: Seq
    implements BubbleSort(input, output)
    ensures Sorted(output)
    ensures Permutation(input, output)

    quasi using style ProceduralAlgorithm:
        let output := input
        while needs_swap(output):
            set output := bubble_pass(output)
        return output
```

The `fn` declaration owns:

- the signature
- contract clauses such as `requires`, `ensures`, and `implements`
- the required `quasi` block

This is the accepted replacement for older top-level quasi ideas.

## Parser Contract

At parse time, the host grammar is responsible for:

- recognizing a quasi clause inside `fn`
- parsing the header, including an optional `using style Name`
- capturing the indented body as raw lines
- recording source spans and indentation metadata

At parse time, the host grammar is not responsible for:

- parsing quasi lines into a general statement grammar
- inventing a built-in proof-sketch or algorithm language
- deciding globally whether unknown lines are acceptable

Those later checks belong to style-driven surface validation.

## Surface Validation

After parse, the validator resolves the style profile for the block.

- `quasi:` uses the implementation's configured default style profile
- `quasi using style Name:` selects the named style profile

The style profile then defines the surface contract for that block, such as:

- allowed leading keywords
- accepted line shapes, typically via regular expressions
- indentation rules
- block-opening and continuation behavior
- whether narrative lines are accepted or rejected

Unknown-line handling is style-dependent. It is not a universal rule of the host grammar.

## Current Concrete Style

The accepted proposal currently defines one concrete example style profile:

- `ProceduralAlgorithm`

That profile is documented here:

- [doc/proposals/accepted/lip-0001-cli-parser-n-validation/quasi-styles.md](../proposals/accepted/lip-0001-cli-parser-n-validation/quasi-styles.md)
- [doc/proposals/accepted/lip-0001-cli-parser-n-validation/procedural-algorithm.quasi-style.yaml](../proposals/accepted/lip-0001-cli-parser-n-validation/procedural-algorithm.quasi-style.yaml)

It is derived from the sorting examples under `examples/quasi/sorting/**` and is intentionally narrow. In the accepted YAML profile, `ProceduralAlgorithm` forbids free-form narrative lines and validates the block purely at the surface level.

## Relationship To `fn` And `spec`

The accepted split is:

- `fn` carries function-like signatures, contract clauses, and required `quasi`
- `spec` remains declarative and uses `given` and `must`

`quasi` therefore belongs with operational or implementation-oriented sketches attached to an `fn`, not with general declarative specification text.

## Scope And Limits

The accepted MVP for quasi is intentionally small.

What is in scope:

- raw quasi block capture in the parser
- optional style metadata in the header
- style-profile loading
- style-specific surface validation after parse

What is not in scope for the host grammar in MVP:

- a built-in universal quasi language
- a dedicated inner quasi statement AST
- semantic interpretation or execution of quasi blocks
- parser extensibility driven directly by quasi content

Additional styles may be added later, but they should still be expressed as style profiles consumed by validation rather than as ad hoc extensions to the host grammar.

## Recommended Reading

For the accepted design, read these documents together:

- [doc/proposals/accepted/lip-0001-cli-parser-n-validation/README.md](../proposals/accepted/lip-0001-cli-parser-n-validation/README.md)
- [doc/proposals/accepted/lip-0001-cli-parser-n-validation/grammar.md](../proposals/accepted/lip-0001-cli-parser-n-validation/grammar.md)
- [doc/proposals/accepted/lip-0001-cli-parser-n-validation/quasi-styles.md](../proposals/accepted/lip-0001-cli-parser-n-validation/quasi-styles.md)

That proposal set is the authoritative description of the current quasi direction for level 1.
