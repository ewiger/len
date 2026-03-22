# len — Meta Programming Language

`len` is a meta-programming language for writing programs that **generate, transform, or reason about other programs and code structures**.

It is specification-first: you describe *what* a program is and *what properties it must satisfy* before describing *how* it runs. This makes it well suited for domain modelling, formal contracts, and code generation targets such as SQL schemas, type-safe APIs, and parsers.

---

## Where to Start

| Goal | Start here |
|------|-----------|
| Understand what len is and why it exists | [Concepts](concepts.md) |
| Install and run your first program | [Hello World tutorial](tutorials/01-hello-world.md) |
| Look up a specific keyword | [Keywords reference](reference/keywords.md) |
| Use the command-line tool | [CLI reference](cli.md) |
| Come from a Java background | [Coming from Java](tutorials/06-from-java.md) |
| Generate a PostgreSQL schema | [Postgres shop tutorial](tutorials/07-postgres-schema-generation.md) |

---

## Language Layers

len is organised into three levels:

| Level | Name | Purpose |
|-------|------|---------|
| **l0** | Natural | Informal pseudocode, design notes, brainstorming |
| **l1** | Structural core | Types, relations, formulas, specifications — the heart of the language |
| **l2** | Evaluation | Contexts, satisfaction, executable semantic rules |

Most day-to-day work happens in `l1`. Files carry the extension `.l1`.

---

## Quick Taste

```len
import core.math.set

type Point

rel OnCircle(p: Point, radius: Nat)

spec circle_definition
    given p: Point
    given r: Nat
    must OnCircle(p, r) iff
        exists x: Nat, y: Nat .
            At(p, x, y) and x * x + y * y = r * r
```

---

## Tutorials

1. [Hello World](tutorials/01-hello-world.md) — install, build, validate
2. [Types and Relations](tutorials/02-types-and-relations.md) — the core modeling primitives
3. [Writing Specifications](tutorials/03-writing-specifications.md) — axioms, definitions, laws
4. [Functions and Quasi](tutorials/04-functions-and-quasi.md) — executable forms and pseudocode
5. [Contracts and Structs](tutorials/05-contracts-and-structs.md) — grouping and record types
6. [Coming from Java](tutorials/06-from-java.md) — OOP mental model mapped to len
7. [PostgreSQL Schema Generation](tutorials/07-postgres-schema-generation.md) — code generation for an online shop

---

## Reference

- [Concepts](concepts.md) — philosophy, layers, module system
- [Keywords](reference/keywords.md) — every reserved word explained
- [Syntax](reference/syntax.md) — grammar, operators, sugar forms
- [CLI](cli.md) — `len-cli` command reference
