# LIP-0002: Struct, Trait, and Relation Refinement Sugar for len.l1

## Metadata

| Field | Value |
| --- | --- |
| LIP | 0002 |
| Title | Struct, Trait, and Relation Refinement Sugar for len.l1 |
| Status | Draft |
| Author | TBD |
| Created | 2026-03-22 |
| Updated | 2026-03-22 |
| Area | Surface Syntax, Core Modeling |
| Target | len.l1 |

## Summary

This proposal clarifies how `struct`, `trait`, and relation-level refinement should appear in `len.l1` without changing the core language model.

The proposal keeps `type` and `rel` as the true L1 primitives. `spec` remains the general mechanism for arbitrary laws, definitions, and semantic constraints. `fn` remains the executable or constructive form, optionally paired with an embedded `quasi` block for open-ended pseudocode.

Within that model, this proposal makes four design moves:

- add relation-level `refines` syntax sugar
- remove `impl` from L1
- keep `trait`, but define it as grouped contract sugar rather than an implementation mechanism
- keep `struct`, but define it as record sugar rather than a class-like construct

The intent is to make the surface language more readable while keeping the semantics minimal, relational, and high-level.

## Proposal Documents

This proposal is split into focused documents:

- `README.md`: proposal overview, relation refinement, and migration away from `impl`
- `trait.md`: detailed design for `trait` as grouped contract sugar
- `struct.md`: detailed design for `struct` as composite type sugar

## Motivation

The current surface contains `trait`, `impl`, and `implements`, but those words pull the language toward programming-language intuitions that do not match `len` well.

In `len`, the important relationships are logical and semantic:

- one relation may refine another relation
- a collection of relations and functions may be grouped under a named contract
- a composite value may be described by named fields
- an executable function may be linked to a relation and may contain a `quasi` sketch

Those ideas do not require object-oriented classes, Rust-style trait implementations, dispatch rules, or nominal instance machinery. The existing words suggest those semantics even when the language does not intend to provide them.

The language should instead expose the higher-level ideas directly:

- logical refinement between relations
- grouped contracts over `rel`, `fn`, and `spec`
- composite data as structured `type` sugar

## Problem Statement

Today the L1 surface has three related issues.

First, the relation between an abstract contract and a more specific relation is described inconsistently. The notion is semantic refinement, but names like `implements` and `impl` imply a programming-language implementation relation.

Second, `impl` is both ambiguous and redundant. In a language that already has `spec`, `rel`, and `fn`, `impl` does not introduce a new essential idea. It mainly imports expectations from other languages.

Third, `trait` and `struct` are useful words for readers, but they are currently underspecified. Without a clear lowering model they risk becoming accidental new primitives instead of shallow sugar over the existing core.

The result is vocabulary drift: the surface starts to look object-oriented even though the actual design is relational and specification-first.

## Proposed Changes

### 1. Add Relation-Level `refines`

L1 should allow a relation declaration to state that it semantically refines another relation.

Example:

```len
rel BubbleSort(input: Seq, output: Seq) refines Sort(input, output)
```

This should be read as logical refinement, not as code-level implementation. A relation with a `refines` clause promises that whenever the refining relation holds, the refined relation also holds.

The canonical lowering is a `spec`:

```len
spec bubble_sort_correct
    given input, output: Seq
    must BubbleSort(input, output) implies Sort(input, output)
```

This proposal intentionally prefers `refines` over `implements` for relations because the relationship is semantic implication between predicates, not implementation in the programming-language sense.

### 2. Drop `impl` From L1

`impl` should be removed from the L1 surface.

The keyword is confusing for three reasons:

- it suggests instance or dispatch semantics that L1 does not define
- it overlaps conceptually with `spec`, `rel`, and `fn`
- it does not say what logical claim is actually being made

After this proposal, the language should express the important relationships directly:

- use `refines` for relation refinement
- use `trait` as group contract sugar or combination of many `fn` and `rel` declarations
- use `fn` clauses in trait but remember that more details can always be expressed in `spec`

If a previous `impl` declaration was intended to mean that a type satisfies a named contract, that meaning should be re-expressed through the desugared contract relations and specs rather than through a dedicated implementation keyword.

### 3. Keep `trait`, but Redefine It Clearly

`trait` should remain in the surface language, but only as grouped contract sugar.

The detailed trait design is specified in `trait.md`. In short:

- a `trait` groups related `rel`, `fn`, and `spec` declarations
- it may look interface-like in surface form
- it must not introduce inheritance, dispatch, or Rust-style implementation machinery
- its semantics remain grounded in ordinary `rel`, `fn`, and `spec`

### 4. Keep `struct`, but Define It as a Composite Type

`struct` should remain in L1 as lightweight composite-type sugar.

The detailed struct design is specified in `struct.md`. In short:

- a `struct` declares a composite type with named fields
- it behaves like a record or named tuple, not a class
- field syntax stays lightweight as direct `name: Type` entries
- the semantics remain reducible to `type` plus lower-level relational structure

