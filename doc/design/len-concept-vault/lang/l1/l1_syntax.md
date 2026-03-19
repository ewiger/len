# `l1` Syntax (refined, minimal, relational)

`l1` is the first *formal* layer of `len`. It is designed to be:

- **relational-first** (facts and relations, not pseudocode)
- **contract/law-first** (truth statements are explicit and central)
- **small** (few keywords, few special forms)
- **strict** (no “empty words”; `l0` is where natural language tolerance lives)

---

## 1) Core ideas (semantic model)

`l1` keeps three roles distinct:

1. **Ontology / vocabulary**
   - `type` declares carrier types (names for sorts of values)
   - `rel` declares relations (predicates); relations are the semantic primitive

2. **Definition / derived concept**
   - `def` introduces *derived* predicates or *collection-like* concepts from formulas  
   - “what this name means”

3. **Truth / law**
   - `spec` states a *law/contract* that must hold (premises ⇒ obligations)
   - “what must be true”

**Key separation:**

- `def` describes *new concepts* (derived predicates / collection-like concepts).
- `spec` describes *truths about concepts* (laws/contracts that must hold).

---

## 2) Minimal keyword set

A small, teachable kernel:

- `import`
- `type`
- `rel`
- `def`
- `spec`
- `given`
- `must`
- `quasi` *(informal guidance only; no formal semantics / no formal logical reasoning)*

Everything else should be library-level or future-layer sugar.

---

## 3) Files, folders, and imports

- **Folder = module boundary.**
- **File = contribution to the folder scope.**
- There is no `module` keyword in `l1`.

Imports are explicit:

```
import foo/bar
import math/set as Set
```

(Exact import resolution rules are handled by tooling; `l1` keeps the surface small.)

---

## 4) Comments and docstrings

- Docstrings are allowed anywhere and **ignored by formal semantics**.
- Prefer docstrings over line comments for structured documentation.

```
"single-line docstring"

"""
Multi-line docstring.
Ignored by formal semantics.
Useful for intent, examples, and notes.
"""
```

---

## 5) Declarations

### 5.1 `type`

Carrier declaration:

```
type Expr
type Int
type Op
```

### 5.2 `rel`

Relation declaration (predicate signature). `rel` is the only semantic primitive besides `type`.

```
rel Const(e: Expr, v: Int)
rel Binary(e: Expr, op: Op, l: Expr, r: Expr)
rel Eval(e: Expr, v: Int)
```

Notes:

- A unary `rel P(x: T)` is set/class-like.
- An n-ary `rel R(...)` is a general predicate.
- If you want “functions”, express them as relations plus laws (e.g. uniqueness/totality) in `spec`,
  or as derived concepts in `def`.

Example (function as relation + laws):

```
type A
type B
type F

rel Fun(A: type, B: type, f: F)
rel Applies(f: F, x: A, y: B)

spec total_functionality(f: F, A: type, B: type)
given Fun(A, B, f)
must forall x: A. exists y: B. Applies(f, x, y)
must forall x: A. forall y1: B. forall y2: B.
    (Applies(f, x, y1) and Applies(f, x, y2)) -> y1 = y2
```

(Here `total + single-valued` gives you the usual extensional “function” behavior without introducing `fn` as a semantic primitive.)

---

## 6) Formulas (lightweight FOL)

`l1` formulas are first-order-logic style.

### 6.0 Arrow tokens: implication vs mapping

`l1` uses arrows in two different “worlds”, and they should not be confused:

- **Logical implication** inside formulas uses `=>` (and `<=>` for iff).
- **Mapping / function type notation** (e.g. `A -> B`) is *not* part of the `l1` core syntax.
  If you need to talk about functions/mappings in `l1`, model them relationally (see examples below).

About `=>`:

- `=>` is **not** a core token in `l1`.
- If you want an alternate implication glyph for readability, keep it as *pure surface sugar* for `->` and normalize it away immediately.
  Do not give `=>` a second meaning (like “maps to”), otherwise the language will drift.

Recommended policy:

- Reserve `=>` and `<->` for logic.
- For “maps to” in docs/examples, use words (`mapsTo(f, x, y)` as a relation) or a named relation like `Applies(f, x, y)`.

### 6.1 Connectives & quantifiers

- `and`, `or`, `not`
- `forall x: T. <formula>`
- `exists x: T. <formula>`
- `->` (implies), `<->` (iff)

### 6.2 Equality

- `=` and `!=`

### 6.3 No “empty words”

