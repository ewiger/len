The missing bridge is real. Right now:

* `struct` gives **data shape**
* `trait` gives **contract shape**

What you still need is a way to say:

* this struct **admits** a contract view
* some parts of that view may be **generated mechanically**
* the rest are still ordinary `rel`, `fn`, and `spec`

The Rust word `derive` points in the right direction, but in len it should not mean “macro magic that implements runtime traits.” It should mean something more like:

> from this declared structure, elaborate the obvious relational contract artifacts and obligations.

## The key design move

Do **not** connect `struct` and `trait` by hidden “implementation”.
Connect them by **elaboration into relations and specs**.

So the bridge should mean one of these two things:

### 1. Contract adoption

A struct declares that it is intended to satisfy a trait-shaped contract.

Example surface idea:

```len
struct Point derives Eq, Show
    x: Real
    y: Real
```

Meaning:

* `Point` is a type
* `Eq` and `Show` are contract groupings
* elaboration generates the relevant relations/functions/spec obligations for `Point`

This is the closest to Rust `derive`, but semantically cleaner.

### 2. Contract generation from structure

A trait can be marked as one that can be synthesized from structural information.

Example idea:

```len
trait Eq
    rel Equal(self, other, value: Bool)

struct Point derives Eq
    x: Real
    y: Real
```

Meaning:

* derive `Eq` by structural equality over the fields of `Point`
* the elaboration inserts the corresponding relation/function/specs

That gives you a genuine bridge.

---

# The most len-like interpretation of `derive`

In len, `derive` should not be primitive semantics. It should be **surface sugar for generated declarations**.

So:

```len
struct Point derives Eq
    x: Real
    y: Real
```

could lower roughly into:

```len
type Point

rel Point_x(p: Point, value: Real)
rel Point_y(p: Point, value: Real)

rel Eq_Point(a: Point, b: Point)

spec eq_point_def
    given a, b: Point
    must Eq_Point(a, b) iff
        exists ax, ay, bx, by: Real.
            Point_x(a, ax) and Point_y(a, ay) and
            Point_x(b, bx) and Point_y(b, by) and
            ax = bx and ay = by
```

Or, if trait members are named contracts:

```len
rel Equal_Point(a: Point, b: Point)

spec equal_point_def
    ...
```

The exact encoding can vary, but the meaning is:

* derive = generate the canonical contract relation/laws from the struct fields

That feels much more like len.

---

# What not to do

I would avoid any of these in L1:

* hidden nominal instance registry
* method lookup
* dispatch semantics
* “impl blocks”
* automatic conformance search
* inheritance between structs and traits

That would drag the language back toward Rust/OO instead of keeping it relational.

---

# Best conceptual model

I think the cleanest model is a **three-layer story**:

## Layer 1: shape

```len
struct Point
    x: Real
    y: Real
```

This introduces a type and field structure.

## Layer 2: contract

```len
trait Eq
    rel Equal(self, other)
```

This introduces a grouped contract vocabulary.

## Layer 3: elaboration bridge

```len
struct Point derives Eq
    x: Real
    y: Real
```

This says:

* instantiate the trait schema for `Point`
* generate the obvious relation/function/spec declarations
* attach generated laws based on the struct fields

So `derive` is not “Point implements Eq”.
It is:

> elaborate the Eq contract canonically for the struct Point.

That is much better for len.
