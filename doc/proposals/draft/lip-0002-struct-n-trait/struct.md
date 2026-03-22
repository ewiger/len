# LIP-0002 Struct Design Notes

## Purpose

This document specifies the `struct` portion of LIP-0002.

In `len.l1`, a `struct` is a first-class surface declaration for a composite type with named fields. It is record-like, not class-like.

A declaration of the form `struct Xyz ...` introduces `Xyz` as a type. In that sense, `struct` is surface sugar for a `type` declaration together with the field-related relational structure that elaboration may generate.

## Design Goal

`struct` should make composite data shapes readable without changing the semantic core of the language.

The key constraints are:

- `struct` must remain surface-level structure, not a new semantic foundation below `type` and `rel`
- `struct` should behave like a record or named tuple
- field syntax should stay lightweight
- `struct` must not imply methods, encapsulation, mutation protocols, or object identity

## Proposed Surface Role

A `struct` declares a composite type with named fields.

Equivalently, `struct Point ...` means that `Point` is a type with a record-like surface description.

Example:

```len
struct Point
    x: Real
    y: Real
```

The intent is a readable surface way to say that `Point` is a type with named components. The form is useful for modeling and tooling, but it does not create OO class semantics.

## Grammar Sketch

```ebnf
StructDecl     = "struct" Identifier StructBody ;
StructBody     = INDENT { FieldDecl } DEDENT ;
FieldDecl      = Identifier ":" TypeExpr ;
```

This proposal treats `struct` as an explicit top-level surface declaration that the parser and validator should preserve.

## Lowering Model

Surface:

```len
struct Point
    x: Real
    y: Real
```

Lowering sketch:

```len
type Point

rel Point_x(p: Point, value: Real)
rel Point_y(p: Point, value: Real)

spec point_field_totality
    given p: Point
    must exists x: Real, y: Real. Point_x(p, x) and Point_y(p, y)
```

The generated names are illustrative only. The important design point is that `struct` elaborates into a `type` declaration plus existing type-and-relation machinery rather than introducing class behavior.

## Example

```len
struct Interval
    start: Real
    end: Real

spec interval_well_formed
    given i: Interval
    must exists s: Real, e: Real. Interval_start(i, s) and Interval_end(i, e) and s <= e
```

This shows the intended division of labor:

- `struct` introduces the composite shape
- `spec` carries additional laws about well-formed instances

## Relation to Traits

See `derive.md` for discussion of how `struct` and `trait` can interact.

## Validation Expectations

The validator should understand `struct` as a real surface construct.

That means validation may need to check:

- duplicate field names
- field type well-formedness
- whether later elaboration can generate a consistent lowered representation
- consistency between a struct declaration and any related field or law declarations

## Non-Goals

`struct` must not introduce:

- class inheritance
- method dispatch semantics
- hidden constructors with special runtime behavior
- mutation or ownership rules

## Open Questions

1. What is the canonical lowered representation for `struct` fields: projection relations, constructor relations, or another equivalent core encoding?
2. Should generated names from struct lowering be visible in source-level tooling, or treated as internal elaboration artifacts?
3. Should future syntax allow optional struct-local invariants, or should all such laws remain ordinary `spec` declarations outside the struct body?