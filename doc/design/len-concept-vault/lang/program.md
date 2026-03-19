# Program In `len`

## The big picture

`len` is a semantic intermediate language for specifications, semantics, and domain knowledge.
A specification language like `len` sits in an unusual position between:

- natural language
- pseudocode
- formal logic
- verification languages

Where by verification languages we mean not just provers or theorem provers like SMT encodings, Coq, Agda, Lean, F*, and so on. And by formal logic we mean not just FOL, but also domain-specific logics, algebraic theories, and so on.

One and first, big idea in the design of `len` is that it should not be a single flat language, but rather a family of related languages that are layered on top of each other with different responsibilities and different intended users.

0. **natural language level or `l0`** is the authoring layer with all the ambiguity and flexibility of natural language, but it sounds like a programming problem statement, math/physics/engineering problem or anything that asks for a rigorous solution rather than arbitrary prose. So, yes `l0` is natural language, but it is a controlled and expected to be LLM-translatable into next layer.

1. **relation level one or `l1`** is the first real formal layer, but it is still designed for humans. It is a relational and contract-first language with a syntax that supports readability and local reasoning. It is not pseudocode, but it is also not a low-level logic. It is influenced by sorted logics, Haskell types, Prolog predicates and relations, and TLA+ specs. It is the main authoring layer for domain declarations, relations, and contracts.

2. **canonical semantic level or `l2`** is the canonical semantic layer. It is more precise than `l1` because it does not merely restate contracts in a lower-level syntax; it identifies what kind of semantic artifact each declaration, relation, law, or implementation attachment actually is. In `l2`, domains are not semantic primitives. They become namespaces and provenance boundaries that contain canonical artifacts such as categories, algebraic structures, theories, effects, morphisms, theorems, and implementation objects. In that sense, `l2` is the exact semantic IR used for normalization, comparison, composition, mapping, guided code generation, backend compilation, proof, and validation. Types, relations, and contracts from `l1`, once normalized into `l2`, become explicit semantic building blocks whose final shape depends on the domain. In programming or infrastructure settings, implementation objects are more likely to be code artifacts such as services, modules, or classes. In mathematics, they are more likely to be proof artifacts.

Another and second big idea is that the layering is not just cosmetic, but its in the hard of development process. We call this a semantic loop (somewhat similar to the print-eval-read loop in programming languages, but with more layers and more AI involvement):

- one can move in both directions in the loop - up and down, e.g.

```
Human Language
      │
      │  S  (semantic interpretation)
      ↓
Knowledge Representation
      │
      │  reasoning / inference
      ↓
Knowledge Representation
      │
      │  I  (semantic interpolation)
      ↓
Human Language
```

- one can start with `l0`, then arrive at `l1`, and if all looks good, immidiately jumop to code generation. Nowerver days this is done with LLMs and prompt engineering, but the idea is that with a more structured language like `len`, the jump from `l0` to code generation can be more direct and more reliable. Furthermore, `l1` has specs (contracts) as the main semantic unit, so this helps with testing and validation of the generated code. 

- if one goes deeper, then the jump from `l1` to  `l2` before wor with the code generation in between can help refine and epxlore deepr deisng properties of the domain.

Essentially, the pragmatics part of len includes:

- code generation and testing
- proof and validation
- knowledge graph construction and maintenance
- semantic search and retrieval
- semantic interoperability and mapping
- semantic change management and evolution
- interpolation and extrapolation of knowledge

## Layering Rationale

The language should not be defined as a single flat syntax.
It should be defined in layers with different responsibilities and different syntax:

- `l0`: authoring and clarification in *problem-stating* natural language
- `l1`: relational contracts and domain declarations
- `l2`: canonical semantic artifacts for classification, normalization, mapping, and backend realization

The point of the layering is not cosmetic.
Each layer has a different job in the semantic loop:

Human input -> `l0` -> `l1` -> `l2` -> code / backend logic / proof / validation

The main principle is:

`l0` is for saying what is meant.
`l1` is for writing the intended relation in a human-oriented formal way.
`l2` is for classifying the semantic object precisely enough to normalize, compare, compose, map, and compile.

## Design Direction

The current direction is:

- `l0` stays close to human language
- `l1` is declaration-first and contract-first
- `l2` is a canonical semantic artifact layer, not just a lower-level restatement of `l1`

This means `l2` should not merely restate `l1` with more quantifiers.
It should expose what kind of semantic artifact something actually is:

- a category-like object
- an algebraic fragment
- a theory bundle
- an effectful context
- a morphism between semantic objects
- a theorem
- an implementation artifact

## Layer Summary

