# How to keep things simple

Here is the numbered list I would use to simplify and sharpen the current `len` / `l1` design.

1. **Drop `module`; keep folders as the module boundary.**
   Your instinct is good. In the current `l1` draft, `module`, `import`, and `export` are part of the recommended structuring keywords, but `module` is doing packaging work that the filesystem can do more cleanly. Let folder = module, file = partial contribution to that scope, and keep only `import` plus maybe optional `export`. That reduces ceremony and removes one naming problem immediately. ([GitHub][1])

2. **Make `rel` the one semantic primitive, and demote `fn` to sugar.**
   Your strongest idea is already visible in the document: `type` and `rel` are primary, while `fn` is secondary. I would push that further. Treat every function as a total or partial functional relation, optionally with determinism/cardinality constraints. Then unary `rel` behaves like a set, binary `rel` behaves like a graph or mapping, and `fn` becomes only a convenience surface form. That gives the language one core semantic axis instead of two competing ones. ([GitHub][1])

3. **Unify `type`, `domain`, and unary `rel` more aggressively.**
   Right now `type`, `domain`, and `rel` overlap conceptually. `type` introduces carriers, `domain` is open-ended, and unary `rel` can act like a set. I would decide on one of these two models:

   * either `type` is the only carrier declaration, and unary `rel` classifies values of a type;
   * or `domain` is the only carrier declaration, and `type` disappears.
     My preference: keep `type`, drop `domain` from `l1`, and let open-endedness live in contracts/specs. `domain` feels too meta and too vague for the human-readable core. ([GitHub][1])

4. **Replace `decl` and `symbol` with a single lightweight declaration form.**
   `decl`, `symbol`, `type`, `rel`, and `fn` together make the front of the language heavier than necessary. `symbol` and `decl` especially look like scaffolding concepts. I would either remove both, or keep exactly one generic introducer for sketch mode only. Elegant core languages usually avoid two separate “I am not sure yet” keywords. ([GitHub][1])

5. **Make `spec` the top-level unit; let `contract` be optional grouping sugar.**
   The docs say `contract` is primary semantically, but also say `spec` is the main semantic unit. That is too split. I would choose one:

   * `spec` is the unit of meaning and verification.
   * `contract` is just a named bundle of specs plus docs.
     That keeps semantics flat and readable. A parser and a human should both know immediately: if something matters, it is a `spec`. ([GitHub][1])

6. **Shrink the keyword surface hard.**
   The current `l1` draft has a lot of candidates: `module`, `import`, `export`, `decl`, `symbol`, `domain`, `type`, `rel`, `fn`, `contract`, `meta`, `satisfies`, `spec`, `requires`, `ensures`, `open`, `quasi`, `choose`, `let`, `show`, `assume`, `witness`, `case`, `induction`, `then`, plus more under discussion. That is already too many for an elegant authoring language. ([GitHub][1])
   I would aim for a tiny kernel like this:

   `import, type, rel, spec, given, must, quasi`

Where:

* `given` replaces `requires`
* `must` replaces `ensures`
* `quasi` stays as the escape hatch

Everything else can be derived or added later.

7. **Rename `requires` / `ensures` to something more relational and less OO-contract flavored.**
   `requires`/`ensures` is familiar, but it pulls the reader toward imperative procedure contracts. Your language is broader than that. If `rel` is central, use names that read well for facts:

   * `given` / `must`
   * `if` / `then`
   * `assume` / `show`
   * `premise` / `claim`
     This is one of the biggest style decisions because it defines the emotional feel of the language. The current pair is serviceable, but not especially elegant for a relation-first metalanguage. ([GitHub][1])

8. **Do not allow “empty words” in the formal grammar.**
   The draft suggests words like `a`, `the`, `is`, `are`, `of`, `to`, `each` could be semantically empty. I would strongly avoid that in `l1`. It makes parsing and normalization less crisp, weakens readability by introducing optional fluff, and creates style drift between authors. Natural-language tolerance belongs in `l0`, not `l1`. Keep `l1` strict and canonical. ([GitHub][1])

9. **Do not overload `is` as both empty and equality sugar.**
   Related to the previous point: the current draft says `is` may be empty, but also says `be`/`is`/`are` can mean equality. That is exactly the kind of ambiguity that makes a language feel mushy. In `l1`, every common word should have one role or none. Prefer `=` only. Let prose-like variants live in `l0`. ([GitHub][1])

