package ast

import (
	"testing"

	"github.com/yy/len/internal/diag"
)

var (
	_ Decl = (*ImportDecl)(nil)
	_ Decl = (*TypeDecl)(nil)
	_ Decl = (*StructDecl)(nil)
	_ Decl = (*RelDecl)(nil)
	_ Decl = (*ConstDecl)(nil)
	_ Decl = (*KeywordDecl)(nil)
	_ Decl = (*SymbolDecl)(nil)
	_ Decl = (*TraitDecl)(nil)
	_ Decl = (*ImplDecl)(nil)
	_ Decl = (*SyntaxDecl)(nil)
	_ Decl = (*SpecDecl)(nil)
	_ Decl = (*ContractDecl)(nil)
	_ Decl = (*FnDecl)(nil)

	_ FnClause = (*RequiresClause)(nil)
	_ FnClause = (*EnsuresClause)(nil)
	_ FnClause = (*ImplementsClause)(nil)

	_ Trivia = (*Comment)(nil)
	_ Trivia = (*Docstring)(nil)

	_ Expr = (*IdentExpr)(nil)
	_ Expr = (*QualifiedExpr)(nil)
	_ Expr = (*StringExpr)(nil)
	_ Expr = (*NumberExpr)(nil)
	_ Expr = (*BoolExpr)(nil)
	_ Expr = (*ApplyExpr)(nil)
	_ Expr = (*UnaryExpr)(nil)
	_ Expr = (*BinaryExpr)(nil)
	_ Expr = (*QuantifiedExpr)(nil)
	_ Expr = (*GroupExpr)(nil)
)

func TestQuasiClauseHasExplicitStyle(t *testing.T) {
	if (*QuasiClause)(nil).HasExplicitStyle() {
		t.Fatal("nil quasi clause reported an explicit style")
	}
	if (&QuasiClause{}).HasExplicitStyle() {
		t.Fatal("empty quasi clause reported an explicit style")
	}
	if !(&QuasiClause{StyleName: "ProceduralAlgorithm"}).HasExplicitStyle() {
		t.Fatal("styled quasi clause did not report an explicit style")
	}
}

func TestDeclGetSpanReturnsStoredSpan(t *testing.T) {
	span := testSpan("decl.l1", 3, 5, 7, 11)
	tests := []struct {
		name string
		decl Decl
	}{
		{name: "import", decl: &ImportDecl{Span: span}},
		{name: "type", decl: &TypeDecl{Span: span}},
		{name: "struct", decl: &StructDecl{Span: span}},
		{name: "rel", decl: &RelDecl{Span: span}},
		{name: "const", decl: &ConstDecl{Span: span}},
		{name: "keyword", decl: &KeywordDecl{Span: span}},
		{name: "symbol", decl: &SymbolDecl{Span: span}},
		{name: "trait", decl: &TraitDecl{Span: span}},
		{name: "impl", decl: &ImplDecl{Span: span}},
		{name: "syntax", decl: &SyntaxDecl{Span: span}},
		{name: "spec", decl: &SpecDecl{Span: span}},
		{name: "contract", decl: &ContractDecl{Span: span}},
		{name: "fn", decl: &FnDecl{Span: span}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.decl.GetSpan(); got != span {
				t.Fatalf("GetSpan() = %#v, want %#v", got, span)
			}
		})
	}
}

func TestFnClauseGetSpanReturnsStoredSpan(t *testing.T) {
	span := testSpan("fn.l1", 10, 1, 10, 20)
	tests := []struct {
		name   string
		clause FnClause
	}{
		{name: "requires", clause: &RequiresClause{Span: span}},
		{name: "ensures", clause: &EnsuresClause{Span: span}},
		{name: "implements", clause: &ImplementsClause{Span: span}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.clause.GetSpan(); got != span {
				t.Fatalf("GetSpan() = %#v, want %#v", got, span)
			}
		})
	}
}

func TestTriviaGetSpanReturnsStoredSpan(t *testing.T) {
	span := testSpan("comments.l1", 2, 1, 2, 12)
	tests := []struct {
		name   string
		trivia Trivia
	}{
		{name: "comment", trivia: &Comment{Span: span}},
		{name: "docstring", trivia: &Docstring{Span: span}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.trivia.GetSpan(); got != span {
				t.Fatalf("GetSpan() = %#v, want %#v", got, span)
			}
		})
	}
}

func TestExprGetSpanReturnsStoredSpan(t *testing.T) {
	span := testSpan("expr.l1", 5, 3, 5, 17)
	tests := []struct {
		name string
		expr Expr
	}{
		{name: "ident", expr: &IdentExpr{Span: span}},
		{name: "qualified", expr: &QualifiedExpr{Span: span}},
		{name: "string", expr: &StringExpr{Span: span}},
		{name: "number", expr: &NumberExpr{Span: span}},
		{name: "bool", expr: &BoolExpr{Span: span}},
		{name: "apply", expr: &ApplyExpr{Span: span}},
		{name: "unary", expr: &UnaryExpr{Span: span}},
		{name: "binary", expr: &BinaryExpr{Span: span}},
		{name: "quantified", expr: &QuantifiedExpr{Span: span}},
		{name: "group", expr: &GroupExpr{Span: span}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.expr.GetSpan(); got != span {
				t.Fatalf("GetSpan() = %#v, want %#v", got, span)
			}
		})
	}
}

func testSpan(file string, startLine int, startColumn int, endLine int, endColumn int) diag.Span {
	return diag.Span{
		Start: diag.Position{File: file, Line: startLine, Column: startColumn},
		End:   diag.Position{File: file, Line: endLine, Column: endColumn},
	}
}
