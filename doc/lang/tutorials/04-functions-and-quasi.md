# Tutorial 4 — Functions and Quasi

A `fn` declaration is the executable or constructive side of len. Every function is always linked to a relation through `implements` and must carry a `quasi` block — a structured pseudocode sketch describing how it works.

---

## The `fn` Form

```len
fn <name>(<params>) -> <result>: <Type>
    implements <Relation>(<args>)
    requires <formula>        # optional, one or more
    ensures  <formula>        # optional, one or more
    quasi [using style <Name>]:
        <pseudocode body>
```

All clauses except `quasi` are optional, but `quasi` is required on every `fn`.

---

## A Simple Function

```len
type Seq

rel Reverse(input: Seq, output: Seq)

fn reverse(input: Seq) -> output: Seq
    implements Reverse(input, output)
    ensures SameLength(input, output)
    quasi using style ProceduralAlgorithm:
        let result := empty_seq()
        for item in input:
            set result := prepend(item, result)
        return result
```

- `implements Reverse(input, output)` — this function is the realisation of the `Reverse` relation.
- `ensures SameLength(input, output)` — the validator records this as a contractual guarantee.
- The `quasi` block is the procedural sketch.

---

## Requires and Ensures

`requires` states a precondition — a formula that must hold **before** the function is called.

`ensures` states a postcondition — a formula that must hold **after** the function returns.

Both are formal contract clauses: they do not get executed, but they are recorded and can be verified by tools that work at higher levels (l2 and beyond).

```len
type Nat

rel Divide(a: Nat, b: Nat, result: Nat)
rel NonZero(n: Nat)

fn divide(a: Nat, b: Nat) -> result: Nat
    implements Divide(a, b, result)
    requires NonZero(b)
    quasi using style ProceduralAlgorithm:
        return a / b
```

Multiple ensures clauses accumulate — all must hold:

```len
fn sort(input: Seq) -> output: Seq
    implements Sort(input, output)
    ensures Sorted(output)
    ensures Permutation(input, output)
    quasi:
        # ...
```

---

## The Quasi Block

### What it is

`quasi` is an embedded pseudocode sketch. It lives inside `fn` and captures *how* the algorithm proceeds. It is not a real programming language — it is validated at the surface level for structural well-formedness.

### Style profiles

A style profile is a YAML file that defines what surface forms are allowed inside a `quasi` block.

The current built-in style is **`ProceduralAlgorithm`**. To use it explicitly:

```len
quasi using style ProceduralAlgorithm:
    ...
```

If you omit `using style`, the default profile (`ProceduralAlgorithm`) is used:

```len
quasi:
    ...
```

### ProceduralAlgorithm syntax

| Form | Example | Meaning |
|------|---------|---------|
| `let x := expr` | `let n := length(xs)` | Introduce a variable |
| `set x := expr` | `set n := n - 1` | Reassign a variable |
| `return expr` | `return result` | Return a value |
| `if cond:` | `if n > 0:` | Conditional |
| `else:` | `else:` | Else branch |
| `while cond:` | `while changed:` | Loop |
| `for x in range:` | `for i in 0 .. n - 1:` | Iteration |
| `append(list, val)` | `append(result, item)` | Grow a sequence |

These keywords (`let`, `set`, `return`, `if`, `else`, `while`, `for`) are **quasi-local** — they are not part of the global len grammar.

---

## Full Example: Sorting Algorithms

Sorting is a running example in the len corpus. Here is the shared vocabulary:

```len
# sorting.l1
type Int
type Elem
type Seq

rel Index(s: Seq, i: Int)
rel At(s: Seq, i: Int, x: Elem)
rel Length(s: Seq, n: Int)
rel Leq(x: Elem, y: Elem)
rel Sorted(s: Seq)
rel Permutation(a: Seq, b: Seq)
rel Sort(input: Seq, output: Seq)

spec sort_def
    given input, output: Seq
    must Sort(input, output) iff
        Sorted(output) and Permutation(input, output)
```

Bubble sort:

```len
# bubble_sort.l1
import sorting

rel BubbleSort(input: Seq, output: Seq)

spec bubble_sort_correct
    given input: Seq
    given output: Seq
    must BubbleSort(input, output) implies
        Sorted(output) and Permutation(input, output)

fn bubble_sort(input: Seq) -> output: Seq
    implements BubbleSort(input, output)
    ensures Sorted(output)
    ensures Permutation(input, output)
    quasi using style ProceduralAlgorithm:
        let xs := input
        let n := length(xs)
        let changed := true

        while changed:
            set changed := false

            for i in 0 .. n - 2:
                if xs[i] > xs[i + 1]:
                    let tmp := xs[i]
                    set xs[i] := xs[i + 1]
                    set xs[i + 1] := tmp
                    set changed := true

        return xs
```

Insertion sort (contrast: it uses a different quasi structure):

```len
# insertion_sort.l1
import sorting

rel InsertionSort(input: Seq, output: Seq)

fn insertion_sort(input: Seq) -> output: Seq
    implements InsertionSort(input, output)
    ensures Sorted(output)
    ensures Permutation(input, output)
    quasi using style ProceduralAlgorithm:
        let result := empty_seq()

        for item in input:
            let i := length(result)
            while i > 0 and result[i - 1] > item:
                set i := i - 1
            set result := insert_at(result, i, item)

        return result
```

---

## The `fn` + `rel` + `spec` Pattern

The idiomatic len pattern for any executable behaviour is:

1. **Declare the relation** (`rel`) — what the function means abstractly.
2. **Constrain it with specs** (`spec`) — what laws it must obey.
3. **Implement it with a function** (`fn`) — how it is done, linked back to the relation.

```len
# Step 1 — relation
rel FindMax(s: Seq, x: Elem)

# Step 2 — laws
spec find_max_def
    given s: Seq, x: Elem
    must FindMax(s, x) iff
        x in s and forall y: Elem . y in s implies Leq(y, x)

# Step 3 — function
fn find_max(s: Seq) -> x: Elem
    implements FindMax(s, x)
    requires NonEmpty(s)
    ensures forall y: Elem . y in s implies Leq(y, x)
    quasi using style ProceduralAlgorithm:
        let best := first(s)
        for item in rest(s):
            if item > best:
                set best := item
        return best
```

---

## What You Learned

- Every `fn` must have a `quasi` block — it is not optional
- `implements` links a function to the relation it realises
- `requires` and `ensures` add preconditions and postconditions
- The `ProceduralAlgorithm` style allows: `let`, `set`, `return`, `if`, `else`, `while`, `for`
- The idiomatic pattern is: `rel` → `spec` → `fn`

**Next:** [Contracts and Structs →](05-contracts-and-structs.md)
