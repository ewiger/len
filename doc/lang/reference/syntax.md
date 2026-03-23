# Syntax Reference

This page covers the complete lexical and grammatical rules for `len.l1`.

---

## File Extension

Level 1 source files use the `.l1` extension.

---

## Whitespace and Indentation

Whitespace separates tokens and is otherwise ignored, **except** for indentation. `len.l1` is indentation-sensitive:

- The body of a `spec`, `fn`, `contract`, `struct`, and `quasi` block is indented relative to its header.
- A dedent (return to an earlier column) closes the block.
- Tabs and spaces may not be mixed within the same file.

---

## Comments

```text
# This is a line comment — starts at the beginning of a line after optional indentation

<#
  This is a block comment.
  It may span multiple lines.
#>
```

- `#` begins a line comment only when it appears at the start of a line (after optional whitespace).
- `<# … #>` is a block comment and may appear anywhere.
- Comments are ignored by the parser and have no effect on semantics.

---

## Docstrings

Docstrings document declarations. They are preserved in the AST for tooling but ignored by validation.

```len
"This is a single-line docstring."

"""
This is a multi-line docstring.
It may span as many lines as needed.
"""
```

Place a docstring immediately before or after the declaration it documents.

---

## Identifiers

An identifier begins with a letter or underscore and continues with letters, digits, and underscores.

```
Identifier = [A-Za-z_][A-Za-z0-9_]*
```

Identifiers are case-sensitive. By convention:

| Form | Used for |
|------|----------|
| `PascalCase` | Types, relations, contracts, structs, constants |
| `snake_case` | Functions, specs, local binders |
| `dot.separated` | Module paths in `import` |

---

## Literals

| Kind | Example | Notes |
|------|---------|-------|
| Integer | `42`, `0`, `100` | Positive integers |
| String | `"hello"` | Double-quoted; used in docstrings and quasi bodies |
| Boolean | `true`, `false` | Reserved keywords; map to `True`/`False` constants |

---

## Operators and Symbolic Tokens

The following symbolic forms are recognised by the lexer:

| Symbol | Role |
|--------|------|
| `.` | Module path separator; field access |
| `:` | Type annotation separator |
| `(` `)` | Parameter lists, formula grouping |
| `,` | Parameter separator |
| `=` | Equality; used in `syntax` sugar and field assignment |
| `<` `>` | Comparison; block comment delimiters (`<#`, `#>`) |
| `\|` | Cardinality notation (`\|s\|`) |
| `/` | Formula division |
| `\` | Backslash in formula notation |
| `*` | Multiplication |
| `+` | Addition |
| `..` | Range in quasi blocks (`0 .. n`) |

---

## Sugar Operators (core.math.logic.sugar)

Import `core.math.logic.sugar` to enable these infix and prefix operators:

| Operator | Desugars to | Meaning |
|----------|------------|---------|
| `A => B` | `A implies B` | Implication |
| `A <=> B` | `A iff B` | Biconditional |
| `A /\ B` | `A and B` | Conjunction |
| `A \/ B` | `A or B` | Disjunction |
| `!A` | `not A` | Negation |

Example:

```len
import core.math.logic.sugar

spec de_morgan
    given A, B: Formula
    must !(A /\ B) <=> (!A \/ !B)
```

---

## Top-Level Declarations

A `.l1` file is a sequence of top-level declarations. Each declaration starts at column 0.

### Type Declaration

```ebnf
TypeDecl = "type" Identifier ;
```

```len
type Point
type Seq
```

### Relation Declaration

```ebnf
RelDecl = "rel" Identifier "(" ParamList ")" [ "refines" QualifiedName ] ;
ParamList = Param { "," Param } ;
Param     = Identifier ":" QualifiedName ;
```

```len
rel Member(x: Set, s: Set)
rel BubbleSort(input: Seq, output: Seq) refines Sort
```

### Constant Declaration

```ebnf
ConstDecl = "const" Identifier ":" QualifiedName ;
```

```len
const True  : Bool
const Pi    : Real
```

### Function Declaration

```ebnf
FnDecl     = "fn" Identifier "(" ParamList ")" "->" Param
             { FnClause }
             QuasiBlock ;
