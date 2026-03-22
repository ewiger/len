# LIP-0001: CLI Parser and Validation for len.l1

## Metadata

| Field | Value |
| --- | --- |
| LIP | 0001 |
| Title | CLI Parser and Validation for len.l1 |
| Status | Accepted |
| Author | TBD |
| Created | 2026-03-21 |
| Updated | 2026-03-21 |
| Area | CLI, Parser, Validation |
| Target | len-cli validate |

## Summary

This proposal introduces a Go-based `len-cli` focused on parsing and validating level 1 `len` source files.

The accepted first milestone is intentionally narrow:

- implement a parser for the current `lang/l1/**` corpus
- implement semantic validation for names, imports, arity, and binder scope
- expose that functionality through `len-cli validate`
- add unit tests and a minimal hello-world example for later integration testing

The first version does not include interpretation, execution, proof checking, or `l1 -> l2` transformation.

## Motivation

The repository already contains a meaningful `lang/l1/**` corpus, but it lacks tooling that can answer basic structural questions:

- does a file parse
- do imported modules resolve
- are declarations internally consistent
- do formulas use names and binders correctly

Without a parser and validator, the language design cannot be exercised with confidence. A CLI validator provides the minimum practical tool for evolving the language while keeping syntax drift visible.

## Goals

- parse the current level 1 corpus under `lang/l1/**`
- provide deterministic diagnostics with file, line, and column information
- separate syntax errors from semantic validation errors
- support module loading for current import forms such as `import core.math.set`
- make parser and validator behavior testable with unit tests and corpus tests
- add a minimal example under `examples/helloworld/**` for later integration coverage

## Non-Goals

- user-extensible parsing driven by `syntax` declarations at runtime
- evaluation or interpretation of formulas
- proof verification or consistency checking between `l1` and `l2`
- code generation
- full language stabilization beyond the current corpus

## Accepted Decisions

- implementation language: Go
- parser style: hand-written lexer plus hand-written parser
- validation style: separate semantic validation pass after parse
- first CLI milestone: `len-cli validate`
- grammar documentation: keep a human-readable grammar document alongside the implementation
- quasi strategy: parse quasi headers in `fn`, capture quasi bodies as raw indented lines, and perform style-specific surface validation after parse

## Proposal Documents

This proposal is split into focused documents:

- `README.md`: proposal overview and accepted decisions
- `plan.md`: exact implementation steps and target files
- `grammar.md`: lexical rules, comments/docstrings, and MVP grammar for `len.l1`
- `quasi-styles.md`: style-profile model and surface-validation strategy for quasi blocks
- `procedural-algorithm.quasi-style.yaml`: first concrete style profile derived from the sorting examples

## Current Language Surface

The current corpus shows a stable enough MVP surface for a first parser.

### Top-Level Declaration Forms

- `import`
- `type`
- `rel`
- `fn`
- `const`
- `spec`
- `syntax`
- `trait`
- `impl`
- `keyword`
- `symbol`

`quasi` is no longer best treated as a top-level declaration form. The preferred design is an embedded `quasi:` block owned by an `fn`.

### Fn Surface

The preferred `fn` direction is function-like:

- parameter list in the header
- optional result binder after `->`
- contract clauses such as `requires`, `ensures`, and `implements`
- optional embedded `quasi:` block for an algorithm sketch or proof-oriented implementation outline
- optional explicit style marker such as `quasi using style ProceduralAlgorithm:`
- quasi bodies are captured as raw indented text and validated later against a style profile rather than fully parsed by the host grammar

### Spec Surface

`spec` remains declarative rather than algorithmic:

- zero or more `given` clauses
- one `must` clause
- no embedded `quasi:` block

### Formula and Expression Forms

- typed binders with `forall` and `exists`
- logical connectives: `and`, `or`, `not`, `implies`, `iff`
- sugar operators: `=>`, `<=>`, `/\`, `\/`, `!`
- equality with `=`
- relation application such as `Member(x, s)`
- infix syntax such as `x in s` and `x subsetof y`
- grouped formulas and expressions

## Comments and Docstrings

Comments and docstrings are allowed in source and ignored by formal semantics.

Accepted documentation forms for `len.l1` are:

- docstring literals as standalone string forms
- line comments starting with `#` at the beginning of the line
- block comments delimited by `<#` and `#>`

Docstrings are preferred over line comments when the intent is structured documentation.

Examples:

```text
"single-line docstring"

"""
Multi-line docstring.
Ignored by formal semantics.
Useful for intent, examples, and notes.
"""

# line comment

<#
block comment
#>
```

### Relevant Reference Files

- [README.md](/Users/yy/code/len/len-feat-cli-l1-validation/README.md)
- [TODO](/Users/yy/code/len/len-feat-cli-l1-validation/TODO)
- [lang/l1/README.md](/Users/yy/code/len/len-feat-cli-l1-validation/lang/l1/README.md)
- [lang/l1/core/syntax/keywords.l1](/Users/yy/code/len/len-feat-cli-l1-validation/lang/l1/core/syntax/keywords.l1)
- [lang/l1/core/syntax/symbols.l1](/Users/yy/code/len/len-feat-cli-l1-validation/lang/l1/core/syntax/symbols.l1)
- [lang/l1/core/syntax/syntax.l1](/Users/yy/code/len/len-feat-cli-l1-validation/lang/l1/core/syntax/syntax.l1)
- [lang/l1/core/types/types.l1](/Users/yy/code/len/len-feat-cli-l1-validation/lang/l1/core/types/types.l1)
- [lang/l1/core/math/logic/logic.l1](/Users/yy/code/len/len-feat-cli-l1-validation/lang/l1/core/math/logic/logic.l1)
- [lang/l1/core/math/logic/sugar.l1](/Users/yy/code/len/len-feat-cli-l1-validation/lang/l1/core/math/logic/sugar.l1)
- [lang/l1/core/math/set/set.l1](/Users/yy/code/len/len-feat-cli-l1-validation/lang/l1/core/math/set/set.l1)
- [lang/l1/core/math/set/nat.l1](/Users/yy/code/len/len-feat-cli-l1-validation/lang/l1/core/math/set/nat.l1)
- [lang/l1/core/math/fun/fun.l1](/Users/yy/code/len/len-feat-cli-l1-validation/lang/l1/core/math/fun/fun.l1)

## Design Rationale

The accepted approach remains a hand-written lexer and parser plus a separate validation phase. That is the best fit for the current corpus because the syntax is still changing and includes context-sensitive constructs such as `spec`, `syntax`, and `quasi`.

For quasi specifically, the accepted direction is intentionally shallow at parse time:

- the host parser recognizes the quasi clause header and its indentation-delimited body
- the body is stored as raw lines with source spans and indentation metadata
- a later validation pass resolves the style and applies a surface-validation profile such as [procedural-algorithm.quasi-style.yaml](/Users/yy/code/len/len-feat-cli-l1-validation/doc/proposals/accepted/lip-0001-cli-parser-n-validation/procedural-algorithm.quasi-style.yaml)

This keeps the core grammar stable while still allowing repository-local or custom quasi styles.

## Acceptance Criteria

This proposal is considered implemented when:

- a Go-based `len-cli validate` command exists
- the current `lang/l1/**` corpus parses under the supported grammar
- semantic validation reports meaningful diagnostics
- unit tests exist for lexer, parser, and validator behavior
- a hello-world example exists under `examples/helloworld/**`

## Future Work

- add `transform` and `verify` commands
- define the interface between `l1` structural validation and future `l2` semantics
- revisit whether `syntax` declarations should ever influence parsing dynamically
- formalize the grammar further once the surface syntax stabilizes