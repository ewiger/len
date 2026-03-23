# Keywords

Every reserved keyword in `len.l1`, grouped by role.

---

## First-Class Primitives

### `type`

Declares a named type.

```len
type Point
type Seq
type Nat
```

A type is an opaque name. Len does not infer structure from the name — you add structure through relations and specs. A type can carry parameters through a `contract`.

---

### `rel`

Declares a relation — a named proposition about one or more typed values.

```len
rel Member(x: Set, s: Set)
rel Sort(input: Seq, output: Seq)
rel Sorted(s: Seq)
```

Relations are the primary modelling tool in len. They replace methods, interfaces, and class properties from OOP languages.

A relation may optionally refine another:

```len
rel BubbleSort(input: Seq, output: Seq) refines Sort
```

---

### `const`

Declares a named constant with a type annotation.

```len
const True  : Bool
const False : Bool
```

Constants are values that are always the same. Unlike relations, they do not take arguments.

---

### `fn`

Declares a function — an executable or constructive form linked to a relation.

```len
fn sort(input: Seq) -> output: Seq
    implements Sort(input, output)
    ensures Sorted(output)
    quasi:
        # pseudocode sketch
```

Every `fn` declaration must contain a `quasi` block. The function signature may include:

- `implements` — the relation this function realises
- `requires` — preconditions that must hold before the call
- `ensures` — postconditions guaranteed after the call

---

### `contract`

Groups related `rel`, `fn`, and `spec` declarations under a named proposition.

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)
    spec eq_reflexive
        given x: T
        must Equal(x, x)
```

`contract` is interface-like in appearance but carries no runtime dispatch. `Eq(T)` is a logical proposition stating that type `T` satisfies the equality contract.

---

## Declarations and Grouping

### `spec`

Declares a logical statement — an axiom, definition, theorem, or law.

```len
spec sort_correct
    given input, output: Seq
    must Sort(input, output) iff
        Sorted(output) and Permutation(input, output)
```

A spec body uses two clauses:

| Clause | Meaning |
|--------|---------|
| `given <var>: <Type>` | Introduce a universally-bound variable |
| `must <formula>` | State the law that must hold |

By convention, use the following labels to indicate the role of a spec:

| Label | Role |
|-------|------|
| `axiom_*` | Fundamental assumption that is not derived |
| `definition_*` | Introduces a new notion in terms of existing ones |
| `theorem_*` | Derivable from other specs |
| `lemma_*` | Auxiliary statement used to prove a theorem |
| `corollary_*` | Follows directly from a theorem |
| `postulate_*` | Local primitive assumption |

---

### `struct`

Declares a composite record type — sugar over `type` and field relations.

```len
struct Product
    name  : String
    price : Nat
    sku   : String
```

`struct` is syntactic sugar. It expands to a `type` declaration plus field relations, so all the regular `spec` and `rel` machinery applies to a struct.

---

## Syntax and Symbols

### `syntax`

Declares a surface-level notation alias. Syntax declarations make formulas more readable without changing the semantic core.

```len
syntax x in s where x: Set, s: Set implies Member(x, s)
syntax A => B where A: Formula, B: Formula implies A implies B
```

Shape: `syntax <pattern> where <bindings> implies <target-relation>`.

The `where` clause introduces typed binders; `implies` names the relation the notation desugars to.

---

### `symbol`

Declares a symbolic token for use in syntax declarations.

```len
symbol =>
symbol /\
```

---

### `keyword`

Reserves a word so the lexer treats it as a keyword rather than an identifier.

```len
keyword type
keyword rel
```

This is used in `core.syntax.keywords` to declare the canonical list of reserved words. User code rarely needs `keyword` directly.

---

## Imports

### `import`

Loads all public declarations from another module.

```len
import core.math.set
import core.math.nat
```

Module paths are dot-separated. They map to file paths on disk (e.g., `core/math/set.l1`).

---

### `from`

Selective import — loads specific names from a module.

```len
from core.math.set import Set, Member
```

---

## Spec Clauses

These keywords appear inside `spec` and `fn` bodies rather than at the top level.

### `given`

Introduces a universally-bound variable inside a `spec` body.

```len
spec extensionality
    given a: Set
    given b: Set
    must (forall x: Set . x in a iff x in b) implies a = b
```

---

### `must`

States the formula that the spec asserts.

```len
spec pairing
    must forall a: Set, b: Set .
        exists p: Set .
            forall x: Set . x in p iff (x = a or x = b)
```

---

### `implements`

Inside an `fn` declaration, links the function to the relation it realises.

```len
fn sort(input: Seq) -> output: Seq
    implements Sort(input, output)
```

---

### `requires`

Inside an `fn` declaration, states a precondition.

```len
fn divide(a: Nat, b: Nat) -> result: Nat
    requires NonZero(b)
    implements Divide(a, b, result)
    quasi:
        return a / b
```

---

### `ensures`

Inside an `fn` declaration, states a postcondition.

```len
fn sort(input: Seq) -> output: Seq
    implements Sort(input, output)
    ensures Sorted(output)
    ensures Permutation(input, output)
    quasi:
        # ...
```

---

## Logical Connectives

### `forall`

Universal quantifier — asserts that a formula holds for every value of a type.

```len
must forall x: Set . x in x implies false
```

---

### `exists`

Existential quantifier — asserts that at least one value satisfies a formula.

```len
must exists s: Set . Empty(s)
```

---

### `and`

Logical conjunction. Both sides must hold.

```len
must Sorted(output) and Permutation(input, output)
```

Also available as `/\` through the sugar module.

---

### `or`

Logical disjunction. At least one side must hold.

```len
must Zero(a) or Positive(a)
```

Also available as `\/` through the sugar module.

---

### `not`

Logical negation.

```len
must not x in s
```

Also available as `!` through the sugar module.

---

### `iff`

Logical biconditional — both directions of implication.

```len
must Empty(s) iff forall x: Set . not x in s
```

Also available as `<=>` through the sugar module.

---

### `implies`

Logical implication — if the left side holds then the right side must hold.

```len
must Sorted(output) implies Permutation(input, output)
```

Also available as `=>` through the sugar module.

---

### `where`

Introduces inline binders in `syntax` declarations and acts as a scoping delimiter in formulas.

```len
syntax x in s where x: Set, s: Set implies Member(x, s)
```

---

### `in`

Membership operator in set-theory contexts (sugar for `Member`).

```len
must forall x: Set . x in a iff x in b
```

---

### `true` / `false`

Boolean constants. Map to `True` and `False` in `core.math.logic.logic`.

---

### `as`

Used for aliased imports (reserved; full form forthcoming).

```len
from core.math.set import Member as In
```

---

## Quasi

### `quasi`

Attaches a pseudocode body to an `fn` declaration. Required on every `fn`.

```len
fn bubble_sort(input: Seq) -> output: Seq
    implements BubbleSort(input, output)
    quasi using style ProceduralAlgorithm:
        let xs := input
        while needs_swap(xs):
            set xs := bubble_pass(xs)
        return xs
```

The `quasi` clause may optionally name a style profile: `quasi using style <Name>:`. If omitted, the default style (`ProceduralAlgorithm`) is used.

Keywords that are **quasi-local** (not global len keywords): `let`, `set`, `return`, `if`, `else`, `while`, `for`, `append`.