FnClause   = ImplementsClause | RequiresClause | EnsuresClause ;
ImplementsClause = "implements" Application ;
RequiresClause   = "requires" Formula ;
EnsuresClause    = "ensures" Formula ;
```

```len
fn sort(input: Seq) -> output: Seq
    implements Sort(input, output)
    requires NonEmpty(input)
    ensures Sorted(output)
    ensures Permutation(input, output)
    quasi:
        # pseudocode
```

### Spec Declaration

```ebnf
SpecDecl  = "spec" Identifier INDENT { GivenClause } MustClause DEDENT ;
GivenClause = "given" Param { "," Param } ;
MustClause  = "must" Formula ;
```

```len
spec sort_correct
    given input, output: Seq
    must Sort(input, output) iff
        Sorted(output) and Permutation(input, output)
```

### Contract Declaration

```ebnf
ContractDecl = "contract" Identifier [ "(" ParamList ")" ]
               INDENT { ContractMember } DEDENT ;
ContractMember = RelDecl | FnDecl | SpecDecl ;
```

```len
contract Ord(T: Type)
    rel Leq(x: T, y: T)
    spec ord_reflexive
        given x: T
        must Leq(x, x)
```

### Struct Declaration

```ebnf
StructDecl = "struct" Identifier INDENT { FieldDecl } DEDENT ;
FieldDecl  = Identifier ":" QualifiedName ;
```

```len
struct Address
    street  : String
    city    : String
    country : String
```

### Syntax Declaration

```ebnf
SyntaxDecl = "syntax" Pattern "where" BindingList "implies" QualifiedApplication ;
Pattern    = token sequence ;
BindingList = Param { "," Param } ;
```

```len
syntax x in s where x: Set, s: Set implies Member(x, s)
syntax A => B where A: Formula, B: Formula implies A implies B
```

### Import Declaration

```ebnf
ImportDecl   = "import" ModulePath ;
FromImport   = "from" ModulePath "import" Identifier { "," Identifier } ;
ModulePath   = Identifier { "." Identifier } ;
```

```len
import core.math.set
from core.math.set import Set, Member
```

---

## Formulas

A formula is a logical expression appearing in `spec` bodies, `fn` clauses, and `syntax` targets.

| Form | Example |
|------|---------|
| Relation application | `Member(x, s)` |
| Equality | `x = y` |
| Conjunction | `P and Q` |
| Disjunction | `P or Q` |
| Negation | `not P` |
| Implication | `P implies Q` |
| Biconditional | `P iff Q` |
| Universal | `forall x: T . P` |
| Existential | `exists x: T . P` |
| Syntax sugar | `x in s`, `A => B` |

Precedence (highest to lowest): `not` > `and` > `or` > `implies` > `iff`.

Use parentheses to override precedence.

---

## Quasi Blocks

A `quasi` block is an indented pseudocode body attached to an `fn` declaration.

```ebnf
QuasiBlock = "quasi" [ "using" "style" Identifier ] ":" INDENT RawLines DEDENT ;
```

The body is captured as raw lines. Allowed constructs depend on the active style profile.

**ProceduralAlgorithm** style (default) allows:

| Keyword | Meaning |
|---------|---------|
| `let <var> := <expr>` | Introduce or bind a variable |
| `set <var> := <expr>` | Reassign a variable |
| `return <expr>` | Return a value |
| `if <cond>:` | Conditional branch |
| `else:` | Else branch |
| `while <cond>:` | Loop |
| `for <var> in <range>:` | Iteration |
| `append(<list>, <val>)` | Append to a list |

These keywords are **quasi-local** — they are not part of the global `len.l1` grammar.

---

## Qualified Names

Declarations from other modules are referenced by their qualified name:

```len
import core.math.set

rel Example(s: core.math.set.Set)
```

Or after importing:

```len
import core.math.set

rel Example(s: Set)
```

Identifiers resolve in this order: local scope → imported modules → error.
