# quasi

Functions must come with a `quasi:` block, but the contents of that block are is validated only syntactically at level 1. The parser should recognize the block structure and host-language slots, but it should preserve unrecognized lines as narrative content instead of rejecting them outright.

`quasi` is an embedded block with a DSL language for algorithm sketches, witness constructions, and proof-oriented pseudocode.
Hence, already at this point it becomes clear that quasi can come in many flavours or styles:
- a proof sketch with `assume`, `have`, and `show` steps
- a procedural algorithm outline with `let` and `choose` steps (but see the examples/quasi/sorting/** sketches for a more structured style)
- a functional algorithm outline with `take`, `yield`, and `emit` steps and more declarative control flow
- a unit-test style sketch with `given` and `assert` steps
- a domain-specific sketch with custom step words and idioms

It is deliberately not a fully formal language in the same sense as core `len.l1` formulas. The goal is different: provide a structured, readable proxy for the style found in algorithm texts, proof sketches, and academic pseudocode, while still connecting to declared `type`, `rel`, `const`, and `fn` declarations from the host language.

Once again, `quasi` is nor a rule engine but a small block or guidance to implement an algorithm in an academic pseudo code style. 

## Design Goals

- provide enough grammatical sketches for many styles to deduce the grammar
- design a high level "check" or validation of the quasi block. Namely we check that each style has a fixed list of keywords or syntactical regexps, which must be followed. So the `len validate` command of the cli tool will break only if those checks do not match.

futhermore, the design goals are:

- express algorithmic or proof-oriented steps in a way that reads like textbook pseudocode
- interact directly with host-level names, types, terms, and formulas
- stay open-ended enough for academic prose and domain-specific idioms
- keep the boundary with `len.l1` explicit so the core parser remains stable
- support gradual validation: parse what is structured, preserve what is narrative

## Non-Goals

- replacing core `len.l1` formulas or types
- making `quasi` a fully executable language in the first iteration
- forcing every valid academic proof sketch into a rigid formal grammar
- making the host parser dynamically extensible based on `quasi` contents

## Design Sources

This design follows common patterns found across classic algorithm texts and mathematical writing:

- imperative step words such as `let`, `set`, `choose`, `return`, `repeat`, and `for each`
- proof-step words such as `assume`, `have`, `show`, `suffices`, and `contradiction`
- indentation and block headers instead of dense punctuation
- inline formulas and term expressions where the mathematical content matters more than machine-level syntax

The intent is to approximate the style of textbook pseudocode, not to reproduce any single source verbatim.

## Host Embedding

The host grammar should treat `quasi` as a delimited block region attached only to an `fn` declaration.

Preferred host form:

```text
fn even_plus_even(a: Int, b: Int) -> sum: Int
	implements EvenPlusEven(a, b, sum)
	requires Even(a)
	requires Even(b)
	ensures sum = add(a, b)
	ensures Even(sum)

quasi:
	choose ka such that a = double(ka)
	choose kb such that b = double(kb)
	let k := add(ka, kb)
	show sum = double(k)
```

The important point is that `quasi` appears as an indented `QUASI_BLOCK` owned by a host declaration that already provides parameters, result names, and declared obligations.

The intended owner is `fn`.

That lets one declaration carry both roles cleanly:

- the contract, via clauses such as `requires`, `ensures`, and `implements`
- the implementation sketch, via `quasi:`

This is a better fit than a top-level `quasi Name` declaration, because the pseudocode can refer directly to function parameters, result binders, and declared obligations.

## QUASI_BLOCK Contract With len.l1

The host parser is responsible for:

- recognizing the `quasi:` header inside an `fn`
- capturing the indented block as a `QUASI_BLOCK`
- recording source spans and indentation structure
- providing the enclosing symbol environment to the quasi block parser

The quasi block parser is responsible for:

- parsing recognized step forms inside the block
- identifying embedded host terms and formulas in well-known slots
- building a quasi AST for tooling and later elaboration
- preserving unknown or narrative lines instead of rejecting them too early

This split keeps `len.l1` stable while allowing `quasi` to evolve, but it does not require a separate language level or a separate top-level parser.

## Fn Contract Shape

The cleanest reading is:

- `requires` for preconditions
- `ensures` for postconditions
- `implements` for the abstract relation or semantic intention that the function implements

Example:

```text
fn bubble_sort(input: Seq) -> output: Seq
	implements BubbleSort(input, output)
	ensures Sorted(output)
	ensures Permutation(input, output)

	quasi:
		let n := length(input)
		let output := input
		...
		return output
```

This is preferable to names such as `bubble_sort_fun`, because the `fn` header already makes the declaration function-like.

If the relation link is not needed, an `fn` may use `ensures` clauses alone.

## Quasi Styles

The next syntax improvement should be to make the quasi style explicit in the header.

Preferred short form with default style:

```text
quasi:
	assume P
	show Q
```

Preferred descriptive form with custom style:

```text
quasi using style ProofSketch:
	let i := zero
	...
```

This gives the validator a fixed style contract for the block.

A style profile should specify:

- the allowed leading keywords
- the allowed block forms
- the expected host-language slots after each keyword
- whether narrative lines are accepted

Style profiles should be loadable from:

- built-in presets
- repository-local declarations
- plugins or user-provided configuration

That means quasi can stay open-ended without becoming arbitrary.

## Relationship To spec

`spec` should remain declarative and work with relations and types using logical formulas.

Use `spec` for statements such as:

- definitions
- claims
- properties
- equivalences

That is where `given` and `must` belong.

Use `fn` when the declaration has:

- feel like algorithmic or operational intent
- parameters and a result binder
- always an attached `quasi:` block (aka body of the pseudo algorithm)

## Names In Scope

Inside a `QUASI_BLOCK`, the following names are visible:

- imported host declarations
- top-level declarations in the current module
- parameters introduced by the enclosing `fn`
- result names introduced by the enclosing `fn`
- each block must have a specific style defined with XYZ keyword

Type checking should be host-driven. A quasi statement does not invent a new notion of type or predicate; it only uses host declarations.

## Core Statement Families

`quasi` should be line-oriented and indentation-sensitive. A full formal grammar is intentionally avoided, but the language should recognize a compact set of statement schemas.

### Binding And Assignment

```text
let x := term
let x: T := term
set x := term
choose x such that formula
choose x: T such that formula
take x := term
```

Use these for witness construction, local names, and algorithm state updates.

### Assertions And Proof Steps

```text
assume formula
have formula
show formula
suffices formula
assert formula
note text...
```

These support mathematical argument without pretending every step is executable.

### Control Flow

```text
if formula:
	...
else:
	...

for each x in term:
	...

while formula:
	...

repeat:
	...
until formula
```

These forms are useful for algorithm descriptions without pulling extra words into the core grammar.

### Procedure-Style Steps

```text
apply theorem_name to term_list
invoke rel_or_fn(args)
return term
emit formula
```

These are optional structured conveniences for later tooling. They should elaborate to host-level references when possible.

## Embedded Host Syntax

Every quasi statement contains one or more host-language slots.

- `term` positions are parsed with the host expression parser
- `formula` positions are parsed with the host formula parser
- `T` positions are parsed with the host type-expression parser

Examples:

- in `let k := add(ka, kb)`, `add(ka, kb)` is a host term
- in `choose k such that a = double(k)`, `a = double(k)` is a host formula
- in `if Even(n):`, `Even(n)` is a host formula

This gives `quasi` semantic leverage without requiring a second theorem language.

## Open-Endedness Without Chaos

`quasi` needs to be open-ended, but not shapeless.

The right compromise is a two-tier model:

- hard syntax for block structure, indentation, and a small set of leading verbs
- soft syntax for the rest of the line, with host expressions and formulas parsed where expected

Unknown lines should be preserved as narrative statements. Tooling may warn that a line is unstructured, but it should not fail the whole declaration unless the user requests strict mode.

That matches how academic pseudocode actually behaves: some steps are precise enough to analyze, some are explanatory glue.

## Validation Model

Validation should be incremental.

Level 1 validation:

- an enclosing `fn` may contain a `quasi:` block
- indentation is well-formed
- referenced host names resolve where the statement schema requires them
- embedded terms and formulas are well-formed

Future quasi validation may add:

- local binder tracking across nested blocks
- checking that `return` matches declared result names or result types
- checking that `show` steps line up with enclosing `ensures`, `implements`, or theorem goals
- checking that the declared style profile accepts each step keyword and line schema
- elaboration from `choose` into explicit existential obligations

## Recommended Direction For len.l1

1. Remove the current top-level `quasi Name / case / when` form.
2. Make `quasi:` an embedded `QUASI_BLOCK` owned only by `fn`.
3. Treat `spec` as declarative and keep algorithmic content in `fn`.
4. Extend `fn` with signature and contract clauses such as `requires`, `ensures`, and `implements`.
5. Add an optional style-bearing quasi header such as `quasi using algorithm/basic:`.
6. Parse quasi as a dedicated block routine inside the main parser, not as a separate language level.
7. Keep quasi-local words such as `choose`, `show`, `set`, and `repeat` out of the global reserved-word set.

## Example: Proof-Oriented Contract

```text
type Int

rel Even(Int)
fn add(a: Int, b: Int) -> Int
fn double(k: Int) -> Int

define EvenDef(x: Int) := exists k: Int. x = double(k)

fn even_plus_even(a: Int, b: Int) -> sum: Int
	requires Even(a)
	requires Even(b)
	implements EvenPlusEven(a, b, sum)
	ensures sum = add(a, b)
	ensures Even(sum)

	quasi:
		choose ka such that a = double(ka)
		choose kb such that b = double(kb)
		let sum := add(a, b)
		let k := add(ka, kb)
		have sum = add(double(ka), double(kb))
		show sum = double(k)
```

## Example: Algorithm Sketch

```text
fn count_steps(start: State) -> count: Nat
	quasi using ProceduralStyle:
		let current := start
		let count := zero
		while not Done(current):
			set current := step(current)
			set count := succ(count)
		return count
```

This is the intended role of quasi: a disciplined pseudocode layer, not a malformed mini-Prolog.
