# len

`len` is a meta programming language.

It is intended for writing programs that generate, transform, or reason about other programs and code structures.

## Status

This project is **early-stage and experimental**.

- Expect breaking changes to syntax, semantics, and tooling.
- The core language and features are still being designed and implemented but you can check out `lang/l1/**`.

## Layers

- **Level 0** or `len.l0` remains in natural language and informal pseudo-code. It is intended for high-level design, brainstorming, and communication.

- **Level 1** or `len.l1` defines the canonical structural core of len: types, relations, formulas, and syntax-level logical constructs. It is intentionally lightweight and declarative. 

- **Level 2** or `len.l2` builds on this core by introducing interpretation and evaluation, such as contexts, satisfaction, and executable semantic rules.

## Developer

The repository now includes the first Go implementation for the accepted CLI parser and validator milestone.

Current implementation:

- `len-cli validate` for parsing and validating `.l1` files
- hand-written lexer and parser for the current MVP declaration surface
- semantic validation for duplicate declarations, unresolved imports, binder scope, unresolved names, and function or relation arity
- quasi surface validation through YAML style profiles

Build the CLI:

```bash
go build ./cmd/len-cli
```

Run all unit tests:

```bash
go test ./...
```

Run the validator on the hello world example:

```bash
go run ./cmd/len-cli validate examples/helloworld/hello.l1
```

Current scope limits:

- `lang/l1/**` is intentionally not the acceptance target yet
- quasi bodies are captured as raw indented lines and surface-validated, not parsed into a dedicated internal statement AST
- the tool validates structure only; it does not execute code, interpret formulas, or perform proof checking