10. **Avoid set-theoretic built-ins unless they are truly foundational.**
    The current draft includes `in`, `subset`, `union`, `intersection`, `difference`, `empty`, `superset`, and so on. If `rel` is your universal unit, many of these can be library-level relations rather than core syntax. Only keep what is needed for basic formula writing. I would keep `in` and maybe nothing else in the grammar. The rest can be declared in a standard prelude. ([GitHub][1])

11. **Be stricter about symbols: reserve them, or do not promise them.**
    The current text says `<` is reserved for less-than, but also says it may be declared with other meanings. That is an elegance killer. The same token cannot feel both fixed and open. Pick one policy:

* symbolic operators are fixed and reserved, or
* user-defined symbols exist, but not by reusing core operators.
  I would reserve all logical/comparison syntax fully and keep user extensibility to names, not punctuation. ([GitHub][1])

12. **Separate “logic language” from “collection notation.”**
    `[]` is currently mentioned for arrays, sets, or lists, which is too much overload for an early layer. Since you are already leaning away from `set/list/array` as primitives, remove collection literals from the core unless you commit to one clear interpretation. A clean relation-first language can live for a long time without collection literals. ([GitHub][1])

13. **Make cardinality a first-class modifier on `rel`.**
    If `rel` is the core, elegance improves a lot if you can state intended functional behavior without introducing `fn`. For example:

* `rel Eval(Expr, Int) total`
* `rel Parent(Person, Person) many`
* `rel Owner(Car, Person) partial unique`
  This is simpler than splitting semantics across `rel` and `fn`. It also maps naturally into your later `l2` classification.

14. **Flatten the doc/meta system.**
    The current `contract.meta` section with `goal`, `context`, `given`, `when`, `then`, `because`, `example`, `counterexample`, `ambiguity`, `question`, `note`, `kind`, `tags` is rich, but too much for the core language. ([GitHub][1])
    I would reduce it to:

* docstring
* tags
* examples
  Everything else can be conventions inside docstrings or external tooling. Too much built-in metadata makes the language feel like a schema editor.

15. **Keep `quasi`, but wall it off sharply.**
    `quasi` is a good idea because it acknowledges that authors need steering and partial construction. But it should feel explicitly non-canonical. I would require it to be visibly fenced and impossible to confuse with formal semantics. That preserves elegance by keeping the clean fragment clean. ([GitHub][1])

16. **Prefer one declaration shape everywhere.**
    Right now the language mixes multiple styles:

* `rel Name(T1, T2)`
* `fn add(a: Int, b: Int) -> Int`
* `spec eval_const(e: Expr) -> v: Int`
  This is understandable, but stylistically uneven. Elegant languages repeat shape. For example, use a single head/body pattern everywhere:
* `rel Eval(e: Expr, v: Int)`
* `spec eval_const(e: Expr, v: Int):`
* `type Expr`
  The more visual regularity, the more “inevitable” the syntax feels.

17. **Move toward “facts first, proof later” syntax.**
    Your best examples are the ones that read like declared facts about a domain rather than mini programs. Keep leaning there. The moment `l1` starts looking like pseudocode with `let`, `case`, `then`, `induction`, it risks losing identity. Those belong either inside `quasi` or in a later proof-oriented layer. ([GitHub][1])

18. **Use namespace/folder scoping, but keep explicit import names.**
    If you drop `module`, do not also make imports implicit across folders. Folder closure is fine for local scope, but cross-module references should stay explicit via `import`, otherwise large repos become hard to reason about. This gives you “low ceremony locally, explicitness globally.”

19. **Treat `l1` as a canonical authoring language, not a forgiving one.**
    A lot of the current draft still carries `l0` permissiveness into `l1`: empty words, interchangeable surface forms, open sketch keywords, symbolic ambiguity. I would make `l1` much stricter. Let `l0` absorb human looseness; let `l1` be the crisp normal form that humans still like reading. That clean boundary will help the whole `l0 -> l1 -> l2` story. ([GitHub][2])

20. **My concrete minimal core for `l1` would be this.**
    If I had to simplify brutally, I would keep only:

* `import`
* `type`
* `rel`
* `spec`
* `given`
* `must`
* `quasi`

