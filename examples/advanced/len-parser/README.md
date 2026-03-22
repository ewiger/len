# len.parser specification tree

This directory contains implementation-facing parser specifications under the
`len.parser.*` namespace.

The namespace is encoded by the filesystem path:

- `examples/advanced/len-parser/len/parser/lexical/*` -> `len.parser.lexical.*`
- `examples/advanced/len-parser/len/parser/declarations/*` -> `len.parser.declarations.*`
- `examples/advanced/len-parser/len/parser/expressions/*` -> `len.parser.expressions.*`
- `examples/advanced/len-parser/len/parser/ast/*` -> `len.parser.ast.*`
- `examples/advanced/len-parser/len/parser/loading/*` -> `len.parser.loading.*`
- `examples/advanced/len-parser/len/parser/resolution/*` -> `len.parser.resolution.*`
- `examples/advanced/len-parser/len/parser/scope/*` -> `len.parser.scope.*`
- `examples/advanced/len-parser/len/parser/quasi/*` -> `len.parser.quasi.*`

## Initial file list

- `len/parser/lexical/tokens.l0`
- `len/parser/lexical/tokens.l1`
- `len/parser/lexical/trivia.l0`
- `len/parser/lexical/trivia.l1`
- `len/parser/declarations/top_level_forms.l0`
- `len/parser/declarations/top_level_forms.l1`
- `len/parser/declarations/functions.l0`
- `len/parser/declarations/functions.l1`
- `len/parser/expressions/core.l0`
- `len/parser/expressions/core.l1`
- `len/parser/expressions/quantifiers_and_precedence.l0`
- `len/parser/expressions/quantifiers_and_precedence.l1`
- `len/parser/ast/surface_nodes.l0`
- `len/parser/ast/surface_nodes.l1`
- `len/parser/loading/import_resolution.l0`
- `len/parser/loading/import_resolution.l1`
- `len/parser/resolution/namespaces.l0`
- `len/parser/resolution/namespaces.l1`
- `len/parser/resolution/arity_and_unknowns.l0`
- `len/parser/resolution/arity_and_unknowns.l1`
- `len/parser/scope/binders.l0`
- `len/parser/scope/binders.l1`
- `len/parser/quasi/procedural_surface.l0`
- `len/parser/quasi/procedural_surface.l1`

## Scope rules for this tree

- `.l0` files describe parser or validator intent in plain language.
- `.l1` files encode parser-facing specimens grounded in `internal/**/*.go`.
- Files may be acceptance-oriented or diagnostic-oriented.
- Acceptance-oriented files should remain inside the current parser surface.
- Diagnostic-oriented files may intentionally fail validation if the file says so.

## Current implementation boundary

This tree mirrors the current Go implementation, not the broader future L1
surface. In particular, `contract`, `struct`, `refines`, `satisfies`, and
`derives` are out of scope here unless explicitly introduced in a future-only
subtree.