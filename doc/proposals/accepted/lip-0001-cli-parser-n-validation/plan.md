# LIP-0001 Implementation Plan

## Scope

This plan describes the exact implementation steps and target files for the accepted `len-cli` parser and validation MVP.

The scope of this plan is:

- parser for the current `lang/l1/**` corpus
- semantic validation after parse
- `len-cli validate`
- unit tests and corpus tests
- a minimal hello-world example

## Proposed Go Layout

```text
cmd/
  len-cli/
    main.go
internal/
  ast/
    ast.go
  diag/
    diag.go
    span.go
  lexer/
    lexer.go
    token.go
  loader/
    loader.go
  parser/
    parser.go
    declarations.go
    expressions.go
  quasi/
    profile.go
    surface.go
  validator/
    symbols.go
    validator.go
examples/
  helloworld/
    hello.l1
testdata/
  valid/
  invalid/
```

## Phase 1: Bootstrap Module and CLI

### Goal

Create the Go module and the initial command entrypoint.

### Files

- `go.mod`
- `cmd/len-cli/main.go`

### Steps

1. initialize the Go module at repository root
2. add a `main.go` entrypoint for `len-cli`
3. define a `validate` command path and command usage output
4. keep CLI execution thin and move logic into internal packages

## Phase 2: Diagnostics and Source Spans

### Goal

Introduce a stable diagnostic model before parser logic grows.

### Files

- `internal/diag/span.go`
- `internal/diag/diag.go`

### Steps

1. define `Position` with file, line, and column
2. define `Span` with start and end positions
3. define `Diagnostic` with code, message, severity, and span
4. add formatting helpers for CLI output

## Phase 3: Lexer

### Goal

Tokenize current `len.l1` source with accurate location tracking.

### Files

- `internal/lexer/token.go`
- `internal/lexer/lexer.go`
- `internal/lexer/lexer_test.go`

### Inputs

- [lang/l1/core/syntax/keywords.l1](/Users/yy/code/len/len-feat-cli-l1-validation/lang/l1/core/syntax/keywords.l1)
- [lang/l1/core/syntax/symbols.l1](/Users/yy/code/len/len-feat-cli-l1-validation/lang/l1/core/syntax/symbols.l1)
- [doc/proposals/accepted/lip-0001-cli-parser-n-validation/grammar.md](/Users/yy/code/len/len-feat-cli-l1-validation/doc/proposals/accepted/lip-0001-cli-parser-n-validation/grammar.md)

### Steps

1. define token kinds for keywords, identifiers, strings, punctuation, operators, comments, and EOF
2. support reserved words from the current corpus, including `fn`, declarative spec words such as `given` and `must`, and contract words such as `requires`, `ensures`, and `implements`
3. support multi-character operators such as `->`, `:=`, `=>`, `<=>`, `/\`, and `\/`
4. support comments:
   `#` line comments only when `#` starts the line
   `<# ... #>` block comments
5. support docstrings:
   single-line string docstrings
   triple-quoted multi-line docstrings
6. track line and column across all tokens and ignored trivia
7. add lexer tests for keywords, comments, docstrings, and operators

## Phase 4: AST

### Goal

Define a canonical AST for the current parser scope.

### Files

- `internal/ast/ast.go`

### Steps

1. add top-level declaration nodes for `import`, `type`, `rel`, `fn`, `const`, `spec`, `syntax`, `trait`, `impl`, `keyword`, and `symbol`
2. extend `fn` nodes with signature information, result binders, contract clauses such as `requires`, `ensures`, and `implements`, and a required embedded raw `QUASI_BLOCK` plus optional style metadata for quasi blocks
3. represent a quasi block as raw lines, indentation metadata, and source spans rather than as a fully parsed statement AST
4. add expression nodes for identifiers, qualified names, applications, equality, infix expressions, unary expressions, quantifiers, and grouped expressions
5. add explicit nodes for comments and docstrings only if the parser needs to preserve them for tooling
6. attach spans to major nodes

## Phase 5: Parser

### Goal

Parse declarations and formulas from the current corpus.

### Files

- `internal/parser/parser.go`
- `internal/parser/declarations.go`
- `internal/parser/expressions.go`
- `internal/parser/parser_test.go`

### Inputs

- [doc/proposals/accepted/lip-0001-cli-parser-n-validation/grammar.md](/Users/yy/code/len/len-feat-cli-l1-validation/doc/proposals/accepted/lip-0001-cli-parser-n-validation/grammar.md)

### Steps

1. parse top-level declarations in source order
2. parse `spec` with zero or more `given` clauses followed by one `must` clause
3. parse `fn` with a signature, optional result binder, contract clauses such as `requires`, `ensures`, and `implements`, and a required embedded `quasi` clause with optional `using style <Name>` header metadata
4. parse `syntax` declarations as surface form, binder list, and canonical form
5. capture embedded quasi bodies as indentation-sensitive raw `QUASI_BLOCK` regions using a dedicated block routine in the main parser, without parsing style-internal statements
6. implement precedence-based expression parsing
7. treat comments and docstrings as non-semantic trivia and ignore them during formal parse
8. add parser tests for every declaration form present in the current corpus