| Layer | Purpose | Main style | Main output |
| --- | --- | --- | --- |
| `l0` | authoring, ambiguity surfacing, refinement dialogue | controlled natural language | candidate meaning |
| `l1` | domain declarations and contracts | relations, types, specs, quasi steps | human-oriented formal spec |
| `l2` | canonical semantic classification and normalization | namespaces, artifact kinds, semantic links | semantic IR for mapping, proof, validation, and backend generation |

## Level 0

`l0` is natural language, but not arbitrary prose.
It is the authoring layer used to state intent, examples, missing assumptions, and clarifications before committing to a stronger formal form.

`l0` should be permissive, but a small set of authoring markers is useful.

### Recommended `l0` Keywords

These are not logical primitives.
They are discourse markers that help an agent produce `l1`.

- `goal`
- `context`
- `given`
- `when`
- `then`
- `because`
- `example`
- `counterexample`
- `ambiguity`
- `question`
- `note`

### Role Of `l0`

`l0` should support:

- stating the problem in domain language
- recording ambiguities without forcing premature precision
- giving examples and counterexamples
- exposing missing assumptions
- preparing an agent-to-`l1` normalization step

### Example `l0`

```text
goal: evaluate calculator expressions safely
context: arithmetic expressions may contain division
given: the application should reuse arithmetic knowledge
then: division by zero must be rejected
example: 2 + 3 should evaluate to 5
ambiguity: should division return an error or an optional value?
question: do we want integer division or rational division?
```

## Level 1

`l1` is the first real formal layer.
It is still designed for humans, so it keeps the syntax local, contract-oriented, and readable.

The stable core of `l1` is:

- declarations first
- specs as the main semantic unit
- formulas in lightweight FOL-style syntax
- optional quasi steps for proof hints or implementation steering

### Recommended Core `l1` Keywords

These should be treated as the main `l1` keyword set.

- `define`
- `decl`
- `symbol`
- `domain`

- `type`
- `rel`
- `fn`
- `spec`
- `requires`
- `ensures`
- `open`
- `quasi`
- `choose`
- `let`
- `show`
- `forall`
- `exists`
- `and`
- `or`
- `not`

### `l1` Operators And Reserved Forms

These are part of the concrete syntax even if they are not word-keywords.

- `->`
- `<->`
- `=`
- `!=`
- `<`
- `<=`
- `>`
- `>=`
- `:`
- `:=`

declared as symbols in the language, but they are reserved for their logical meaning.

### `l1` Scope

`l1` should cover:

- types and signatures
- relations and helper functions
- contracts, policies, queries, and validation rules through `spec`
- local quantification
- witness-oriented statements
- partially constructive quasi blocks

### Example `l1`

```len
type Expr
type Op
type Int

rel Const(Expr, Int)
rel Binary(Expr, Op, Expr, Expr)
rel AddOp(Op)
rel DivOp(Op)
rel Eval(Expr, Int)

fn add(a: Int, b: Int) -> Int

spec eval_const(e: Expr) -> v: Int
	requires Const(e, v)
	ensures Eval(e, v)

spec eval_add(e: Expr, op: Op, l: Expr, r: Expr) -> v: Int
	requires Binary(e, op, l, r)
	requires AddOp(op)
	requires exists a: Int. Eval(l, a)
	requires exists b: Int. Eval(r, b)
	ensures exists a: Int. exists b: Int.
		Eval(l, a) and Eval(r, b) and v = add(a, b)
	ensures Eval(e, v)

spec no_div_by_zero(e: Expr, op: Op, l: Expr, r: Expr) -> ok: Int
	requires Binary(e, op, l, r)
	requires DivOp(op)
	ensures exists b: Int. Eval(r, b) and b != 0
	open DivisionResultPolicy
```

### `l1` Design Constraint

`l1` should remain relational and contract-first.
It should not collapse into ordinary pseudocode.

This means:

- `spec` remains primary
- `fn` is secondary and usually interpreted
- `quasi` is allowed for steering, not for replacing the contract language

## Level 2

`l2` is the canonical semantic layer.
It is more precise than `l1` because it identifies what kind of object each artifact is.

In `l2`, domains are not semantic primitives.
They become namespaces that contain canonical artifacts such as categories, algebraic structures, theories, effects, morphisms, theorems, and implementation objects.

### Canonical `l2` Artifact Keywords

This is the recommended small family for `l2`.

- `namespace`
- `cat`
- `alg`
- `theory`
- `effect`
- `morph`
- `theorem`
- `impl`

### Intended Meaning Of Each `l2` Keyword

- `namespace`: packaging and provenance boundary for a domain
- `cat`: structured semantic object or semantic domain object
- `alg`: algebraic signature together with its main laws
- `theory`: bundle of logical laws, constraints, or derived rule families
- `effect`: contextual computation structure such as error, state, nondeterminism, or time
- `morph`: mapping or interpretation between canonical objects
- `theorem`: checked or accepted truth within the normalized semantic graph
- `impl`: implementation-level artifact derived from or attached to the semantic layer

