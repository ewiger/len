# LIP-0002: Struct, Contract, and Relation Refinement Sugar for len.l1

## Metadata

| Field | Value |
| --- | --- |
| LIP | 0002 |
| Title | Struct, Contract, and Relation Refinement Sugar for len.l1 |
| Status | Draft |
| Author | TBD |
| Created | 2026-03-22 |
| Updated | 2026-03-22 |
| Area | Surface Syntax, Core Modeling |
| Target | len.l1 |

## Summary

This proposal clarifies how `struct`, `contract`, contract satisfaction, and relation-level refinement should appear in `len.l1` without changing the core language model.

The proposal keeps `type` and `rel` as the true L1 primitives. `spec` remains the general mechanism for arbitrary laws, definitions, and semantic constraints. `fn` remains the executable or constructive form, optionally paired with an embedded `quasi` block for open-ended pseudocode.

Within that model, this proposal makes six design moves:

- add relation-level `refines` syntax sugar
- remove `impl` from L1
- replace `trait` with `contract`, defined as grouped contract sugar
- keep `struct`, but define it as record sugar rather than a class-like construct
- add `satisfies` for explicit contract adoption
- keep `derives` as the structural companion to `satisfies`

The intent is to make the surface language more readable while keeping the semantics minimal, relational, and high-level.

## Proposal Documents

This proposal is split into focused documents:

- `README.md`: proposal overview, relation refinement, and migration away from `impl`
- `contract.md`: detailed design for `contract` as grouped contract sugar
- `struct.md`: detailed design for `struct` as composite type sugar
- `derive.md`: detailed design for `satisfies` and `derives` as the bridge from structure to contracts

## Motivation

The current surface contains `trait`, `impl`, and `implements`, but those words pull the language toward programming-language intuitions that do not match `len` well.

In `len`, the important relationships are logical and semantic:

- one relation may refine another relation
- a collection of relations and functions may be grouped under a named contract
- a composite value may be described by named fields
- an executable function may be linked to a relation and always contains a `quasi` sketch

Those ideas do not require object-oriented classes, Rust-style trait implementations, dispatch rules, or nominal instance machinery. The existing words suggest those semantics even when the language does not intend to provide them.

The language should instead expose the higher-level ideas directly:

- logical refinement between relations
- grouped contracts over `rel`, `fn`, and `spec`
- composite data as structured `type` sugar

## Problem Statement

Today the L1 surface has three related issues.

First, the relation between an abstract contract and a more specific relation is described inconsistently. The notion is semantic refinement, but names like `implements` and `impl` imply a programming-language implementation relation.

> Yes, `implements` is still used in `fn` declarations, it is still a great fit for the intended meaning there. A function can be read as implementing a relation, and the `implements` clause can be read as a contract linking the executable form to the abstract relation. See the open question about whether `ensures Relation(...)` should also be allowed as equivalent surface sugar for that purpose. The key point is that `implements` is appropriate for linking a function to a relation as its contract, while `refines` is more appropriate for stating that one relation semantically implies another.

Second, `impl` is both ambiguous and redundant. In a language that already has `spec`, `rel`, and `fn`, `impl` does not introduce a new essential idea. It mainly imports expectations from other languages.

Third, `trait` is still carrying unnecessary foreign-language baggage even though its intended meaning is already grouped contract sugar. `struct` also needs a clearer bridge to those grouped contracts. Without a clear lowering model both forms risk becoming accidental new primitives instead of shallow sugar over the existing core.

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
- use `contract` as grouped contract sugar over `rel`, `fn`, and `spec`
- use `satisfies` when a type or structure adopts a contract
- use `fn` clauses inside contracts when appropriate, while keeping `spec` available for additional laws

If a previous `impl` declaration was intended to mean that a type satisfies a named contract, that meaning should be re-expressed through `satisfies`, `derives`, and the desugared contract relations and specs rather than through a dedicated implementation keyword.

