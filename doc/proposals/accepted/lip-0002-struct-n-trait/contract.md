# LIP-0002 Contract Design Notes

## Purpose

This document specifies the `contract` portion of LIP-0002.

In `len.l1`, a `contract` is a first-class surface declaration for grouping a contract. It is interface-like in presentation, but it does not carry object-oriented runtime semantics.

A declaration of the form `contract Xyz(...) ...` introduces a same-named relation representing satisfaction of that contract. In that sense, `contract` is surface sugar for a contract relation together with the grouped member declarations and any related laws that elaboration may generate.

## Design Goal

`contract` should make contract-oriented modeling easier to read without changing the semantic core of the language.

The key constraints are:

- `contract` must remain surface-level structure, not a new semantic foundation below `type`, `rel`, and `spec`
- `contract` may group `rel`, `fn`, and `spec` declarations under one named contract
- `contract` should support interface-like examples when they improve readability
- `contract` must not imply inheritance, dispatch, hidden receivers, instance lookup, or Rust-style implementation rules

## Proposed Surface Role

A `contract` names a grouped contract.

Equivalently, `contract Eq(T: Type) ...` means that `Eq(T)` names a contract proposition with a grouped surface description.

The contract body groups ordinary `rel`, `fn`, and related law declarations under that contract-oriented surface form.

Conceptually, a `contract` may contain:

- required relations
- required functions
- laws expressed as ordinary `spec` declarations

The grouping is useful for organization, readability, and tooling. The semantics of the grouped members remain ordinary `rel`, `fn`, and `spec` semantics.

## Example

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)
    spec eq_reflexive
        given x: T
        must Equal(x, x)
```

This is intentionally interface-like. A reader from an OO background can recognize the idea of a reusable equality contract.

The crucial difference is semantic: `Eq(T)` is not a dispatch table or typeclass instance slot. It is a contract proposition saying that the type `T` comes equipped with the required relation and laws.

## Grammar Sketch

```ebnf
ContractDecl   = "contract" Identifier [ "(" [ ParamList ] ")" ] ContractBody ;
ContractBody   = INDENT { ContractMember } DEDENT ;
ContractMember = RelDecl
               | FnDecl
               | SpecDecl ;
```

This proposal treats `contract` as an explicit top-level surface declaration that the parser and validator should preserve.

## Lowering Model

Surface:

```len
contract Figure(T: Type)
    rel Area(self: T, value: Real)
    fn area(self: T) -> value: Real
        ensures Area(self, value)
```

Lowering sketch:

- a same-named contract relation `Figure(T)`
- lowered member declarations equivalent to ordinary `rel` and `fn`
- any contract laws lowered to ordinary `spec` declarations

One acceptable strategy is elaboration into prefixed or namespaced declarations around that contract relation. Another is to retain a contract node in the AST for tooling and validation, then elaborate later. The important requirement is that `contract` remains explicit in the surface model even if later compiler stages elaborate it.

## Satisfaction

Contracts are adopted with `satisfies`.

Example:

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)

struct Point satisfies Eq(Point)
    x: Real
    y: Real
```

This should be read as a declaration that `Point` is intended to satisfy the contract proposition `Eq(Point)`. It is not an `impl` block, and it does not create hidden instance machinery.

## Validation Expectations

The validator should understand `contract` as a real surface construct.

That means validation may need to check:

- member declaration shape
- duplicate member names within the contract
- parameter well-formedness
- whether embedded `fn` and `spec` declarations are well-formed in contract context
- whether later elaboration has enough information to generate a consistent lowered form

## Non-Goals

`contract` must not introduce:

- object-oriented inheritance semantics
- virtual dispatch or method lookup rules
- implicit receiver binding beyond explicit parameters such as `self`
- Rust-style `impl` blocks or trait-instance semantics

## Open Questions

1. Should `contract` permit nested `fn` declarations with `quasi`, or should contracts be restricted to signatures and laws only? it does permit nesting of `fn` with `quasi`.
2. Should `contract` introduce a namespace boundary, or should it elaborate into prefixed top-level names? namespace boundary. Instantiated contracts act as namespace containers, with the detailed import and alias rules specified in `LIP-0003`.
3. Should contract-local names be addressable directly in source-level tooling, or only through their elaborated form? they should be addressable through the contract namespace within the defining module and from other modules only via explicit import of the desired member.