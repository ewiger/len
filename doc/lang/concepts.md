# Concepts

This page explains the core ideas behind `len` and why it is designed the way it is.

---

## What is len?

`len` is a **meta-programming language**. Its programs do not run applications — they describe, generate, or reason about other programs and code structures.

Typical uses:

- Describing a domain model precisely enough to generate type-safe code from it
- Writing formal contracts that a code generator must respect
- Specifying algorithms at an abstract level before committing to an implementation language
- Self-hosting: expressing the len parser and validator in len itself

---

## Specification First

In most languages you write code and infer the specification later (if at all). In `len` you write the specification first and derive the code from it.

A **relation** describes a property or relationship between values. A **spec** states a law about those relations. A **function** is an executable form that is always linked to a relation through `implements`. The specification is not optional documentation — it is the primary artefact.

```len
rel Sort(input: Seq, output: Seq)

spec sort_correct
    given input, output: Seq
    must Sort(input, output) iff
        Sorted(output) and Permutation(input, output)

fn sort(input: Seq) -> output: Seq
    implements Sort(input, output)
    ensures Sorted(output)
    quasi:
        # ... algorithm sketch
```

---

## Language Layers

len is deliberately layered. Each layer adds expressiveness while keeping the layer below it stable.

### l0 — Natural Language

Files: none (or `.l0` by convention)

l0 is informal. It uses prose, pseudocode, and diagrams. Its purpose is communication, design exploration, and onboarding.

```text
# l0 sketch
Sort takes a list and returns a new list that has the same elements
in non-decreasing order.
```

### l1 — Structural Core

Files: `.l1`

l1 is the canonical heart of len. It is intentionally **small and declarative**. The entire vocabulary fits on one page:

| Declaration | Purpose |
|-------------|---------|
| `type` | Introduce a named type |
| `rel` | Introduce a relation (property or relationship) |
| `const` | Introduce a named constant with a type |
| `fn` | Declare a function with a contract and quasi body |
| `spec` | State a law, axiom, definition, or theorem |
| `contract` | Group related relations and specs under a named contract |
| `struct` | Describe a composite record type |
| `syntax` | Define surface-level notation sugar |
| `import` | Load declarations from another module |
| `keyword` | Reserve a word in the surface language |
| `symbol` | Reserve a symbolic token |

l1 is **relational**, not object-oriented. There are no methods, classes, inheritance trees, or dispatch tables. Instead:

- Types are named sets of values
- Relations state propositions about values
- Specs constrain those relations with logical laws
- Functions are linked to relations through `implements`

### l2 — Evaluation

Files: `.l2` (future)

l2 adds interpretation: contexts, satisfaction checking, executable semantic rules. It builds on l1 without changing it.

---

## Module System

len organises code into modules identified by dot-separated paths:

```len
import core.math.set
import core.math.nat
```

A module path maps to a file path on disk. For example, `core.math.set` maps to `core/math/set.l1` (or `core/math/set/set.l1` by convention for directories).

Selective imports pull specific names from a module:

```len
from core.math.set import Set, Member, Subset
```

The standard library ships under `core`:

| Module | Contents |
|--------|----------|
| `core.types.types` | Primitive `Type` and `Relation` types, subtype laws |
| `core.syntax.syntax` | Reflective meta-level carriers: `Expr`, `Identifier`, `Equals` |
| `core.syntax.keywords` | Canonical reserved keyword declarations |
| `core.math.logic.logic` | `Bool`, `Formula`, `Eval`, `Holds` |
| `core.math.logic.sugar` | Operator sugar: `=>`, `<=>`, `/\`, `\/`, `!` |
| `core.math.set` | `Set`, `Member`, `Subset`, `Empty`, axioms of set theory |
| `core.math.set.nat` | `Nat`, `Zero`, `Successor`, Peano axioms |
| `core.math.fun` | `Fn`, `Injective`, `Surjective`, `Bijective` |
| `core.struct` | `Struct` contract: `Field` relation |

---

## Relation Refinement

One relation can refine another — that is, it is a more specific version of the same idea:

```len
rel Sort(input: Seq, output: Seq)
rel BubbleSort(input: Seq, output: Seq) refines Sort
```

`BubbleSort` inherits all the obligations of `Sort` plus any additional laws stated in its own specs.

---

## Contracts

A `contract` groups related declarations under a named proposition:

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)
    spec eq_reflexive
        given x: T
        must Equal(x, x)
    spec eq_symmetric
        given x, y: T
        must Equal(x, y) implies Equal(y, x)
```

`Eq(T)` is a proposition stating that type `T` satisfies the equality contract. Unlike a Java interface, there is no dispatch — satisfying a contract is a logical statement, not a runtime object.

---

## Quasi Blocks

Every `fn` declaration must include a `quasi` block: an embedded pseudocode sketch that describes *how* the function works. The quasi block is not interpreted by l1 — it is captured as raw text and validated against a **style profile**.

```len
fn bubble_sort(input: Seq) -> output: Seq
    implements BubbleSort(input, output)
    ensures Sorted(output)
    quasi using style ProceduralAlgorithm:
        let xs := input
        while needs_swap(xs):
            set xs := bubble_pass(xs)
        return xs
```

The default style is `ProceduralAlgorithm`. Custom styles can be declared and validated through YAML profile files.

---

## Design Principles

1. **Specification before implementation.** The `rel` + `spec` pair is the primary artefact. The `fn` + `quasi` pair is secondary.
2. **Small core, extensible surface.** The language reserves a minimal set of keywords. New notation is introduced explicitly through `syntax` declarations.
3. **No hidden semantics.** Syntax sugar is always declared with an explicit `implies` target. There is no magic.
4. **Layer discipline.** l1 is stable. l2 builds on it. Nothing in l1 depends on l2.
5. **Relational, not object-oriented.** Types, relations, and specs replace classes, methods, and interfaces.