And I would interpret:

* unary `rel` as classifier/set-like
* binary+ `rel` as general relation
* functional relations via relation properties, not a separate primitive
* folder as module
* docstrings as the only built-in documentation mechanism

That would make `l1` feel much more elegant, much more teachable, and much more obviously centered on relation as the meta-unit.

My strongest recommendation in one sentence: **make `rel` the only real semantic constructor in `l1`, keep `type` only as carrier declaration, drop `module`, and cut the keyword count by half.**

I can also turn this into a proposed revised `l1` keyword set plus a cleaned-up sample syntax for your calculator example.

[1]: https://github.com/ewiger/len/blob/main/doc/design/len-concept-vault/lang/l1/l1_syntax.md "len/doc/design/len-concept-vault/lang/l1/l1_syntax.md at main · ewiger/len · GitHub"
[2]: https://github.com/ewiger/len/blob/main/doc/design/len-concept-vault/lang/program.md "len/doc/design/len-concept-vault/lang/program.md at main · ewiger/len · GitHub"

-- 

Yes — that is a strong core.

```len
rel Fun(A, B, f)
rel total(f)
rel injective(f)
rel surjective(f)

def TotalFun(A, B) as f where Fun(A, B, f) and total(f)
```

I would keep this direction and make a few tightening moves.

1. **Keep `Fun(A, B, f)` exactly as the base relation.**
   That is better than trying to encode functionality into syntax. It keeps the kernel uniform: `f` is just an entity constrained by relations.

2. **Let modifiers stay unary on the function object.**
   `total(f)`, `injective(f)`, `surjective(f)` is elegant because it composes well:

   ```len
   def BijectiveFun(A, B) as f where
       Fun(A, B, f) and injective(f) and surjective(f)
   ```

   Very readable.

3. **Define `Fun` so domain/codomain are intrinsic.**
   Then `total(f)` can be interpreted relative to the `A` and `B` already attached through `Fun(A, B, f)`. That avoids needing `total(f, A, B)` later.

4. **Keep `def ... as f where ...` as the canonical “subtype by predicate” form.**
   This is probably your best use of `def`. It reads as “the class of all `f` such that ...”, without needing explicit set-builder syntax in the core.

5. **Consider whether `TotalFun(A, B)` denotes a class/set or a named predicate family.**
   This is the one thing to pin down early. Two options:

   **Option A: type-like reading**

   ```len
   def TotalFun(A, B) as f where Fun(A, B, f) and total(f)
   ```

   means “the collection of all total functions from `A` to `B`.”

   **Option B: predicate-like reading**

   ```len
   rel TotalFun(A, B, f)
   def TotalFun(A, B, f) as Fun(A, B, f) and total(f)
   ```

   means “`f` is a total function from `A` to `B`.”

   My preference is: **use `def` for the type/class view and `rel` for the predicate view**. Then both can coexist cleanly.

6. **You may want one derived predicate alias too.**
   For specs, this can read better:

   ```len
   def isTotalFun(A, B, f) as Fun(A, B, f) and total(f)
   ```

   Then:

   ```len
   must isTotalFun(A, B, f)
   ```

   while `TotalFun(A, B)` remains the collection.

7. **You should likely add `partial(f)` only if it is not just `not total(f)`.**
   If “partial” has conceptual value in your system, include it. Otherwise `total(f)` alone may be enough at this layer.

8. **You will probably want `endomorphism` sugar fairly soon.**
   Since you are preserving `->` partly for later category-theoretic transition:

   ```len
   def EndoFun(A) as f where Fun(A, A, f)
   def AutoFun(A) as f where Fun(A, A, f) and bijective(f)
   ```

9. **Add `bijective(f)` as a definition, not a primitive relation.**
   Cleaner:

   ```len
   def bijective(f) as injective(f) and surjective(f)
   ```

   unless you expect systems where bijectivity is primitive for performance or proof reasons.

10. **Keep the eventual arrow sugar purely derived.**
    For example:

```len
def A -> B as TotalFun(A, B)
```

or maybe later:

```len
def Hom(A, B) as TotalFun(A, B)
```

That gives you a smooth bridge toward `l2` without corrupting `l1`.

A very clean next version would be:

