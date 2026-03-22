# LIP-0002 Satisfaction and Derivation Notes

## Purpose

This document specifies the bridge between `struct` and `contract` in LIP-0002.

In `len.l1`, the bridge should be elaboration into ordinary relations and specs rather than hidden implementation machinery.

## Design Goal

`struct` gives data shape. `contract` gives contract shape. The language still needs a way to say:

- this struct adopts a contract view
- some parts of that view may be generated mechanically
- the rest remain ordinary `rel`, `fn`, and `spec`

The right split is:

- `satisfies` for explicit contract adoption
- `derives` for canonical structural generation

## The Key Design Move

Do not connect `struct` and `contract` by hidden implementation.

Connect them by elaboration into relations and specs.

That means the bridge has two distinct roles.

### 1. Contract Adoption with `satisfies`

A struct may declare that it adopts a contract explicitly.

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)

struct Point satisfies Eq(Point)
    x: Real
    y: Real
```

Meaning:

- `Point` is a type
- `Eq(Point)` is the intended contract proposition
- elaboration or validation must account for the corresponding relations, functions, and laws

This is the direct, explicit form. It names the contract view without implying any hidden registry, instance search, or dispatch rule.

### 2. Contract Generation with `derives`

Use `derives` when the contract view can be synthesized canonically from structure.

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)

struct Point derives Eq(Point)
    x: Real
    y: Real
```

Meaning:

- `Point` adopts `Eq(Point)`
- the adoption is structural rather than fully handwritten
- elaboration may generate the corresponding relation and defining specs from the fields of `Point`

`derives` is therefore not a replacement for `satisfies`. It is a specialized form of contract adoption that authorizes canonical generation.

## Lowering Model

In len, `derives` should not be primitive semantics. It should be surface sugar for generated declarations.

So:

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)

struct Point derives Eq(Point)
    x: Real
    y: Real
```

could lower roughly into:

```len
type Point

rel Point_x(p: Point, value: Real)
rel Point_y(p: Point, value: Real)

rel Eq(Point)
rel Equal_Point(a: Point, b: Point)

spec equal_point_def
    given a, b: Point
    must Equal_Point(a, b) iff
        exists ax, ay, bx, by: Real.
            Point_x(a, ax) and Point_y(a, ay) and
            Point_x(b, bx) and Point_y(b, by) and
            ax = bx and ay = by
```

The exact encoding can vary, but the meaning is stable:

- `satisfies` declares contract adoption
- `derives` generates the canonical contract relation and laws from the struct fields

## What Not to Do

L1 should avoid any of these:

- hidden nominal instance registries
- method lookup
- dispatch semantics
- `impl` blocks
- automatic conformance search
- inheritance between structs and contracts

Those would drag the language toward Rust or OO semantics instead of keeping it relational.

## Best Conceptual Model

The cleanest model is a three-layer story.

### Layer 1: Shape

```len
struct Point
    x: Real
    y: Real
```

This introduces a type and field structure.

### Layer 2: Contract

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)
```

This introduces a grouped contract vocabulary.

### Layer 3: Bridge

```len
struct Point derives Eq(Point)
    x: Real
    y: Real
```

This says:

- instantiate the contract schema for `Point`
- generate the obvious relation and spec declarations
- attach generated laws based on the struct fields

So `derives` does not mean `Point implements Eq`.

It means: elaborate the `Eq(Point)` contract canonically for the struct `Point`.