## Grammar Sketch

The following sketch shows the intended surface direction. It is illustrative rather than final.

```ebnf
TopLevelDecl   = ImportDecl
               | TypeDecl
               | RelDecl
               | FnDecl
               | SpecDecl
               | StructDecl
               | TraitDecl
               | KeywordDecl
               | SymbolDecl ;

RelDecl        = "rel" Identifier "(" [ ParamList ] ")" [ RefinesClause ] ;
RefinesClause  = "refines" Formula ;
```

Notable consequences:

- `impl` is removed
- `implements` is not used for relation refinement, only in `fn` clauses when appropriate
- `trait` and `struct` remain explicit top-level surface forms
- detailed trait and struct grammar is specified in the companion documents

This proposal does require `trait` and `struct` to be represented explicitly in the parser and validator as top keywords.

## Desugaring / Lowering Model

### Relation Refinement

Surface:

```len
rel BubbleSort(input: Seq, output: Seq) refines Sort(input, output)
```

Lowering sketch:

```len
rel BubbleSort(input: Seq, output: Seq)

spec bubble_sort_refines_sort
    given input, output: Seq
    must BubbleSort(input, output) implies Sort(input, output)
```

The generated `spec` name is illustrative. Implementations may synthesize names deterministically.

### Trait as Grouped Contract

The detailed trait surface and lowering model are specified in `trait.md`.

At the overview level, the important point is that `trait` remains an explicit surface declaration while its meaning is still reducible to ordinary `rel`, `fn`, and `spec` declarations.

### Struct as Composite Type Sugar

The detailed struct surface and lowering model are specified in `struct.md`.

At the overview level, the important point is that `struct` remains an explicit surface declaration while its semantics stay grounded in `type` plus lower-level relational structure rather than class semantics.

## Examples

### Before / After: Sorting Relation

Before, the semantic intent is present but indirect:

```len
rel Sort(input: Seq, output: Seq)
rel BubbleSort(input: Seq, output: Seq)

spec bubble_sort_correct
    given input, output: Seq
    must BubbleSort(input, output) implies Sort(input, output)
```

After, the same intent can be stated directly at the relation declaration:

```len
rel Sort(input: Seq, output: Seq)
rel BubbleSort(input: Seq, output: Seq) refines Sort(input, output)
```

The lowered meaning remains the same `spec`-level implication.

### Before / After: Function Contract

Before:

```len
fn bubble_sort(input: Seq) -> output: Seq
    implements BubbleSort(input, output)
    ensures Sorted(output)
    ensures Permutation(input, output)
    quasi:
        let output := input
        note reorder output until sorted
        return output
```

After:

```len
rel BubbleSort(input: Seq, output: Seq) refines Sort(input, output)

fn bubble_sort(input: Seq) -> output: Seq
    ensures BubbleSort(input, output)
    ensures Sorted(output)
    ensures Permutation(input, output)
    quasi:
        let output := input
        ...
        return output
```

The executable form remains `fn`, and the open-ended algorithm sketch remains inside `quasi`. The abstract semantic relation is still expressed relationally.

Detailed trait and struct examples are provided in the companion documents.

## Rationale

This proposal preserves the existing `len` philosophy.

`type` and `rel` remain the core modeling tools. `spec` remains the general mechanism for laws, definitions, and semantic constraints. `fn` remains the executable or proof-oriented form, and may still include an embedded `quasi` section when a structured but open-ended procedure sketch is useful.

Within that architecture:

- `refines` gives the right name for semantic implication between relations
- removing `impl` reduces misleading programming-language baggage
- `trait` remains available as a readable grouping device without becoming a new semantic tower
- `struct` remains available as readable data-shape sugar without becoming an OO class

This keeps the language minimal while still making common patterns concise.

## Backward-Compatibility Notes

This proposal changes the preferred surface language and therefore requires migration for currently observed syntax.

- `impl` declarations should be rejected or deprecated in L1 once replacement lowering is implemented
- relation-level semantics previously described informally through `implements` should move to `refines` or explicit `spec`
- `fn ... implements ...` clauses remains as is
- existing parsers and validators that recognize `impl` and `implements` do not need a transition strategy - it's MVP

Because current usage is still limited, this is a good point to simplify the vocabulary before those forms become entrenched.

## Open Questions

1. Should `fn` keep a dedicated relation-linking clause, or is `ensures` alone sufficient once relation refinement is explicit?
2. Should `refines` accept only relation applications, or any formula whose free variables match the relation signature?
3. What transition strategy should parsers and validators use while `impl` and older `implements` usage still exists in the repository?

## Conclusion

L1 should remain a high-level specification language centered on `type`, `rel`, `spec`, and `fn`, with `quasi` as an embedded auxiliary notation rather than a competing language tier.

Under that design, `trait` and `struct` are useful only if they stay as shallow surface sugar. Likewise, semantic implication between relations should be named `refines`, not described with implementation vocabulary.

This proposal therefore keeps the readable surface forms while removing the misleading semantics that normally travel with them in other languages.