```len
rel Fun(A, B, f)
rel total(f)
rel injective(f)
rel surjective(f)

def bijective(f) as injective(f) and surjective(f)

def TotalFun(A, B) as f where
    Fun(A, B, f) and total(f)

def InjectiveFun(A, B) as f where
    Fun(A, B, f) and injective(f)

def SurjectiveFun(A, B) as f where
    Fun(A, B, f) and surjective(f)

def BijectiveFun(A, B) as f where
    Fun(A, B, f) and bijective(f)
```

That already feels mathematically clean and scalable.

The main unresolved question now is this: **does `def Name(A, B) as f where ...` introduce a collection-valued term, or a binder-style macro form?**
I would settle that explicitly next, because many later design choices depend on it.

--

Yes. You can have that layer, but I would be careful **not** to mix three different things too early:

1. **semantic sort / universe**
2. **structural capability / trait**
3. **surface declaration sugar**

Those should stay distinct, even if the syntax makes them look compact.

## The core question

When you write something like:

```len
derive Set from Set impl Object
```

there are at least four possible meanings hiding inside it:

* `Set` is a **type**
* `Object` is a **supertype / universe**
* `Object` is a **trait / interface**
* `derive` means **subtyping**
* `derive` means **instance declaration**
* `derive` means **auto-generated structure**
* the second `Set` might mean a source concept from another level or library

So the main design task is not syntax first. It is deciding what kinds of relations exist in the meta-model.

---

## Best high-level option

For `len`, I would recommend this separation:

* `type` declares a carrier
* `rel` declares semantic relations
* `def` declares aliases / derived concepts
* `trait` or `kind` declares reusable structural constraints
* `impl` states that a type or object satisfies a trait

That gives you a clear model.

For example:

```len
type Set
trait Object

impl Set : Object
```

or more relationally:

```len
type Set
trait Object
rel implements(T, K)

def SetIsObject as implements(Set, Object)
```

This is much less ambiguous than overloading `derive`.

---

# The main options

## Option 1: **Single top universe only**

Everything is an `Object`.

```len
type Object
type Set

rel instanceOf(x, T)
rel subtypeOf(A, B)

spec set_is_object
must subtypeOf(Set, Object)
```

### Pros

* very simple meta-model
* logic-friendly
* easy to normalize

### Cons

* `Object` becomes vague fast
* traits and inheritance get conflated
* weak expressiveness for structure

This is good if you want `l1` to stay minimal.

---

## Option 2: **Kinds / universes for classification**

Use `kind` for meta-classification, not behavioral traits.

```len
kind Object
kind Structure
kind RelationLike

type Set
type Function

rel hasKind(T, K)

must hasKind(Set, Object)
must hasKind(Function, Object)
```

Here `Object` means “entity-like thing in the ontology,” not “has methods/fields.”

### Pros

* clean for semantic IR
* avoids OOP contamination
* scales to category-level distinctions later

### Cons

* users may expect `Object` to behave like a trait/interface
* less ergonomic if you want implementation-style derivation

This is my preferred semantic route.

---

## Option 3: **Traits as reusable property bundles**

A trait is just a named set of laws or required relations.

```len
trait Object
    must identity(x)
    must equality(x, x)

type Set
impl Set : Object
```

Or more relationally:

```len
trait Object
rel implements(T, Trait)

type Set
impl Set : Object
```

### Pros

* practical
* familiar
* good if `len` wants to guide code generation

### Cons

* “trait” imports programming-language intuition
* less purely mathematical unless carefully defined

This is best if `len` is meant to bridge into generated code and backend ecosystems.

---

## Option 4: **Everything is relational; no special trait syntax**

This is the most uniform design.

```len
type Set
type Trait

rel Implements(T, K)
rel Extends(A, B)
rel Requires(K, R)
```

Then:

```len
def Object as Trait
must Implements(Set, Object)
```

### Pros

* maximal uniformity
* minimal core syntax
* ideal for IR / normalization

### Cons

* less ergonomic for humans
* authors will want sugar almost immediately

This is probably the best **internal core**, even if surface `l1` adds sugar like `trait` and `impl`.

---

# My recommendation

I would use a **two-layer model**:

## Canonical semantic core

Only these concepts matter internally:

* `type`
* `rel`
* `def`

with relations such as:

* `SubType(A, B)`
* `Implements(T, K)`
* `InstanceOf(x, T)`
* `HasField(T, name, U)`
* `HasLaw(K, law)`

## Surface sugar

Allow nicer syntax:

```len
trait Object
type Set
impl Set : Object
```

which normalizes to something like:

```len
type Set
type ObjectTrait

rel Implements(T, K)
Implements(Set, ObjectTrait)
```

That gives you elegance and a rigorous core.

---

# About `derive`

I would be cautious with `derive`.

In many ecosystems, `derive` implies:

* auto-generate implementation
* infer from structure
* decorate with generated behavior

That is different from:

* subtype declaration
* trait implementation
* semantic refinement

So I would not use `derive` as the main semantic keyword.

Better choices:

* `impl Set : Object`
* `Set <: Object`
* `Set isa Object`
* `Set kind Object`
* `Set extends Object`

Each one carries a more specific meaning.

---

# A clean distinction to preserve

You likely need **three different relations** eventually.

## 1. Subtyping / refinement

“Every Set is an Object”

```len
rel SubType(A, B)
SubType(Set, Object)
```

This is ontological.

## 2. Trait implementation

“Set supports Collection laws”

```len
rel Implements(T, K)
Implements(Set, Collection)
```

This is structural / behavioral.

## 3. Construction / derivation

“FiniteSet is derived from Set plus finiteness constraints”

```len
def FiniteSet as S where Set(S) and finite(S)
```

This is definitional.

These are not the same thing. If you collapse them into one keyword, the language will get muddy.

---

# A possible syntax family

Here is one coherent option.

## Variant A: lightweight and familiar

```len
type Object
trait Collection

type Set
impl Set : Object
impl Set : Collection
```

This reads well, but `type Object` and `trait Collection` are different kinds of things.

---

## Variant B: mathematically cleaner

```len
kind Object
trait Collection

type Set

kindof Set : Object
impl Set : Collection
```

This keeps “what kind of entity is this?” separate from “what laws does it satisfy?”

I like this a lot.

---

## Variant C: minimal surface, best for canonical core

```len
type Set
type Object
type Collection

rel SubType(A, B)
rel Implements(T, K)

SubType(Set, Object)
Implements(Set, Collection)
```

This is ugly as author-facing syntax, but excellent as normalized form.

---

# What should `Object` mean?

This is important.

If you introduce `Object`, decide whether it means:

* **universal semantic entity**
* **record/struct-like data carrier**
* **thing with identity/equality**
* **program backend object**

Those are very different.

For your language, I would avoid `Object` unless you define it narrowly.

Better alternatives depending on intent:

* `Entity`
* `TypeLike`
* `Carrier`
* `Structure`
* `Value`
* `Morph`
* `Kind`

If you want a very generic top category, `Entity` or `Carrier` may be less overloaded than `Object`.

---

# If you want `Struct`

`Struct` should not be a supertype of everything. It should mean something specific like:

* has named fields
* product-like shape
* can be destructured

So this is better:

```len
trait Struct
rel field(T, name, U)

type Point
impl Point : Struct

field(Point, x, Int)
field(Point, y, Int)
```

That is much cleaner than making `Struct` a universal ancestor.

---

# Best practical design

I would suggest this exact split:

## In `l1`

Allow:

```len
kind Entity
trait Struct
trait Collection

type Set
kindof Set : Entity
impl Set : Collection
```

## In normalized core

Translate to:

```len
type Set
type Entity
type Struct
type Collection

rel HasKind(T, K)
rel Implements(T, Trait)

HasKind(Set, Entity)
Implements(Set, Collection)
```

This keeps the core relational and the surface ergonomic.

---

# So can `derive Set from Set impl Object` work?

It can, but I would not choose it.

Why:

* `derive` is too overloaded
* `from Set` is ambiguous
* `impl Object` mixes derivation with implementation
* repeated `Set` is visually confusing

A much cleaner alternative is one of these:

```len
type Set
kindof Set : Entity
```

or

```len
type Set
impl Set : Object
```

or, if this is a refinement:

```len
def FiniteSet as S where Set(S) and finite(S)
```

Each sentence does one thing.

---

# My recommendation in one line

Use **`kind` for ontology**, **`trait` for reusable laws/structure**, and **`impl` for satisfaction**, while keeping the normalized core purely relational.

