package parser_test

import (
	"testing"

	"github.com/yy/len/internal/ast"
	"github.com/yy/len/internal/diag"
	"github.com/yy/len/internal/parser"
)

func TestParseFnSpecAndQuasi(t *testing.T) {
	source := `type Seq
rel Sorted(s: Seq)
rel At(s: Seq, i: Seq)

fn bubble_sort(input: Seq) -> output: Seq
    ensures Sorted(output)
    quasi using style ProceduralAlgorithm:
        let xs := input
        if true:
            return xs

spec bubble_sort_correct
    given input, output: Seq
    must Sorted(output) implies
        At(output, input)
`

	file, diags := parser.Parse("sample.l1", source)
	if len(diags) != 0 {
		t.Fatalf("Parse returned diagnostics: %#v", diags)
	}
	if len(file.Decls) != 5 {
		t.Fatalf("got %d declarations, want 5", len(file.Decls))
	}
	fn, ok := file.Decls[3].(*ast.FnDecl)
	if !ok {
		t.Fatalf("fourth decl = %T, want *ast.FnDecl", file.Decls[3])
	}
	if fn.Quasi == nil || fn.Quasi.StyleName != "ProceduralAlgorithm" {
		t.Fatalf("fn quasi = %#v, want style ProceduralAlgorithm", fn.Quasi)
	}
	if len(fn.Quasi.Block.Lines) != 3 {
		t.Fatalf("quasi lines = %d, want 3", len(fn.Quasi.Block.Lines))
	}
	spec, ok := file.Decls[4].(*ast.SpecDecl)
	if !ok {
		t.Fatalf("fifth decl = %T, want *ast.SpecDecl", file.Decls[4])
	}
	if len(spec.Given) != 2 {
		t.Fatalf("spec given binders = %d, want 2", len(spec.Given))
	}
	if _, ok := spec.Must.(*ast.BinaryExpr); !ok {
		t.Fatalf("spec must = %T, want *ast.BinaryExpr", spec.Must)
	}
}

func TestParseSyntaxAndQuantifier(t *testing.T) {
	source := `type Int
type Seq
syntax x in s where x: Int, s: Seq implies Member(x, s)
spec ordered
    given s: Seq
    must forall i: Int, j: Int . i < j implies Member(i, s)
`

	file, diags := parser.Parse("sample.l1", source)
	if len(diags) != 0 {
		t.Fatalf("Parse returned diagnostics: %#v", diags)
	}
	if _, ok := file.Decls[2].(*ast.SyntaxDecl); !ok {
		t.Fatalf("third decl = %T, want *ast.SyntaxDecl", file.Decls[2])
	}
	spec := file.Decls[3].(*ast.SpecDecl)
	if _, ok := spec.Must.(*ast.QuantifiedExpr); !ok {
		t.Fatalf("must = %T, want quantified expression", spec.Must)
	}
}

func TestParseFnRequiresQuasi(t *testing.T) {
	source := `type Seq
fn bubble_sort(input: Seq) -> output: Seq
    ensures output = input
`

	_, diags := parser.Parse("sample.l1", source)
	if !hasDiagCode(diags, "parser.fn.quasi.required") {
		t.Fatalf("expected parser.fn.quasi.required, got %#v", diags)
	}
}

func TestParseStructAndContract(t *testing.T) {
	source := `type Int
struct Pair
    left: Int
    right: Int

contract Arithmetic(T: Int)
    rel Add(left: T, right: T, value: T)
    spec identity
        given x: T
        must Add(x, x, x)
`

	file, diags := parser.Parse("sample.l1", source)
	if len(diags) != 0 {
		t.Fatalf("Parse returned diagnostics: %#v", diags)
	}
	if len(file.Decls) != 3 {
		t.Fatalf("got %d declarations, want 3", len(file.Decls))
	}
	strct, ok := file.Decls[1].(*ast.StructDecl)
	if !ok {
		t.Fatalf("second decl = %T, want *ast.StructDecl", file.Decls[1])
	}
	if len(strct.Fields) != 2 {
		t.Fatalf("struct fields = %d, want 2", len(strct.Fields))
	}
	contract, ok := file.Decls[2].(*ast.ContractDecl)
	if !ok {
		t.Fatalf("third decl = %T, want *ast.ContractDecl", file.Decls[2])
	}
	if len(contract.Params) != 1 {
		t.Fatalf("contract params = %d, want 1", len(contract.Params))
	}
	if len(contract.Members) != 2 {
		t.Fatalf("contract members = %d, want 2", len(contract.Members))
	}
	if _, ok := contract.Members[0].(*ast.RelDecl); !ok {
		t.Fatalf("first contract member = %T, want *ast.RelDecl", contract.Members[0])
	}
	if _, ok := contract.Members[1].(*ast.SpecDecl); !ok {
		t.Fatalf("second contract member = %T, want *ast.SpecDecl", contract.Members[1])
	}
}

func hasDiagCode(diags []diag.Diagnostic, code string) bool {
	for _, item := range diags {
		if item.Code == code {
			return true
		}
	}
	return false
}