### 3. Replace `trait` with `contract`

`contract` should replace `trait` in the surface language as the explicit grouped-contract form.

The detailed contract design is specified in `contract.md`. In short:

- a `contract` groups related `rel`, `fn`, and `spec` declarations
- it may still look interface-like in surface form
- it introduces a same-named satisfaction relation `Eq(T)` 
that satisfy the contract 
- it also introduces a generic relation predicate or `Eq(T,x)` (where `x` is a value of type `T`) that can be used in Formula and `spec` declarations
- it must not introduce inheritance, dispatch, instance lookup, or Rust-style implementation machinery
- its semantics remain grounded in ordinary `rel`, `fn`, and `spec`

### 4. Keep `struct`, but Define It as a Composite Type

`struct` should remain in L1 as lightweight composite-type sugar.

The detailed struct design is specified in `struct.md`. In short:

- a `struct` declares a composite type with named fields
- it behaves like a record or named tuple, not a class
- field syntax stays lightweight as direct `name: Type` entries
- the semantics remain reducible to `type` plus lower-level relational structure

### 5. Add `satisfies` for Contract Adoption

A type or structure should be able to declare that it adopts a contract explicitly.

Example:

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)

struct Point satisfies Eq(Point)
    x: Real
    y: Real
```

This should be read as contract adoption, not instance registration. It says that `Point` is intended to satisfy the contract proposition `Eq(Point)`, and that elaboration or validation must account for the corresponding relations and laws.

### 6. Keep `derives` as the Structural Companion to `satisfies`

`derives` should remain the mechanism for generated satisfaction artifacts when the contract view follows canonically from structure.

Example:

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)

struct Point derives Eq(Point)
    x: Real
    y: Real
```

Here `derives` means that the contract adoption is structural and that elaboration may generate the corresponding relations and defining specs from the fields of `Point`.

## Grammar Sketch

The following sketch shows the intended surface direction. It is illustrative rather than final.

```ebnf
TopLevelDecl   = ImportDecl
               | TypeDecl
               | RelDecl
               | FnDecl
               | SpecDecl
               | StructDecl
               | ContractDecl
               | KeywordDecl
               | SymbolDecl ;

RelDecl        = "rel" Identifier "(" [ ParamList ] ")" [ RefinesClause ] ;
RefinesClause  = "refines" Formula ;
ContractDecl   = "contract" Identifier [ "(" [ ParamList ] ")" ] ContractBody ;
StructDecl     = "struct" Identifier [ SatisfiesClause ] [ DerivesClause ] StructBody ;
SatisfiesClause = "satisfies" Formula ;
DerivesClause  = "derives" FormulaList ;
```

Notable consequences:

- `impl` is removed
- `implements` is not used for relation refinement, only in `fn` clauses when appropriate
- `contract` and `struct` remain explicit top-level surface forms
- `satisfies` names contract adoption directly
- `derives` names structural generation directly
- detailed contract and struct grammar is specified in the companion documents

This proposal does require `contract` and `struct` to be represented explicitly in the parser and validator as top keywords.

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

### Contract as Grouped Contract

The detailed contract surface and lowering model are specified in `contract.md`.

At the overview level, the important point is that `contract` remains an explicit surface declaration while its meaning is still reducible to ordinary `rel`, `fn`, and `spec` declarations.

### Contract Adoption and Structural Derivation

Surface:

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)

struct Point satisfies Eq(Point)
    x: Real
    y: Real
```

Lowering sketch:

```len
type Point

rel Point_x(p: Point, value: Real)
rel Point_y(p: Point, value: Real)

rel Eq(Point)
rel Equal_Point(x: Point, y: Point)
```

If the declaration instead uses `derives Eq(Point)`, elaboration may synthesize the canonical defining specs for `Equal_Point` from the struct fields.

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

### Example: Canonicalization Contract

This proposal keeps relation refinement and function implementation distinct even when they talk about the same domain.

```len
rel CanonicalPath(raw: Path, normalized: Path)
rel SlashNormalizedPath(raw: Path, normalized: Path) refines CanonicalPath(raw, normalized)

