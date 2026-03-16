Great direction. And yes: **IR = Intermediate Representation**.

Here it means:

- not raw human language
    
- not yet full Lean / SMT / prover syntax
    
- but a **structured semantic layer in between**
    

So `len` would be a **semantic intermediate language** for specifications and knowledge, designed to be:

- easier to produce from natural language than Lean
    
- easier to compile into logic than natural language
    
- easier to read back into human language than low-level proof terms
    

A good mental picture is:

[  
\text{HumanLanguage} \xrightarrow{S_1} \texttt{len} \xrightarrow{S_2} \text{Lean / FOL / SMT}  
]

and back:

[  
\text{Lean / FOL / SMT} \xrightarrow{I_1} \texttt{len} \xrightarrow{I_2} \text{HumanLanguage}  
]

So `len` is the **semantic bridge language**.

---

# 1. What `len` should be

From what you said, `len` should combine three styles:

1. **predicate style**  
    `Sorted(x)`, `Member(x, xs)`
    
2. **mathematical / set-theoretic style**  
    `x ∈ xs`, `∀ y ∈ xs : x ≤ y`
    
3. **high-level controlled natural language**  
    `result is in xs`  
    `for every y in xs, result <= y`
    

That suggests `len` should not be purely symbolic or purely natural.  
It should be a **controlled semantic language** with:

- typed entities
    
- predicates / relations
    
- quantifiers
    
- contracts
    
- definitions
    
- explicit ambiguity markers
    
- optionally both symbolic and readable surface forms
    

---

# 2. Design goal for `len`

A good `len` should satisfy these principles:

### P1. Semantic explicitness

No important assumption should stay implicit.

### P2. Structural simplicity

It should be much simpler than Lean.

### P3. Bidirectional readability

Humans should be able to read and edit it.

### P4. Deterministic compilation

It should compile to formal logic with minimal guesswork.

### P5. Relational by default

Specifications should describe relations, not only functions.

### P6. Partiality-friendly

It should allow underspecification, ambiguity, and open slots.

This is important: `len` should not force false precision too early.

---

# 3. Core idea: two-layer syntax

I would recommend `len` have two surfaces:

## A. Canonical symbolic form

Good for compiler and formal semantics.

## B. Sugared readable form

Good for humans and LLMs.

Example:

Canonical:

```text
spec min(xs: List<Int>) -> result: Int
requires NonEmpty(xs)
ensures Member(result, xs)
ensures Forall(y in xs, result <= y)
```

Sugared:

```text
spec min(xs: List<Int>) -> result: Int
requires xs is non-empty
ensures result in xs
ensures for every y in xs, result <= y
```

These should mean the same thing.

---

# 4. First draft syntax for `len`

Here is a compact core.

## 4.1 Declarations

```text
type Int
type Bool
type List<T>
type Set<T>
```

```text
const empty : List<Int>
const zero : Int
```

```text
fn min(xs: List<Int>) -> Int
rel Member(x: Int, xs: List<Int>)
rel Sorted(xs: List<Int>)
rel Unique(xs: List<Int>)
```

Where:

- `fn` means function symbol
    
- `rel` means predicate / relation
    
- `type` introduces types
    
- `const` introduces constants
    

---

## 4.2 Terms

```text
x
f(x)
head(xs)
length(xs)
result
```

---

## 4.3 Formulas

```text
P(x)
not P(x)
P(x) and Q(x)
P(x) or Q(x)
P(x) -> Q(x)
P(x) <-> Q(x)
x = y
x != y
x in xs
x <= y
```

---

## 4.4 Quantifiers

```text
forall x: Int . P(x)
exists x: Int . P(x)
forall x in xs . P(x)
exists x in xs . P(x)
```

---

## 4.5 Definitions

```text
define NonEmpty(xs: List<T>) :=
    exists x in xs . true
```

or

```text
define NonEmpty(xs) := length(xs) > 0
```

---

## 4.6 Specs

```text
spec min(xs: List<Int>) -> result: Int
requires NonEmpty(xs)
ensures result in xs
ensures forall y in xs . result <= y
```

---

## 4.7 Ambiguity / open slots

Very important.

```text
spec sort(xs: List<Int>) -> ys: List<Int>
ensures Sorted(ys)
ensures Permutation(ys, xs)
open Stability(ys, xs)
open OrderingRelation
```

This says: the spec is incomplete in these dimensions.

That is extremely useful.

---

# 5. Suggested minimal grammar

A rough EBNF-style sketch:

```text
Program     ::= Decl*
Decl        ::= TypeDecl | ConstDecl | FnDecl | RelDecl | DefDecl | SpecDecl

TypeDecl    ::= "type" IDENT
ConstDecl   ::= "const" IDENT ":" TypeExpr
FnDecl      ::= "fn" IDENT "(" ParamList? ")" "->" TypeExpr
RelDecl     ::= "rel" IDENT "(" ParamList? ")"
DefDecl     ::= "define" IDENT "(" ParamList? ")" ":=" Formula
SpecDecl    ::= "spec" IDENT "(" ParamList? ")" "->" ResultDecl Newline SpecBody

ParamList    ::= Param ("," Param)*
Param        ::= IDENT ":" TypeExpr
ResultDecl   ::= IDENT ":" TypeExpr
SpecBody     ::= Clause*
Clause       ::= "requires" Formula
              | "ensures" Formula
              | "invariant" Formula
              | "reads" TermList
              | "writes" TermList
              | "open" IDENT
              | "note" STRING

Formula      ::= Atomic
              | "not" Formula
              | Formula "and" Formula
              | Formula "or" Formula
              | Formula "->" Formula
              | Formula "<->" Formula
              | "forall" Binder "." Formula
              | "exists" Binder "." Formula

Binder       ::= IDENT ":" TypeExpr
              | IDENT "in" Term

Atomic       ::= IDENT "(" ArgList? ")"
              | Term RelOp Term
              | Term "in" Term

Term         ::= IDENT
              | NUMBER
              | IDENT "(" ArgList? ")"

RelOp        ::= "=" | "!=" | "<" | "<=" | ">" | ">="
```

This is already enough to express a surprising amount.

---

# 6. Semantics style: relational, not implementation-first

I strongly agree with making `len` **spec-first**.

Instead of saying:

```text
fn dedup_sort(xs) = ...
```

say:

```text
spec canonicalize(xs: List<Int>) -> ys: List<Int>
ensures Sorted(ys)
ensures Unique(ys)
ensures forall y in ys . y in xs
ensures forall x in xs . x in ys <-> exists i . xs[i] = x
```

This keeps `len` focused on meaning, not code.

---

# 7. Three syntax styles inside `len`

You mentioned liking three styles. I would preserve all three, but assign them roles.

## Style A: predicate/core form

Best for canonical semantics.

```text
Sorted(xs)
Unique(xs)
Member(x, xs)
```

## Style B: mathematical sugar

Best for readability.

```text
x in xs
forall y in xs . result <= y
```

## Style C: controlled natural form

Best for interface and round-trip checks.

```text
result is in xs
for every y in xs, result <= y
```

Recommendation: define **Style A as canonical**, and let B/C desugar into A.

Example:

```text
x in xs
```

desugars to

```text
Member(x, xs)
```

and

```text
result is in xs
```

also desugars to

```text
Member(result, xs)
```

That keeps semantics clean.

---

# 8. Properties we should want from `len`

Now for your part (b): state properties.

## 8.1 Parseability

Every valid `len` program should have an unambiguous parse tree.

## 8.2 Surface equivalence

Different surface forms can map to the same canonical form.

Example:

```text
result in xs
Member(result, xs)
result is in xs
```

all normalize to the same internal representation.

## 8.3 Normalization

There should be a normalization function

[  
N : \texttt{len} \to \texttt{len}_{core}  
]

so semantic comparison happens on normalized forms.

## 8.4 Partiality support

`len` should represent incomplete specs without inventing details.

## 8.5 Verbalizability

Every normalized `len` fragment should have a readable interpolation into controlled language.

## 8.6 Logic-compilability

There should be a deterministic compiler from normalized `len` to target KR:

[  
C_{\text{Lean}} : \texttt{len}_{core} \to \text{Lean}  
]

[  
C_{\text{FOL}} : \texttt{len}_{core} \to \text{FOL}  
]

## 8.7 Semantic preservation

Compilation should preserve meaning.

Informally:

[  
\llbracket e \rrbracket_{\texttt{len}} = \llbracket C(e) \rrbracket_{\text{KR}}  
]

## 8.8 Round-trip stability

For some interpolation map (I),

[  
N(S(I(e))) = N(e)  
]

at least for a stable subset of `len`.

This is a very important property.

---

# 9. Suggested semantic layers of `len`

I would define `len` in three levels.

## len-0: Core logic/spec layer

- types
    
- predicates
    
- terms
    
- quantifiers
    
- requires / ensures / invariant
    

## len-1: Data/domain layer

- lists, sets, maps
    
- indexing
    
- sequences
    
- arithmetic
    
- common predicates like `Sorted`, `Unique`, `Subset`
    

## len-2: Interface layer

- controlled natural sugar
    