### `l2` Relation Keywords

These are the main relation clauses that make `l2` useful as a graph of semantic artifacts.

- `uses`
- `extends`
- `implements`
- `refines`
- `maps`
- `preserves`
- `via`
- `laws`

Not every project needs all of them immediately, but this is the right direction for the canonical layer.

### Why `l2` Needs More Than `cat`

Not everything in canonical form is best modeled as a category.

Examples:

- parity laws are better represented as a `theory`
- additive and multiplicative signatures are better represented as `alg`
- maybe/error behavior is better represented as an `effect`
- evaluation from syntax to arithmetic is better represented as a `morph`
- an extracted evaluator is better represented as an `impl`

If everything became a `cat`, the representation would be too coarse.

### Example `l2`: Arithmetics

```len
namespace math.arithmetics

alg additive
	maps:
		carrier -> Int
		zero -> 0
		combine -> add
	laws:
		forall a: Int. add(a, 0) = a
		forall a: Int. add(0, a) = a
		forall a: Int. forall b: Int. forall c: Int.
			add(add(a, b), c) = add(a, add(b, c))

alg multiplicative
	maps:
		carrier -> Int
		one -> 1
		combine -> mul

theory parity
	laws:
		forall x: Int. Even(x) <-> exists k: Int. x = mul(2, k)
		forall x: Int. Odd(x) <-> exists k: Int. x = add(mul(2, k), 1)

cat integer_semantics
	uses: additive
	uses: multiplicative
	uses: parity
```

### Example `l2`: Calculator App

```len
namespace app.calculator

theory syntax
	laws:
		forall e: Expr. ConstExpr(e) -> WellFormed(e)
		forall e: Expr. forall l: Expr. forall r: Expr.
			AddExpr(e, l, r) -> WellFormed(e)

effect maybe_error
	laws:
		ErrorPropagates
		SuccessCarriesValue

cat eval_semantics
	uses: syntax
	uses: math.arithmetics.integer_semantics
	uses: maybe_error

morph eval : eval_semantics -> math.arithmetics.integer_semantics
	via: interpret
	preserves:
		additive

impl recursive_eval
	refines: eval
```

## Normalization Direction

The main normalization path should be:

`l0` -> `l1` -> `l2`

with the following intent:

- `l0` resolves ambiguity enough to produce candidate declarations and contracts
- `l1` expresses the domain with readable relations and specs
- `l2` classifies the result into canonical artifact kinds

### Typical `l1` To `l2` Rewrites

Common cases:

- a reusable group of operations and laws in `l1` normalizes into `alg`
- a logical bundle of invariants or closure statements normalizes into `theory`
- a domain object with structured semantic meaning normalizes into `cat`
- an evaluation or interpretation relation normalizes into `morph`
- an operational context normalizes into `effect`
- a proved property normalizes into `theorem`
- a chosen algorithm or extracted executable artifact normalizes into `impl`

### Example Rewrite Shape

`l1` domain fragment:

```len
type Int
rel Even(Int)
fn add(a: Int, b: Int) -> Int

spec even_sum_closed(a: Int, b: Int) -> sum: Int
	requires Even(a)
	requires Even(b)
	ensures sum = add(a, b)
	ensures Even(sum)
```

Normalizes into mixed `l2` artifacts:

```len
namespace math.arithmetics

alg additive
theory parity

theorem even_sum_closed
	uses: additive
	uses: parity
```

This is the core reason to keep `l2` richer than only `cat`.

## Practical Layer Constraints

To keep the language coherent:

1. `l0` should not pretend to be logically complete.
2. `l1` should be the main authoring format for contracts.
3. `l2` should be small, typed, and canonical.
4. `l2` should classify artifacts rather than flatten everything into one container type.
5. `spec` is central in `l1`, but in `l2` a `spec` may normalize into `theory`, `theorem`, `morph`, `effect`, or `impl` depending on meaning.

## Current Recommendation

For now, the exact keyword inventory should be treated as:

### `l0`

- `goal`
- `context`
- `given`
- `when`
- `then`
- `because`
- `example`
- `counterexample`
- `ambiguity`
- `question`
- `note`

### `l1`

- `type`
- `rel`
- `fn`
- `spec`
- `requires`
- `ensures`
- `define`
- `open`
- `quasi`
- `choose`
- `let`
- `show`
- `forall`
- `exists`
- `and`
- `or`
- `not`

### `l2`

- `namespace`
- `cat`
- `alg`
- `theory`
- `effect`
- `morph`
- `theorem`
- `impl`

This is a compact and defensible split.
It keeps `l1` centered on relations, types, and contracts, while making `l2` a true canonical semantic layer rather than a more decorated version of `l1`.