fn normalize_slashes(raw: Path) -> normalized: Path
    implements SlashNormalizedPath(raw, normalized)
    quasi:
        let normalized := raw
        ...
        return normalized
```

Here the roles stay separate:

- `refines` states that `SlashNormalizedPath` semantically implies `CanonicalPath`
- `implements` states that `normalize_slashes` is the executable form linked to `SlashNormalizedPath`
- the example does not add an `ensures CanonicalPath(...)` clause because that would be redundant once the implemented relation already refines `CanonicalPath`

### Example: Generic Contract Satisfaction

```len
contract Ord(T: Type)
    rel LessEq(x: T, y: T)
    spec ord_reflexive
        given x: T
        must LessEq(x, x)
    spec ord_transitive
        given x, y, z: T
        must LessEq(x, y) and LessEq(y, z) implies LessEq(x, z)

struct NatBox satisfies Ord(NatBox)
    value: Nat
```

This keeps the carrier type explicit. `Ord(NatBox)` is a contract proposition, not a hidden trait-instance table.

## Rationale

This proposal preserves the existing `len` philosophy.

`type` and `rel` remain the core modeling tools. `spec` remains the general mechanism for laws, definitions, and semantic constraints. `fn` remains the executable or proof-oriented form, and may still include an embedded `quasi` section when a structured but open-ended procedure sketch is useful.

Within that architecture:

- `refines` gives the right name for semantic implication between relations
- `contract` gives the right name for grouped contract declarations and their satisfaction relation
- `satisfies` gives the right name for explicit contract adoption
- `derives` provides a companion mechanism for generating derived specs from struct fields and contract members
- removing `impl` reduces misleading programming-language baggage
- replacing `trait` with `contract` removes unnecessary OO and Rust-adjacent vocabulary
- `struct` remains available as readable data-shape sugar without becoming an OO class

This keeps the language minimal while still making common patterns concise.

## Backward-Compatibility Notes

This proposal changes the preferred surface language and therefore requires migration for currently observed syntax.

- `impl` declarations should be rejected or deprecated in L1 once replacement lowering is implemented
- relation-level semantics previously described informally through `implements` should move to `refines` or explicit `spec`
- existing `trait` declarations should migrate to `contract`
- existing type-to-contract intent should migrate to `satisfies` or `derives`, depending on whether the contract view is declared or generated
- `fn ... implements ...` clauses remain as is
- existing parsers and validators that recognize `trait` and `impl` will need a transition strategy once the surface syntax is updated

Because current usage is still limited, this is a good point to simplify the vocabulary before those forms become entrenched.

## Open Questions

1. Should `fn` keep `implements` as the dedicated relation-linking clause, or should `ensures Relation(...)` also be allowed as equivalent surface sugar? The current direction keeps both, since they do different things: `implements` links the function to a relation as its contract, while `ensures` states a postcondition that may or may not be the same as the contract relation.
2. Should `refines` accept only relation applications, or any formula whose free variables match the relation signature?
3. Should `satisfies` accept only contract applications, or any formula recognized as a contract proposition?
4. What transition strategy should parsers and validators use while `trait`, `impl`, and older usage still exists in the repository?

## Conclusion

L1 should remain a high-level specification language centered on `type`, `rel`, `spec`, and `fn`, with `quasi` as an embedded auxiliary notation rather than a competing language tier.

Under that design, `contract` and `struct` are useful only if they stay as shallow surface sugar. Likewise, semantic implication between relations should be named `refines`, function-to-relation linkage should remain `implements`, and type-to-contract adoption should be named `satisfies` or `derives`.

This proposal therefore keeps the readable surface forms while removing the misleading semantics that normally travel with them in other languages.