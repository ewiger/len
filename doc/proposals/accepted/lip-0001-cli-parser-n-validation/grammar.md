# LIP-0001 Grammar Notes for len.l1

## Purpose

This document defines the MVP lexer and parser input for `len.l1` as accepted by LIP-0001.

The grammar here is intentionally scoped to the current repository corpus under `lang/l1/**`.

## Lexical Rules

### Whitespace

Whitespace separates tokens and is otherwise ignored except where it helps terminate identifiers, operators, and comments.

### Comments

Comments are allowed anywhere and ignored by formal semantics.

Supported comment forms:

```text
# line comment

<#
block comment
#>
```

Rules:

- `#` starts a line comment only when it appears at the beginning of the line after optional indentation
- `<#` starts a block comment
- `#>` ends a block comment
- block comments are not part of formal semantics

### Docstrings

Docstrings are allowed anywhere and ignored by formal semantics.

Docstrings are preferred over line comments for structured documentation.

Supported docstring forms:

```text
"single-line docstring"

"""
Multi-line docstring.
Ignored by formal semantics.
Useful for intent, examples, and notes.
"""
```

Rules:

- a standalone string literal may be treated as a docstring token or ignored trivia
- triple-quoted strings are multi-line docstrings
- docstrings must not affect parsing or validation outcome
- parser implementations may preserve docstrings in the AST for tooling, but validation must ignore them

## Reserved Words

The lexer should recognize at least the current observed reserved words:

- `type`
- `rel`
- `syntax`
- `keyword`
- `symbol`
- `implies`
- `import`
- `from`
- `spec`
- `given`
- `must`
- `forall`
- `exists`
- `true`
- `false`
- `and`
- `or`
- `not`
- `iff`
- `where`
- `in`
- `struct`
- `trait`
- `impl`
- `quasi`
- `const`
- `case`
- `when`
- `as`

## Symbolic Tokens

The lexer should recognize at least these symbolic forms:

- `.`
- `:`
- `(`
- `)`
- `,`
- `=`
- `<`
- `>`
- `|`
- `/`
- `\`
- `*`
- `+`
- `-`
- `!`
- `=>`
- `<=>`
- `/\`
- `\/`
- `<#`
- `#>`

## Identifiers and Paths

- `Identifier` is an unqualified name used for declarations and local binders
- `QualifiedName` is a dotted sequence of identifiers
- `ModulePath` is a dotted sequence of identifiers used by `import`

## Parser Overview

The parser should treat comments and docstrings as non-semantic trivia.

The supported MVP top-level forms are:

- `import`
- `type`
- `rel`
- `const`
- `spec`
- `syntax`
- `trait`
- `impl`
- `keyword`
- `symbol`
- `quasi`

## EBNF

```ebnf
Program        = { Trivia | TopLevelDecl } ;

Trivia         = Comment | Docstring ;

TopLevelDecl   = ImportDecl
               | TypeDecl
               | RelDecl
               | ConstDecl
               | SpecDecl
               | SyntaxDecl
               | TraitDecl
               | ImplDecl
               | KeywordDecl
               | SymbolDecl
               | QuasiDecl ;

ImportDecl     = "import" ModulePath ;
TypeDecl       = "type" Identifier ;
RelDecl        = "rel" Identifier "(" [ ParamList ] ")" ;
ConstDecl      = "const" Identifier ":" TypeExpr ;
TraitDecl      = "trait" Identifier ;
ImplDecl       = "impl" TypeExpr ":" Identifier ;
KeywordDecl    = "keyword" Identifier ;
SymbolDecl     = "symbol" Identifier "as" StringLiteral ;

SpecDecl       = "spec" Identifier { Trivia | GivenClause } MustClause ;
GivenClause    = "given" BinderList ;
MustClause     = "must" Formula ;

SyntaxDecl     = "syntax" Expr "where" BinderList "implies" Expr ;

QuasiDecl      = "quasi" Identifier { Trivia | QuasiCase } ;
QuasiCase      = "case" Expr { Trivia | QuasiRule } ;
QuasiRule      = Expr [ "when" Formula ] ;

BinderList     = Binder { "," Binder } ;
Binder         = Identifier ":" TypeExpr ;
ParamList      = Binder { "," Binder } ;
ModulePath     = Identifier { "." Identifier } ;

Formula        = Quantified
               | BinaryFormula
               | UnaryFormula
               | Atom ;

Quantified     = Quantifier BinderList "." Formula ;
Quantifier     = "forall" | "exists" ;

BinaryFormula  = Expr BinaryOp Expr ;
UnaryFormula   = UnaryOp Formula ;
UnaryOp        = "not" | "!" ;
BinaryOp       = "="
               | "in"
               | "subsetof"
               | "and"
               | "/\\"
               | "or"
               | "\\/"
               | "implies"
               | "=>"
               | "iff"
               | "<=>" ;

Expr           = Application | Primary ;
Application    = Identifier "(" [ ExprList ] ")" ;
ExprList       = Expr { "," Expr } ;

Primary        = Identifier
               | QualifiedName
               | StringLiteral
               | "(" Formula ")" ;

TypeExpr       = Identifier | QualifiedName ;

Comment        = LineComment | BlockComment ;
LineComment    = "#" { Character } ;
BlockComment   = "<#" { Character } "#>" ;
Docstring      = StringLiteral | MultiLineDocstring ;
```

## Operator Precedence

Recommended precedence from tightest to loosest:

1. grouping and application
2. unary `not` and `!`
3. equality and infix relations such as `in` and `subsetof`
4. `and` and `/\`
5. `or` and `\/`
6. `implies` and `=>`
7. `iff` and `<=>`

## Semantic Notes

- comments and docstrings are ignored by formal semantics
- `syntax` declarations do not make the parser dynamically extensible in the MVP
- semantic validation runs after parse and handles arity, symbol resolution, and binder scope

## Known MVP Limits

- grammar is intentionally fixed to the current corpus under `lang/l1/**`
- import path resolution is defined by the implementation plan, not by this grammar alone
- comments and docstrings may be retained in the AST for tooling, but they must not influence validation