- ambiguity markers
    
- notes
    
- provenance
    
- confidence / alternative interpretations
    

This decomposition keeps the formal kernel small.

---

# 10. Examples

## Example 1: minimum element

```text
spec min(xs: List<Int>) -> result: Int
requires NonEmpty(xs)
ensures result in xs
ensures forall y in xs . result <= y
```

Possible Lean-ish compilation:

```text
def MinSpec (xs : List Int) (result : Int) : Prop :=
  NonEmpty xs ∧
  result ∈ xs ∧
  ∀ y, y ∈ xs -> result ≤ y
```

---

## Example 2: sort

```text
spec sort(xs: List<Int>) -> ys: List<Int>
ensures Sorted(ys)
ensures Permutation(ys, xs)
open Stability
```

This is nice because it captures underspecification directly.

---

## Example 3: remove duplicates preserving order

```text
spec dedup(xs: List<Int>) -> ys: List<Int>
ensures Unique(ys)
ensures forall y in ys . y in xs
ensures forall x in xs . x in ys
ensures Subsequence(ys, xs)
```

This now distinguishes it from sorting.

---

## Example 4: stack push/pop law

```text
type Stack<T>

fn push(s: Stack<T>, x: T) -> Stack<T>
fn pop(s: Stack<T>) -> Stack<T>
fn top(s: Stack<T>) -> T
rel NonEmpty(s: Stack<T>)

axiom pop_push:
    forall s: Stack<T>, x: T .
        pop(push(s, x)) = s

axiom top_push:
    forall s: Stack<T>, x: T .
        NonEmpty(push(s, x))

axiom top_of_push:
    forall s: Stack<T>, x: T .
        top(push(s, x)) = x
```

This shows `len` can express knowledge, not only function contracts.

---

## Example 5: natural language-friendly surface

```text
spec min(xs: List<Int>) -> result: Int
requires xs is non-empty
ensures result is in xs
ensures for every y in xs, result <= y
```

Normalized core:

```text
spec min(xs: List<Int>) -> result: Int
requires NonEmpty(xs)
ensures Member(result, xs)
ensures forall y in xs . result <= y
```

---

# 11. Suggested notion of equivalence in `len`

This is probably central for you.

We likely need at least three equivalence relations.

## 11.1 Syntactic equivalence

Exact same text.

[  
e_1 \equiv_{\text{syn}} e_2  
]

Too strict.

## 11.2 Normalized equivalence

After desugaring and normalization, same core form.

[  
N(e_1) = N(e_2)  
]

Useful operationally.

## 11.3 Semantic equivalence

They denote the same meaning in all models.

[  
e_1 \equiv_{\text{sem}} e_2  
]

This is the real notion, but harder to check.

So in practice:

- use normalized equivalence as an approximation
    
- target semantic equivalence as the ideal
    

---

# 12. A compact first version of `len`

Here is a concrete minimal proposal.

```text
type Int
type Bool
type List<T>

rel Member(x: T, xs: List<T>)
rel Sorted(xs: List<Int>)
rel Unique(xs: List<T>)
rel Permutation(xs: List<T>, ys: List<T>)
rel NonEmpty(xs: List<T>)

spec min(xs: List<Int>) -> result: Int
requires NonEmpty(xs)
ensures Member(result, xs)
ensures forall y in xs . result <= y

spec sort(xs: List<Int>) -> ys: List<Int>
ensures Sorted(ys)
ensures Permutation(ys, xs)
open Stability

define SetEqual(xs: List<T>, ys: List<T>) :=
    forall x: T . Member(x, xs) <-> Member(x, ys)
```

This is already a viable seed language.

---

# 13. My recommendation for next refinement

The next useful step is to freeze a **tiny core calculus** for `len`, not too big.

I would choose exactly these primitives first:

- `type`, `fn`, `rel`, `define`, `spec`
    
- `requires`, `ensures`, `invariant`, `open`
    
- equality and order relations
    
- `forall`, `exists`
    
- `in`
    
- conjunction, implication, negation
    

That is enough to start generating examples and proving properties without getting lost.

Then later we can add:

- state / heap / mutation
    
- temporal aspects
    
- probabilistic claims
    
- refinement
    
- effect systems
    

Those should come later.

---

# 14. One clean slogan for `len`

A concise identity for the language would be:

**`len` is a lazy semantic specification language: more explicit than natural language, less rigid than Lean, and structured enough to compile into logic.`**

That is a very good design target.

Next, the most productive move is to do exactly what you proposed:  
define `len` more sharply in three passes:

1. core syntax,
    
2. formal semantics,
    
3. worked examples with normalization and compilation.