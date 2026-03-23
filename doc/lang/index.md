# len

<div class="hero" markdown>

## A specification-first language for programs about programs

`len` is a meta-programming language for writing systems that **generate, transform, or reason about code structures**.

Instead of treating the specification as a comment written after the fact, `len` makes it the primary artifact. You state the model, the relations, and the laws first, then connect executable forms to those declarations.

[Read the concepts guide](concepts.md){ .md-button .md-button--primary }
[Start with Hello World](tutorials/01-hello-world.md){ .md-button }

</div>

<div class="grid cards" markdown>

-   :material-compass-outline: __Orient yourself__

    Understand the language model, module system, and layer boundaries.

    [Open Concepts](concepts.md)

-   :material-play-circle-outline: __Build and validate__

    Install the toolchain and run your first `.l1` file through the validator.

    [Open the Hello World tutorial](tutorials/01-hello-world.md)

-   :material-book-search-outline: __Look up syntax__

    Jump straight into reserved words, grammar, and language surface details.

    [Open the reference](reference/keywords.md)

-   :material-console-line: __Use the CLI__

    See command syntax, exit codes, diagnostics, and validation behavior.

    [Open the CLI reference](cli.md)

</div>

## Where to Start

| Goal | Start here |
|------|-----------|
| Understand what len is and why it exists | [Concepts](concepts.md) |
| Install and run your first program | [Hello World tutorial](tutorials/01-hello-world.md) |
| Look up a specific keyword | [Keywords reference](reference/keywords.md) |
| Use the command-line tool | [CLI reference](cli.md) |
| Come from a Java background | [Coming from Java](tutorials/06-from-java.md) |
| Generate a PostgreSQL schema | [Postgres shop tutorial](tutorials/07-postgres-schema-generation.md) |

## Language Layers

len is organised into three levels:

| Level | Name | Purpose |
|-------|------|---------|
| **l0** | Natural | Informal pseudocode, design notes, brainstorming |
| **l1** | Structural core | Types, relations, formulas, specifications |
| **l2** | Evaluation | Contexts, satisfaction, executable semantic rules |

Most day-to-day work happens in `l1`. Files carry the extension `.l1`.

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

## Documentation Map

### Tutorials

1. [Hello World](tutorials/01-hello-world.md) — install, build, validate
2. [Types and Relations](tutorials/02-types-and-relations.md) — core modelling primitives
3. [Writing Specifications](tutorials/03-writing-specifications.md) — axioms, definitions, and laws
4. [Functions and Quasi](tutorials/04-functions-and-quasi.md) — executable forms and pseudocode
5. [Contracts and Structs](tutorials/05-contracts-and-structs.md) — grouping and record types
6. [Coming from Java](tutorials/06-from-java.md) — map OOP ideas into len
7. [PostgreSQL Schema Generation](tutorials/07-postgres-schema-generation.md) — generate an online shop schema

### Reference

- [Concepts](concepts.md) — philosophy, layers, and the module system
- [Keywords](reference/keywords.md) — reserved words and their meanings
- [Syntax](reference/syntax.md) — grammar, operators, and sugar forms
- [CLI](cli.md) — `len-cli` command reference
