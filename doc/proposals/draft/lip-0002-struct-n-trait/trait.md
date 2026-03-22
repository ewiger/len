# LIP-0002 Trait Design Notes

## Purpose

This document specifies the `trait` portion of LIP-0002.

In `len.l1`, a `trait` is a first-class surface declaration for grouping a contract. It is interface-like in presentation, but it does not carry object-oriented runtime semantics.

A declaration of the form `trait Xyz ...` introduces a same-named relation representing satisfaction of that contract. In that sense, `trait` is surface sugar for a contract relation together with the grouped member declarations and any related laws that elaboration may generate.

## Design Goal

`trait` should make contract-oriented modeling easier to read without changing the semantic core of the language.

The key constraints are:

- `trait` must remain surface-level structure, not a new semantic foundation below `type`, `rel`, and `spec`
- `trait` may group `rel`, `fn`, and `spec` declarations under one named contract
- `trait` should support interface-like examples familiar from OO vocabulary
- `trait` must not imply inheritance, dispatch, hidden receivers, or Rust-style implementation rules

## Proposed Surface Role

A `trait` names a grouped contract.

Equivalently, `trait Figure ...` means that `Figure` names a contract relation with an interface-like surface description.

The trait body groups ordinary `rel`, `fn`, and related law declarations under that contract-oriented surface form.

Conceptually, a `trait` may contain:

- required relations
- required functions
- laws expressed as ordinary `spec` declarations

The grouping is useful for organization, readability, and tooling. The semantics of the grouped members remain ordinary `rel`, `fn`, and `spec` semantics.

## Example

```len
trait Figure
    rel Area(self: Figure, value: Real)
    fn area(self: Figure) -> value: Real
        ensures Area(self, value)
```

This is intentionally interface-like. A reader from an OO background can recognize the idea of a shape with an `area` operation.

The crucial difference is semantic: `area` is not a dispatch rule or a method slot. The actual contract is still carried by `Area(self, value)`.

## Grammar Sketch

```ebnf
TraitDecl      = "trait" Identifier TraitBody ;
TraitBody      = INDENT { TraitMember } DEDENT ;
TraitMember    = RelDecl
               | FnDecl
               | SpecDecl ;
```

This proposal treats `trait` as an explicit top-level surface declaration that the parser and validator should preserve.

## Lowering Model

Surface:

```len
trait Figure
    rel Area(self: Figure, value: Real)
    fn area(self: Figure) -> value: Real
        ensures Area(self, value)
```

Lowering sketch:

- a same-named contract relation `Figure(...)`
- lowered member declarations equivalent to ordinary `rel` and `fn`
- any trait laws lowered to ordinary `spec` declarations

One acceptable strategy is elaboration into prefixed or namespaced declarations around that contract relation. Another is to retain a trait node in the AST for tooling and validation, then elaborate later. The important requirement is that `trait` remains explicit in the surface model even if later compiler stages elaborate it.

## Validation Expectations

The validator should understand `trait` as a real surface construct.

That means validation may need to check:

- member declaration shape
- duplicate member names within the trait
- whether embedded `fn` and `spec` declarations are well-formed in trait context
- whether later elaboration has enough information to generate a consistent lowered form

## Non-Goals

`trait` must not introduce:

- object-oriented inheritance semantics
- virtual dispatch or method lookup rules
- implicit receiver binding beyond explicit parameters such as `self`
- Rust-style `impl` blocks or trait-instance semantics

## Open Questions

1. Should `trait` permit nested `fn` declarations with `quasi`, or should traits be restricted to signatures and laws only?
2. Should `trait` introduce a namespace boundary, or should it elaborate into prefixed top-level names?
3. Should trait-local names be addressable directly in source-level tooling, or only through their elaborated form?