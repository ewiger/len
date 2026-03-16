Good approach. Before defining syntax, the **design questions must be made explicit**.  
A specification language like **`len`** sits in an unusual position between:

- natural language
    
- pseudocode
    
- formal logic
    
- verification languages
    

So the first step is to **identify the design constraints** and possible answers.

Below is a **working table of key questions**. The goal is not to answer them immediately but to narrow the design space.

---

# Design Questions for `len`

|#|Question|Why it matters|Possible options|
|---|---|---|---|
|1|Is `len` compiled or interpreted?|Determines whether `len` is just documentation or a formal intermediate language|(A) documentation only, (B) validated IR, (C) compiled to Lean/SMT|
|2|Is `len` executable?|If executable, it becomes a programming language rather than a spec language|(A) non-executable spec only, (B) partially executable, (C) fully executable|
|3|What is the core semantic foundation?|Needed to avoid ambiguity|(A) FOL semantics, (B) type theory semantics, (C) relational semantics|
|4|Does `len` require types?|Types strongly affect expressiveness and compilation|(A) untyped, (B) optional types, (C) strict typed|
|5|Are specifications relational or functional?|Major modeling decision|(A) functions only, (B) relations only, (C) relations with optional functions|
|6|Can specs be partial / underspecified?|Critical for early-stage requirements|(A) no, must be precise, (B) yes with "open" slots|
|7|Should ambiguity be representable?|Key if the language is produced from NL|(A) disallow ambiguity, (B) explicit ambiguity markers|
|8|Does `len` allow state / mutation?|Determines if imperative systems can be specified|(A) pure relational, (B) state via pre/post conditions, (C) full state logic|
|9|What is the minimal logic supported?|Determines reasoning capability|(A) propositional, (B) first-order logic, (C) higher-order|
|10|Can `len` express recursion and induction?|Needed for many program specs|(A) no, (B) recursive definitions allowed|
|11|Should `len` be deterministic after normalization?|Important for equivalence and reasoning|(A) yes canonical core form, (B) multiple representations allowed|
|12|Should the language support modules / namespaces?|Important for scaling knowledge bases|(A) flat namespace, (B) modules|
|13|Can `len` express algorithms or only specs?|Distinguishes pseudocode from spec language|(A) specs only, (B) optional algorithm blocks|
|14|Should there be a small trusted core?|Needed for verification soundness|(A) no formal core, (B) small logical kernel|
|15|Should `len` support proof hints / lemmas?|Useful for proof assistants|(A) none, (B) optional hints|
|16|Should normalization remove syntactic sugar?|Important for equivalence checking|(A) yes canonical IR, (B) preserve surface structure|
|17|Can `len` reference external theories?|Needed for math libraries and domain knowledge|(A) self-contained, (B) importable theories|
|18|Should `len` be machine-first or human-first?|Affects syntax style|(A) logic-first syntax, (B) human-readable DSL|
|19|Should every `len` expression map to a logical formula?|Guarantees semantic grounding|(A) yes strict mapping, (B) partially informal|
|20|What is the unit of knowledge?|Defines how memory / KB evolves|(A) facts, (B) specs, (C) theories|

---

# First Observations

A few of these decisions are **almost forced** by your earlier architecture.

### 1. `len` must be compilable

Otherwise the pipeline

```
Human → len → Lean
```

breaks.

So the most consistent choice is:

**`len` is a validated IR that compiles to formal logic.**

---

### 2. `len` should not be a full programming language

If it becomes executable code, you recreate:

- Python
    
- Rust
    
- ML
    

The real goal is **semantic specification**.

So likely:

**`len` = relational spec language**

---

### 3. `len` must support partial specifications

Natural language requirements are incomplete.

Example:

```
sort(xs)
```

Open questions:

- stable?
    
- ordering?
    
- duplicates?
    

Therefore:

**underspecification must be representable.**

---

### 4. The semantic foundation should be relational

Functions alone are too restrictive.

Instead:

```
Spec(input, output)
```

is more flexible.

---

### 5. Canonical normalization is critical

You previously mentioned **semantic equivalence**.

To reason about equivalence, we likely need:

```
surface syntax
→ normalization
→ canonical core
```

---

# Likely Core Shape of `len`

From the constraints above, a plausible core emerges.

```
len surface syntax
        ↓
normalization
        ↓
len-core
        ↓
logic compiler
        ↓
Lean / SMT
```

So the **true language** is actually `len-core`.

Everything else is sugar.

---

# Key Strategic Question

The most important question for the next step is this:

### What is `len-core`?

Two realistic options exist.

---

## Option A — FOL-style core

Example:

```
forall x in xs .
    Member(x, ys) -> P(x)
```

Pros

- simple
    
- maps well to SMT
    
- well understood
    

Cons

- weaker abstraction
    
- less expressive
    

---

## Option B — Type theory style core

Example:

```
∀ x : Int, x ∈ xs → P x
```

Pros

- expressive
    
- matches Lean
    

Cons

- harder to generate automatically
    

---

# Another Critical Question

This one determines whether `len` becomes truly useful.

### Should `len` allow _unknown semantics_?

Example:

```
open OrderingRelation
```

This is extremely powerful.

It allows the spec to exist **before all decisions are fixed**.

Most formal languages cannot do this.

---

# A Good Next Step

The most productive next move is:

1. Reduce the table to **~8 core constraints**
    
2. Decide them explicitly
    
3. Derive `len-core` from those decisions
    

This prevents the language from becoming incoherent.

---

# My Suggested Next Question

Among all these, **three decisions determine almost everything**.

We should settle them next.

### Q1

Is `len` fundamentally

- **logic with sugar**  
    or
    
- **structured natural language compiled to logic**?
    

---

### Q2

Is the semantic core

- **first-order relational logic**  
    or
    
- **type-theoretic logic (Lean-like)**?
    

---

### Q3

Is underspecification (open slots) a **first-class feature**?

---

Once these three are fixed, the rest of the design becomes much clearer.

If you want, the next step can be something extremely useful:

**derive the smallest possible `len-core` calculus (maybe ~8 constructs total).**

That exercise usually exposes hidden design contradictions immediately.