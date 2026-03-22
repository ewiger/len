package ast

import "github.com/yy/len/internal/diag"

// This file sketches the core AST for the accepted LIP-0001 parser direction.
// Quasi blocks are stored as raw lines plus header metadata; style-internal
// statements are not parsed into a dedicated AST at host-grammar time.

// File is the root AST node for one parsed source file.
type File struct {
	Decls []Decl
	Span  diag.Span
}

// Decl is implemented by all top-level declarations.
type Decl interface {
	declNode()
	GetSpan() diag.Span
}

// ImportDecl models `import a.b.c`.
type ImportDecl struct {
	ModulePath []string
	Span       diag.Span
}

func (d *ImportDecl) declNode()          {}
func (d *ImportDecl) GetSpan() diag.Span { return d.Span }

// TypeDecl models `type Name`.
type TypeDecl struct {
	Name string
	Span diag.Span
}

func (d *TypeDecl) declNode()          {}
func (d *TypeDecl) GetSpan() diag.Span { return d.Span }

// RelDecl models `rel Name(params...)`.
type RelDecl struct {
	Name   string
	Params []Binder
	Span   diag.Span
}

func (d *RelDecl) declNode()          {}
func (d *RelDecl) GetSpan() diag.Span { return d.Span }

// ConstDecl models `const Name : TypeExpr`.
type ConstDecl struct {
	Name string
	Type Expr
	Span diag.Span
}

func (d *ConstDecl) declNode()          {}
func (d *ConstDecl) GetSpan() diag.Span { return d.Span }

// KeywordDecl models `keyword name`.
type KeywordDecl struct {
	Name string
	Span diag.Span
}

func (d *KeywordDecl) declNode()          {}
func (d *KeywordDecl) GetSpan() diag.Span { return d.Span }

// SymbolDecl models `symbol Name as "..."`.
type SymbolDecl struct {
	Name  string
	Value string
	Span  diag.Span
}

func (d *SymbolDecl) declNode()          {}
func (d *SymbolDecl) GetSpan() diag.Span { return d.Span }

// TraitDecl models `trait Name`.
type TraitDecl struct {
	Name string
	Span diag.Span
}

func (d *TraitDecl) declNode()          {}
func (d *TraitDecl) GetSpan() diag.Span { return d.Span }

// ImplDecl models `impl TypeExpr : TraitName`.
type ImplDecl struct {
	Type      Expr
	TraitName string
	Span      diag.Span
}

func (d *ImplDecl) declNode()          {}
func (d *ImplDecl) GetSpan() diag.Span { return d.Span }

// SyntaxDecl models `syntax surface where binders implies canonical`.
type SyntaxDecl struct {
	Surface   Expr
	Binders   []Binder
	Canonical Expr
	Span      diag.Span
}

func (d *SyntaxDecl) declNode()          {}
func (d *SyntaxDecl) GetSpan() diag.Span { return d.Span }

// SpecDecl models a declarative `spec` block.
type SpecDecl struct {
	Name   string
	Given  []Binder
	Must   Expr
	Trivia []Trivia
	Span   diag.Span
}

func (d *SpecDecl) declNode()          {}
func (d *SpecDecl) GetSpan() diag.Span { return d.Span }

// FnDecl models a function-like declaration with a required quasi body.
type FnDecl struct {
	Name    string
	Params  []Binder
	Result  *Binder
	Clauses []FnClause
	Quasi   *QuasiClause
	Trivia  []Trivia
	Span    diag.Span
}

func (d *FnDecl) declNode()          {}
func (d *FnDecl) GetSpan() diag.Span { return d.Span }

// Binder names a local variable and optional type slot.
type Binder struct {
	Name string
	Type Expr
	Span diag.Span
}

// FnClause is implemented by requires, ensures, and implements clauses.
type FnClause interface {
	fnClauseNode()
	GetSpan() diag.Span
}

// RequiresClause models `requires Formula`.
type RequiresClause struct {
	Formula Expr
	Span    diag.Span
}

func (c *RequiresClause) fnClauseNode()      {}
func (c *RequiresClause) GetSpan() diag.Span { return c.Span }

// EnsuresClause models `ensures Formula`.
type EnsuresClause struct {
	Formula Expr
	Span    diag.Span
}

func (c *EnsuresClause) fnClauseNode()      {}
func (c *EnsuresClause) GetSpan() diag.Span { return c.Span }

