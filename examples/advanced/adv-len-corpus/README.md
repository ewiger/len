# len parser-spec corpus

This directory is a parser-design corpus for `len.l1`.

It is organized as a specification tree rather than a single runnable example. The
goal is to expose the current parser surface, the current language corpus, and the
intended next-surface in one place so future parser work has a concrete target.

## Project stage analysis

- The Go implementation is still at the accepted LIP-0001 MVP stage.
- The current parser accepts these top-level forms: `import`, `type`, `rel`,
	`const`, `spec`, `syntax`, `trait`, `impl`, `keyword`, `symbol`, and `fn`.
- `fn` currently requires an indented body with optional `requires`, `ensures`,
	and `implements` clauses followed by a required `quasi` block.
- Inline expressions are already richer than the declaration surface: names,
	qualified names, calls, strings, numbers, booleans, grouping, quantifiers,
	unary operators, and binary operators are implemented.
- The repository documentation and `lang/l1/**` corpus are ahead of the parser in
	some places. In particular, `contract` and `struct` appear in docs and in the
	language corpus, but the current Go parser does not yet accept them.
- The existing runnable examples under `examples/basic/**` and `examples/quasi/**`
	stay inside the current MVP parser surface.

## Hierarchy

- `10-current-mvp`: lexical trivia and top-level declarations that the current
	parser already accepts.
- `20-current-mvp`: `spec` blocks and inline formula syntax that the current
	expression parser already handles.
- `30-current-mvp`: `fn` bodies and `quasi` blocks aligned with the current host
	parser and quasi style validator.
- `40-target-l1`: grouped `contract` syntax as described in the L1 docs and
	proposals.
- `50-target-l1`: explicit `struct` syntax and field bodies.
- `60-target-l1`: relation refinement plus `satisfies` and `derives` bridging.
- `70-target-l1`: module and namespace-oriented examples for future import and
	qualification work.

## How To Use This Tree

- Treat `10-current-mvp` through `30-current-mvp` as acceptance-oriented examples
	for the parser and validator that exist today.
- Treat `40-target-l1` through `70-target-l1` as design targets for the broader
	`len.l1` surface described by docs, proposals, and the `lang/l1` corpus.
- Do not expect the whole tree to validate with the current CLI yet. The target
	subtree is intentionally ahead of the current implementation.

## Recommended parser roadmap implied by this corpus

1. Keep the current MVP declaration surface stable.
2. Preserve expression compatibility with the current `spec` and `syntax` forms.
3. Add explicit AST support for `contract` and `struct` instead of treating them
	 as ad hoc sugar outside the parser.
4. Add optional relation-level `refines` and struct-level `satisfies` and
	 `derives` clauses.
5. Extend module import and qualification rules so future example trees can be
	 validated as multi-file L1 corpora rather than only as isolated sketches.