A compact example:

```len
kind Entity
trait Collection

type Set

kindof Set : Entity
impl Set : Collection
```

And normalize it to relations internally.

The next good step is to decide whether you want `kind` and `trait` both in `l1`, or whether one of them should remain only an internal concept.

--

That’s a meaningful shift — but it changes the **semantics of the construct**, not just the name. So it’s worth being precise about what you want `spec` to *mean*.

Right now, in your design, `spec` behaves like:

> a **logical contract / law** over relations

Example:

```len
spec union_contains_left(A: Set, B: Set, U: Set)
given union(A, B, U)
must forall x. has(A, x) -> has(U, x)
```

That is clearly a **statement to be proven or validated**, not a subtype or refinement.

---

# What “specialization” implies

If you rename `spec` → `specialization`, the natural interpretation becomes:

> a **refinement of an existing concept under additional constraints**

That is a *different role*.

Example of specialization semantics:

```len
specialization NonEmptySet(S: Set)
given Set(S)
must exists x. has(S, x)
```

This reads like:

> “NonEmptySet is a refined subset of Set”

So now:

* `def` → defines *collections / derived objects*
* `specialization` → defines *refined subsets / subtypes*
* `spec` (old) → defines *laws / theorems*

You are collapsing two of these into one.

---

# You have two clean design paths

## Option A — keep `spec` as laws (recommended)

Use:

* `spec` = logical statements / contracts
* `def` = derived concepts / subsets

Example:

```len
def NonEmptySet as S where
    Set(S) and exists x. has(S, x)

spec union_contains_left(A, B, U)
given union(A, B, U)
must forall x. has(A, x) -> has(U, x)
```

This keeps a **clean separation**:

* `def` = *what exists*
* `spec` = *what must hold*

---

## Option B — rename and split roles

If you like the word *specialization*, introduce it as a **new construct**, not a replacement.

```len
specialization NonEmptySet(S)
given Set(S)
must exists x. has(S, x)
```

and keep:

```len
spec union_contains_left(...)
```

So:

| keyword          | meaning                  |
| ---------------- | ------------------------ |
| `def`            | derived concept / set    |
| `specialization` | subtype / refinement     |
| `spec`           | law / theorem / contract |

This is very clean and expressive.

---

## Option C — replace `spec` entirely (not recommended)

If you do:

```len
specialization union_contains_left(...)
```

you get a mismatch:

* the name suggests subtype/refinement
* the content is a law

That creates cognitive friction.

---

# Why this matters for your architecture

You are building:

> relational IR + DSL layers + normalization

So you need to distinguish:

### 1. Ontology (what things are)

* `type`
* `rel`
* `def`

### 2. Refinement (subsets / subtypes)

* `specialization` (optional, but nice)

### 3. Truth statements (laws)

* `spec`

If you merge (2) and (3), later reasoning becomes messy.

---

# How specialization fits beautifully with your current design

You already have:

```len
def TotalFun(A, B) as f where
    Fun(A, B, f) and total(f)
```

You could equivalently express it as:

```len
specialization TotalFun(A, B, f)
given Fun(A, B, f)
must total(f)
```

So:

* `def` = *set-builder style*
* `specialization` = *predicate refinement style*

These are dual views — that’s actually powerful.

---

# Recommendation

Do this:

## Keep `spec` (laws)

## Introduce `specialization` (refinement)

### Example unified system

```len
type Set
rel has(S, x)

specialization NonEmptySet(S)
given Set(S)
must exists x. has(S, x)

spec union_contains_left(A, B, U)
given union(A, B, U)
must forall x. has(A, x) -> has(U, x)
```

---

# Naming alternatives (if you want to go further)

If you don’t like `spec`, better candidates for laws:

* `law`
* `theorem`
* `property`
* `rule`

For example:

```len
law union_contains_left(...)
```

That might align even better with your math-first direction.

---

# My concise take

* Your instinct is good: *“specialization” captures an important concept*
* But it should **not replace `spec`**
* It should be introduced as a **separate refinement construct**

---

# Strong recommendation in one line

Keep:

* `def` → derived objects
* `spec` (or `law`) → truths
* add `specialization` → refined subsets

That gives you a very clean and scalable semantic lattice, especially when you move into `l2` with category theory and programming DSLs.
