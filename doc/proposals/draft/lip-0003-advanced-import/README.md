# LIP-0003: Instantiated Contract Namespaces and Alias-Style Member Imports

## Metadata

| Field | Value |
| --- | --- |
| LIP | 0003 |
| Title | Instantiated Contract Namespaces and Alias-Style Member Imports |
| Status | Draft |
| Author | TBD |
| Created | 2026-03-22 |
| Updated | 2026-03-22 |
| Area | Imports, Namespaces, Formula Surface |
| Target | len.l1 |

## Summary

This proposal defines how members of a `contract` are referenced after contract instantiation and how those members should be brought into local scope.

The key rule is that an instantiated contract proposition such as `Eq(Point)` acts as a namespace container for its members. If `contract Eq(T: Type)` declares `rel Equal(x: T, y: T)`, then `Eq(Point).Equal` is the fully qualified member relation specialized to `Point`.

That fully qualified form is semantically valid, but it is not the preferred formula surface. In `spec` and other formula contexts, dotted member predicates such as `Eq(Point).Equal(p, p)` are treated as an anti-pattern. The preferred style is to import the member into a local alias and then use the local name directly.

This proposal therefore extends `import` so that it can bind a named member from a module or contract namespace into a local identifier. The goal is to keep formulas clean while still treating contracts as real namespace boundaries.

## Motivation

LIP-0002 already leans toward treating `contract` as a namespace boundary. Once that decision is taken seriously, contract members need a clear source-level access path.

The natural fully qualified form is easy to explain:

- `Eq(Point)` is the contract proposition for the carrier type
- `Eq(Point).Equal` is the namespaced relation member specialized to `Point`
- `Eq(Point).Equal(p, p)` is a formula-level use of that specialized member relation

However, the last form is too heavy for ordinary proofs and specifications. It mixes namespace navigation with logical content, makes formulas harder to scan, and encourages a member-access style that looks more object-oriented than the language intends.

The language needs both of these properties at once:

- contract instantiations must remain real namespace containers
- formulas should still read in the flat relational style used elsewhere in `len`

Alias-style member imports provide that split cleanly.

## Problem Statement

Without an explicit import rule for contract members, users are pushed toward one of two bad outcomes.

First, they may write fully qualified member predicates directly inside formulas.

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)

spec point_eq_reflexive
    given p: Point
    must Eq(Point).Equal(p, p)
```

This is mechanically understandable but stylistically poor. The logical statement is about equality on points, but the surface form is dominated by qualification syntax.

Second, users may want to import the dotted path directly as if it were a single qualified name.

```len
import Eq(Point).Equal from core.math.logic.eq
```

That is also the wrong model. `Eq(Point)` is not a prefix fragment in a long dotted identifier. It is the instantiated contract proposition, and it should be treated as the namespace container from which members are selected.

If the language leaves this area vague, tooling and human style will drift apart. Some code will flatten contract members into ad hoc top-level names, some code will overuse dotted formulas, and imports will become ambiguous.

## Proposed Changes

### 1. Instantiated Contracts Are Namespace Containers

For a contract declaration such as:

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)
```

the instantiated proposition `Eq(Point)` is both:

- the contract proposition for the carrier type `Point`
- the namespace container for the contract members specialized to `Point`

So the canonical member reference is:

```len
Eq(Point).Equal
```

The same rule applies to other contract-local declarations such as `fn` and `spec` members.

### 2. Dotted Member Predicates Are Allowed but Discouraged in Formulas

The following formula is valid in meaning:

```len
Eq(Point).Equal(p, p)
```

It means: apply the `Equal` relation member from the instantiated contract namespace `Eq(Point)` to the arguments `p` and `p`.

But this proposal treats that surface form as a non-idiomatic fallback. It may still appear in explanations, diagnostics, elaboration output, or intermediate lowering descriptions, but it should not be the preferred author-facing formula style.

The preferred surface is to import the member relation into a local name and use that local name in formulas.

### 3. `import` Binds Members from a Namespace into a Local Alias

Add a member-import form:

```len
import Equal from core.math.logic.eq.Eq(Point) as Equal
```

This means:

- resolve the module path `core.math.logic.eq`
- resolve the instantiated contract namespace `Eq(Point)` within that module
- select the member `Equal` from that namespace
- bind it into the current module under the local alias `Equal`

After that import, formulas should use the local alias:

```len
spec point_eq_reflexive
    given p: Point
    must Equal(p, p)
```

The shorter form without renaming may be accepted as surface sugar:

```len
import Equal from core.math.logic.eq.Eq(Point)
```

but its meaning is still alias binding, equivalent to importing `Equal` as the same local name.

### 4. Do Not Import Fully Qualified Dotted Member Paths Directly

The following style should be rejected or, at minimum, documented as invalid surface:

```len
import Eq(Point).Equal from core.math.logic.eq
```

