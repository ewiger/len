# len.l1

`len.l1` is the level 1 layer of the `len` meta-language for expressing domain models, contracts, and specifications. Its semantic core remains small: `type`, `rel`, `const`, and `spec` carry the language, while `contract`, `struct`, `fn`, `syntax`, and `quasi` organize that core at the surface level.

## Core

Core `len.l1` is organized around a contract-centric surface model:

- **core.syntax** = the reflective vocabulary for expressions, identifiers, symbols, and syntax-level aliases
- **core.math.logic.logic** = the base logical carriers and defining laws for formulas and boolean constants
- **core.struct** = contract-oriented surface vocabulary for structured data

The current corpus treats `contract` as a central L1 surface form alongside `type`, `rel`, and `spec`. Older vocabulary based on `trait` and `impl` is no longer part of the language corpus.

Then it is extended by various domain-specific libraries, e.g. `core.math.set` for set theory and `core.math.nat` for natural numbers.

### Syntax 

The syntax of `len.l1` is intentionally reflective rather than parser-complete. The core syntax inventory reserves the current surface forms and supplies meta-level carriers for expression-like entities.

Within that reflective layer:

- `contract` is a first-class keyword used to group related `rel`, `fn`, and `spec` declarations
- `const` is part of the canonical surface and is used for named constants such as `True` and `False`
- syntax aliases are defined explicitly through equality, with `syntax x = y` treated as surface sugar for `x equals y`
- the named equality target for that sugar is the reflective relation `Equals(x, y)` in `core.syntax.syntax`

### Status of Legacy Surface Forms

The keyword inventory still reserves `fn`, `requires`, `ensures`, and `implements`, but they are not the organizing backbone of the current core corpus.

- `fn` remains available for executable or constructive surface forms
- `requires` and `ensures` remain available for function-level specification clauses
- `implements` remains associated with linking a `fn` declaration to a relation when that surface is used

By contrast, `trait`, `impl`, and `note` are not part of the current `len.l1` corpus.