# len.l1

`len.l1` is the level 1 layer of a meta-language called `len` for expressing design and specification of domain knowledge in programming, mathematics, etc. This is done by means of relational logic with primitives like `type` and `rel`. 

## Core

Core len.l1 consists of two main parts:

- **core.syntax** = meta-language tokens plus core carriers such as `Expr`, `Type`, `Relation`, and `Formula`

- **core.math.logic.logic** = human-facing logical aliases built on top of len.l1 core connectives

Then it is extended by various domain-specific libraries, e.g. `core.math.set` for set theory and `core.math.nat` for natural numbers.

### Syntax 

The syntax of len.l1 is designed to be minimal and flexible, allowing for the expression of a wide range of concepts and relationships. 

However our module here deals mostly with reflection rather than building a parser, so we have a very minimal syntax that includes:

- a reflection vocabulary
- a meta-level ontology of surface forms
- a place to say “there are identifiers, formulas, qualified names”