Reason:

- `Eq(Point)` is the namespace container
- `Equal` is the imported member
- `import` should make that split explicit instead of treating the whole expression as a dotted path token

This keeps namespace structure visible and avoids inventing a special parser rule for importing an application-shaped prefix.

### 5. Exported Names Must Not Conflict Across Import-Relevant Kinds

Because imported members are bound into a local alias, the exported surface of a module or contract namespace must avoid ambiguous identifier collisions.

The minimum rule introduced by this proposal is:

- a module or contract namespace must not export the same identifier as both a `type` and a `rel`

That disjointness is enough to support clear lookup for the cases this proposal introduces. Tooling may later widen this rule to include `fn` and `const` as the import surface becomes richer, but this proposal commits at least to the `type` versus `rel` separation.

## Grammar Sketch

Illustrative grammar only:

```ebnf
ImportDecl        = ModuleImportDecl | MemberImportDecl ;
ModuleImportDecl  = "import" ModulePath ;
MemberImportDecl  = "import" Identifier "from" NamespaceRef [ "as" Identifier ] ;

NamespaceRef      = ModulePath "." NamespaceTail ;
NamespaceTail     = Identifier
                  | ContractInstance ;
ContractInstance  = Identifier "(" [ ArgList ] ")" ;
ArgList           = Expr { "," Expr } ;

QualifiedMember   = ContractInstance "." Identifier ;
```

Notes:

- `QualifiedMember` is the explanatory source form for contract-member selection
- `MemberImportDecl` imports a member from a namespace reference, not an already-dotted member path
- the exact split between `ModulePath` and post-module namespace resolution may be parser-specific

## Desugaring / Lowering Model

Surface contract:

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)
```

Source-level qualified use:

```len
Eq(Point).Equal(p, p)
```

Preferred imported use:

```len
import Equal from core.math.logic.eq.Eq(Point) as Equal

spec point_eq_reflexive
    given p: Point
    must Equal(p, p)
```

Conceptually, the import introduces a local binding equivalent to:

- local name `Equal`
- target member `core.math.logic.eq.Eq(Point).Equal`

The formula then typechecks and validates against that local binding rather than repeatedly spelling the qualified member expression.

This proposal does not require runtime objects, instance search, or method dispatch. The namespace role is purely structural and static.

## Examples

### Before / After

Before, with direct dotted formula use:

```len
spec point_eq_reflexive
    given p: Point
    must Eq(Point).Equal(p, p)
```

After, with local alias import:

```len
import Equal from core.math.logic.eq.Eq(Point) as Equal

spec point_eq_reflexive
    given p: Point
    must Equal(p, p)
```

### Canonical Example

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)

import Equal from core.math.logic.eq.Eq(Point) as Equal

spec point_eq_reflexive
    given p: Point
    must Equal(p, p)
```

Reading guide:

- `Eq(Point)` is the contract proposition for `Point`
- `Eq(Point).Equal` is the namespaced member specialized to `Point`
- `Equal` is the local alias imported from that namespace
- `Equal(p, p)` is the preferred formula-level use

## Rationale

This proposal fits `len` because it preserves the existing relational reading instead of turning contract members into methods.

- `type` remains the core carrier notion
- `rel` remains the core predicate notion
- `spec` remains the place where laws are stated
- `fn` remains the executable or constructive form when needed

The proposal adds only source-level namespace and import structure. It does not add hidden receivers, instance resolution, object fields, or runtime semantics.

Treating `Eq(Point)` as a namespace container is also the cleanest explanation of specialized contract members. It says exactly where the member comes from while still allowing formulas to stay visually flat through local alias imports.

## Backward-Compatibility Notes

- current MVP grammar only supports `import ModulePath`; this proposal extends import syntax beyond the accepted LIP-0001 parser scope
- validator and loader work will need to distinguish module imports from member imports
- name-resolution rules will need to track alias bindings for imported members
- duplicate-name checks will need at least one new disjointness rule for exported `type` and `rel` identifiers in the same importable namespace

## Open Questions

1. Should `as` be mandatory for member imports, or should omitting it remain sugar for importing under the same local name?
2. Should the anti-pattern status of dotted member predicates be enforced by the parser, by style linting, or only by documentation?
3. Should member import support extend immediately to contract-local `fn` and `spec`, or should the first implementation cover `rel` only?
4. How should nested namespace cases be handled if contracts are later allowed inside richer module structures?

## Non-Goals

- object-oriented member access semantics
- method dispatch or instance lookup
- wildcard imports from contracts
- implicit opening of all contract members into scope
- changing the meaning of `Eq(Point)` as a proposition

## Future Proposal Notes

- this proposal should be read together with `LIP-0002`, especially `contract.md`
- if accepted, a follow-up implementation plan may be added once parser and validator work for advanced imports is scheduled