Words like `a`, `the`, `is`, `are`, `of`, `to`, `each` are **not** part of the formal grammar.
Use explicit logical structure instead. This keeps parsing and normalization crisp.

---

## 7) `def` — definitions (derived concepts)

`def` introduces a *meaning* by expansion to a formula. Use it for aliases and derived predicates.

Two common patterns are supported conceptually:

### 7.0 Examples inspired by “shaving the yak”

A minimal “function-stack” vocabulary (purely relational core):

```
type A
type B
type F

rel Fun(A: type, B: type, f: F)
rel total(f: F)
rel injective(f: F)
rel surjective(f: F)

def bijective(f: F) as injective(f) and surjective(f)

def TotalFun(A, B, f) as Fun(A, B, f) and total(f)
def BijectiveFun(A, B, f) as Fun(A, B, f) and bijective(f)
```

Collection-like “Hom-set” style (still just `def`):

```
def Hom(A, B) as f where TotalFun(A, B, f)
```

This stays aligned with the recommendation: keep the semantic core as `type`, `rel`, `def`, and express structure via relations + definitions.

### 7.1 Predicate alias (most common)

```
def bijective(f) as injective(f) and surjective(f)

def TotalFun(A, B, f) as Fun(A, B, f) and total(f)
```

### 7.2 “Collection-like” definition (set-builder view)

If you want a type/class-like reading, use a binder in the definition:

```
def TotalFuns(A, B) as f where
    Fun(A, B, f) and total(f)
```

This is a *derived concept* whose intended reading is “the collection of all `f` such that …”.
(Exactly how “collections” are represented internally is a normalization detail; the surface syntax stays simple.)

---

---

## 8) `spec` — laws / contracts (truth statements)

`spec` introduces a *truth obligation*. It is the unit of reasoning and verification.

Structure:

- `spec <name>(...)`
- zero or more `given <formula>` premises
- one or more `must <formula>` obligations

Interpretation:

- The spec asserts: **(given₁ and … and givenₙ) -> (must₁ and … and mustₖ)**

Example:

```
spec union_contains_left(A: Set, B: Set, U: Set)
given union(A, B, U)
must forall x. has(A, x) -> has(U, x)
```

Notes:

- Use multiple `must` lines for readability; they are conjunctive.
- `spec` is for laws/theorems/contracts/policies—statements that must hold.

---

## 9) `quasi` — open-ended implementation/proof guidance (walled off)

`quasi` is an escape hatch for proof hints, implementation steering, or constructive sketches.
It has **no formal semantics** and is **not part of formal logical reasoning**.

```
spec eval_const(e: Expr, v: Int)
given Const(e, v)
must Eval(e, v)

quasi
    "Implementation/proof sketch goes here."
```

Rule of thumb:

- If it matters semantically, it must be in `given`/`must`/`def`/`spec`.
- If it’s guidance, it belongs under `quasi`.

---

## 10) Style constraints (to keep `l1` crisp)

- Prefer `rel` over `fn`. (If `fn` exists at all, treat it as surface sugar for a relation plus laws.)
- Avoid heavy built-in set/collection syntax in the core. Prefer relations (`has`, `member`, etc.)
  and define conveniences in a prelude library.
- Reserve punctuation operators consistently; don’t allow redefining core logical/equality operators.

---

## 11) End-to-end example (calculator-style)

```
type Expr
type Op
type Int

rel Const(e: Expr, v: Int)
rel Binary(e: Expr, op: Op, l: Expr, r: Expr)

rel AddOp(op: Op)
rel DivOp(op: Op)

rel Eval(e: Expr, v: Int)
rel add(a: Int, b: Int, out: Int)

spec eval_const(e: Expr, v: Int)
given Const(e, v)
must Eval(e, v)

spec eval_add(e: Expr, op: Op, l: Expr, r: Expr, v: Int)
given Binary(e, op, l, r)
given AddOp(op)
given exists a: Int. Eval(l, a)
given exists b: Int. Eval(r, b)
must exists a: Int. exists b: Int.
    Eval(l, a) and Eval(r, b) and add(a, b, v)
must Eval(e, v)

spec no_div_by_zero(e: Expr, op: Op, l: Expr, r: Expr)
given Binary(e, op, l, r)
given DivOp(op)
must exists b: Int. Eval(r, b) and b != 0
```

This demonstrates the intended feel:

- vocabulary via `type`/`rel`
- truths via `spec` (`given`/`must`)
- no imperative pseudocode in the core
- constructive or proof guidance (if needed) goes into `quasi`

---
