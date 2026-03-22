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
- `fn`
- `syntax`
- `keyword`
- `symbol`
- `implies`
- `import`
- `from`
- `spec`
- `given`
- `must`
- `requires`
- `ensures`
- `implements`
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
- `as`

Notes:

- `fn` is a top-level grammar form for function-like declarations
- `quasi` is a dedicated embedded block for pseudo-algorithmic syntax inside `fn` declarations
- `quasi` is not a separate level of the language and does not require a separate top-level parser; the main parser should parse it via a dedicated block routine
- words used only inside quasi statements such as `let`, `choose`, `return`, `if`, `else`, `while`, and `repeat` should remain quasi-local rather than become global `len.l1` keywords

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
- `->`
- `:=`
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
- `fn`
- `const`
- `spec`
- `syntax`
- `trait`
- `impl`
- `keyword`
- `symbol`

Preferred direction after the current MVP:

- keep quasi as an embedded `quasi:` block owned only by `fn`
- let the host parser capture a `QUASI_BLOCK`
- parse the block contents with a dedicated quasi parsing routine inside the main parser

## EBNF

```ebnf
Program        = { Trivia | TopLevelDecl } ;

Trivia         = Comment | Docstring ;

TopLevelDecl   = ImportDecl
               | TypeDecl
               | RelDecl
               | FnDecl
               | ConstDecl
               | SpecDecl
               | SyntaxDecl
               | TraitDecl
               | ImplDecl
               | KeywordDecl
               | SymbolDecl ;

ImportDecl     = "import" ModulePath ;
TypeDecl       = "type" Identifier ;
RelDecl        = "rel" Identifier "(" [ ParamList ] ")" ;
FnDecl         = "fn" Identifier Signature FnBody ;
ConstDecl      = "const" Identifier ":" TypeExpr ;
TraitDecl      = "trait" Identifier ;
ImplDecl       = "impl" TypeExpr ":" Identifier ;
KeywordDecl    = "keyword" Identifier ;
SymbolDecl     = "symbol" Identifier "as" StringLiteral ;

SpecDecl       = "spec" Identifier { Trivia | GivenClause } MustClause ;
GivenClause    = "given" BinderList ;
MustClause     = "must" Formula ;

Signature      = "(" [ ParamList ] ")" [ "->" ResultBinder ] ;
ResultBinder   = Identifier ":" TypeExpr ;
FnBody         = INDENT { Trivia | FnClause } DEDENT ;
FnClause       = RequiresClause
               | EnsuresClause
               | ImplementsClause
               | QuasiClause ;
RequiresClause = "requires" Formula ;
EnsuresClause  = "ensures" Formula ;
ImplementsClause = "implements" Formula ;
QuasiClause    = "quasi" ":" QuasiBlock ;
QuasiBlock     = INDENT { QuasiLine } DEDENT ;

SyntaxDecl     = "syntax" Expr "where" BinderList "implies" Expr ;

QuasiLine      = QuasiStmt | QuasiNarrative ;
QuasiStmt      = LetStmt
               | SetStmt
               | ChooseStmt
               | AssumeStmt
               | HaveStmt
               | ShowStmt
               | ReturnStmt
               | ForEachStmt
               | IfStmt
               | WhileStmt
               | RepeatStmt ;
LetStmt        = "let" Identifier [ ":" TypeExpr ] ":=" Expr ;
SetStmt        = "set" Identifier ":=" Expr ;
ChooseStmt     = "choose" Identifier [ ":" TypeExpr ] "such" "that" Formula ;
AssumeStmt     = "assume" Formula ;
HaveStmt       = "have" Formula ;
ShowStmt       = "show" Formula ;
ReturnStmt     = "return" Expr ;
ForEachStmt    = "for" "each" Identifier "in" Expr ":" QuasiBlock ;
IfStmt         = "if" Formula ":" QuasiBlock [ "else" ":" QuasiBlock ] ;
WhileStmt      = "while" Formula ":" QuasiBlock ;
RepeatStmt     = "repeat" ":" QuasiBlock "until" Formula ;
QuasiNarrative = TextLine ;

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
- `fn` owns signature-level contract clauses and may include a `quasi:` implementation sketch
- `spec` remains a declarative statement form built from `given` and `must`
- quasi is parsed as an embedded block language by the main parser rather than as a top-level declaration form
- `implements` is the preferred clause for linking an `fn` to an abstract relation, while `ensures` remains appropriate for postconditions

Preferred future extension for configurable quasi styles:

- allow `quasi using CustomStyle:` as the clause header
- interpret `CustomStyle` as a style profile that fixes the accepted step keywords and line schemas for that block
- allow style profiles to come from built-in presets, repository-local declarations, or plugins
- Also `quasi:` means default style, which is currently `ProcStyle` but may be redefined in the future

## Known MVP Limits

- grammar is intentionally fixed to the current corpus under `lang/l1/**`
- import path resolution is defined by the implementation plan, not by this grammar alone
- comments and docstrings may be retained in the AST for tooling, but they must not influence validation
- quasi is intentionally only partially formalized; unknown lines inside a quasi block may be preserved as narrative rather than rejected outright