## Phase 5.5: Quasi Style Profiles

### Goal

Load quasi style profiles and prepare them for surface validation.

### Files

- `internal/quasi/profile.go`
- `internal/quasi/surface.go`
- `internal/quasi/profile_test.go`

### Inputs

- [doc/proposals/accepted/lip-0001-cli-parser-n-validation/quasi-styles.md](/Users/yy/code/len/len-feat-cli-l1-validation/doc/proposals/accepted/lip-0001-cli-parser-n-validation/quasi-styles.md)
- [doc/proposals/accepted/lip-0001-cli-parser-n-validation/procedural-algorithm.quasi-style.yaml](/Users/yy/code/len/len-feat-cli-l1-validation/doc/proposals/accepted/lip-0001-cli-parser-n-validation/procedural-algorithm.quasi-style.yaml)

### Steps

1. define Go structs for style profiles, slot definitions, rule definitions, and validation settings
2. load YAML style profiles with `gopkg.in/yaml.v3`
3. compile regular expressions once during profile load
4. expose a surface-validation API that accepts raw quasi lines plus a style profile
5. add tests for profile loading, regex compilation, and block-structure validation

## Phase 6: Loader

### Goal

Resolve imports from current module paths to repository files.

### Files

- `internal/loader/loader.go`
- `internal/loader/loader_test.go`

### Steps

1. map module path `a.b.c` to the current `lang/l1` file layout
2. load a root file and its transitive imports
3. detect missing imports and import cycles if present
4. return source plus diagnostics to the parser and validator

## Phase 7: Validator

### Goal

Validate semantic well-formedness after parse.

### Files

- `internal/validator/symbols.go`
- `internal/validator/validator.go`
- `internal/validator/validator_test.go`

### Steps

1. define symbol tables for types, relations, constants, traits, and local binders
2. detect duplicate top-level names per namespace
3. validate relation arity and parameter references
4. validate binder scope in formulas
5. validate unresolved references in `syntax` and, where enabled, in host expressions or formulas captured from quasi style rules
6. validate `spec` structure and required `must` clause
7. validate `fn` structure, including contract clauses, result binders, required quasi presence, quasi header correctness, and style-profile conformance
8. surface-validate raw quasi blocks by resolving the declared style, applying its profile, and reporting lexical or structural diagnostics with source spans

## Phase 8: CLI Wiring

### Goal

Expose parser and validator through `len-cli validate`.

### Files

- `cmd/len-cli/main.go`

### Steps

1. accept one or more `.l1` paths
2. invoke loader, parser, and validator
3. print diagnostics in a stable human-readable format
4. return non-zero on failure

## Phase 9: Tests and Fixtures

### Goal

Cover the parser and validator with unit tests and corpus-backed tests.

### Files

- `internal/lexer/lexer_test.go`
- `internal/parser/parser_test.go`
- `internal/loader/loader_test.go`
- `internal/validator/validator_test.go`
- `testdata/valid/*.l1`
- `testdata/invalid/*.l1`

### Steps

1. add table-driven lexer tests
2. add parser tests for precedence, declaration structure, quasi headers, and raw quasi block capture
3. add validator tests for duplicates, unresolved names, arity mismatches, and quasi style-profile failures
4. add corpus tests that parse `lang/l1/**`
5. add negative fixtures that assert expected diagnostic codes

## Phase 10: Hello World Example

### Goal

Create a minimal valid module for later integration testing.

### Files

- `examples/helloworld/hello.l1`

### Suggested Content

```text
"""
Minimal hello-world style module for parser and validator bring-up.
Ignored by formal semantics.
"""

type Greeting
rel Hello(x: Greeting)

spec hello_exists
  ensures exists x: Greeting . Hello(x)
```

## Phase 11: Documentation

### Goal

Keep repository docs aligned with implementation.

### Files

- `README.md`
- `TODO`
- `doc/proposals/accepted/lip-0001-cli-parser-n-validation/README.md`
- `doc/proposals/accepted/lip-0001-cli-parser-n-validation/plan.md`
- `doc/proposals/accepted/lip-0001-cli-parser-n-validation/grammar.md`

### Steps

1. document the implemented `validate` command
2. document the supported grammar subset
3. update milestones as parser work lands

## Exit Criteria

- `len-cli validate` exists and runs
- the current `lang/l1/**` corpus parses under the supported subset
- semantic validation emits useful diagnostics
- unit tests cover lexer, parser, loader, and validator
- `examples/helloworld/hello.l1` exists for integration use