// ImplementsClause models `implements Formula`.
type ImplementsClause struct {
	Formula Expr
	Span    diag.Span
}

func (c *ImplementsClause) fnClauseNode()      {}
func (c *ImplementsClause) GetSpan() diag.Span { return c.Span }

// QuasiClause stores only header information plus raw indented lines.
type QuasiClause struct {
	StyleName  string
	HeaderSpan diag.Span
	Block      RawQuasiBlock
	Span       diag.Span
}

// HasExplicitStyle reports whether the quasi header declared a style name.
func (q *QuasiClause) HasExplicitStyle() bool {
	return q != nil && q.StyleName != ""
}

// RawQuasiBlock preserves the quasi body exactly as captured by the host parser.
type RawQuasiBlock struct {
	Lines []RawQuasiLine
	Span  diag.Span
}

// RawQuasiLine preserves the original text and indentation of one quasi line.
type RawQuasiLine struct {
	Text         string
	TrimmedText  string
	IndentColumn int
	LineSpan     diag.Span
}

// Trivia models ignored but retained source material for tooling.
type Trivia interface {
	triviaNode()
	GetSpan() diag.Span
}

// Comment stores a comment token if tooling chooses to retain it.
type Comment struct {
	Text string
	Span diag.Span
}

func (t *Comment) triviaNode()        {}
func (t *Comment) GetSpan() diag.Span { return t.Span }

// Docstring stores a docstring token if tooling chooses to retain it.
type Docstring struct {
	Text string
	Span diag.Span
}

func (t *Docstring) triviaNode()        {}
func (t *Docstring) GetSpan() diag.Span { return t.Span }

// Expr is a deliberately small placeholder for the host expression tree.
// The parser can refine this hierarchy as expression parsing is implemented.
type Expr interface {
	exprNode()
	GetSpan() diag.Span
}

// IdentExpr models a single identifier.
type IdentExpr struct {
	Name string
	Span diag.Span
}

func (e *IdentExpr) exprNode()          {}
func (e *IdentExpr) GetSpan() diag.Span { return e.Span }

// QualifiedExpr models dotted names such as a.b.c.
type QualifiedExpr struct {
	Parts []string
	Span  diag.Span
}

func (e *QualifiedExpr) exprNode()          {}
func (e *QualifiedExpr) GetSpan() diag.Span { return e.Span }

// StringExpr models string literals and docstring-like standalone values.
type StringExpr struct {
	Value string
	Span  diag.Span
}

func (e *StringExpr) exprNode()          {}
func (e *StringExpr) GetSpan() diag.Span { return e.Span }

// NumberExpr models integer literals.
type NumberExpr struct {
	Value string
	Span  diag.Span
}

func (e *NumberExpr) exprNode()          {}
func (e *NumberExpr) GetSpan() diag.Span { return e.Span }

// BoolExpr models true and false.
type BoolExpr struct {
	Value bool
	Span  diag.Span
}

func (e *BoolExpr) exprNode()          {}
func (e *BoolExpr) GetSpan() diag.Span { return e.Span }

// ApplyExpr models callee(arg1, arg2, ...).
type ApplyExpr struct {
	Callee Expr
	Args   []Expr
	Span   diag.Span
}

func (e *ApplyExpr) exprNode()          {}
func (e *ApplyExpr) GetSpan() diag.Span { return e.Span }

// UnaryExpr models unary host operators such as `not`.
type UnaryExpr struct {
	Op   string
	Expr Expr
	Span diag.Span
}

func (e *UnaryExpr) exprNode()          {}
func (e *UnaryExpr) GetSpan() diag.Span { return e.Span }

// BinaryExpr models infix host expressions and formulas.
type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
	Span  diag.Span
}

func (e *BinaryExpr) exprNode()          {}
func (e *BinaryExpr) GetSpan() diag.Span { return e.Span }

// QuantifiedExpr models forall/exists binders.
type QuantifiedExpr struct {
	Quantifier string
	Binders    []Binder
	Body       Expr
	Span       diag.Span
}

func (e *QuantifiedExpr) exprNode()          {}
func (e *QuantifiedExpr) GetSpan() diag.Span { return e.Span }

// GroupExpr models parenthesized subexpressions.
type GroupExpr struct {
	Inner Expr
	Span  diag.Span
}

func (e *GroupExpr) exprNode()          {}
func (e *GroupExpr) GetSpan() diag.Span { return e.Span }
