package validator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yy/len/internal/diag"
	"github.com/yy/len/internal/loader"
	"github.com/yy/len/internal/validator"
)

func TestValidatorReportsArityAndUnknownNames(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "sample.l1")
	mustWrite(t, path, `type Seq
rel Sorted(s: Seq)
fn bubble_sort(input: Seq) -> output: Seq
    ensures Sorted(output, input)
spec check
    given input: Seq
    must Missing(input)
`)

	program := loader.Loader{Root: root}.LoadPaths([]string{path})
	diags := validator.Validator{ProfileDir: filepath.Join(root, "profiles")}.Validate(program)
	assertCodes(t, diags, "validator.arity.rel", "validator.name.unresolved")
}

func TestValidatorReportsDuplicateBindersAndQuasiIssues(t *testing.T) {
	root := t.TempDir()
	writeProceduralProfile(t, root)
	path := filepath.Join(root, "sample.l1")
	mustWrite(t, path, `type Seq
fn bubble_sort(input: Seq) -> output: Seq
    ensures output = input
    quasi using style ProceduralAlgorithm:
        let xs := input
      else:
        return xs
spec bad
    given input, input: Seq
    must output = input
`)

	program := loader.Loader{Root: root}.LoadPaths([]string{path})
	diags := validator.Validator{ProfileDir: filepath.Join(root, "doc", "proposals", "accepted", "lip-0001-cli-parser-n-validation")}.Validate(program)
	assertCodes(t, diags, "quasi.continuation.orphan", "validator.binder.duplicate", "validator.name.unresolved")
}

func assertCodes(t *testing.T, diags []diag.Diagnostic, want ...string) {
	t.Helper()
	seen := map[string]bool{}
	for _, item := range diags {
		seen[item.Code] = true
	}
	for _, code := range want {
		if !seen[code] {
			t.Fatalf("expected diagnostic %q in %#v", code, diags)
		}
	}
}

func mustWrite(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func writeProceduralProfile(t *testing.T, root string) {
	t.Helper()
	path := filepath.Join(root, "doc", "proposals", "accepted", "lip-0001-cli-parser-n-validation", "procedural-algorithm.quasi-style.yaml")
	mustWrite(t, path, `version: 1
style:
	name: ProceduralAlgorithm
layout:
	indentation:
		mode: spaces
		width: 4
	allowBlankLines: true
keywords:
	blockOpeners: [if, while, for]
	blockContinuations: [else, else if]
	simpleStatements: [let, set, append, return]
slots:
	identifier: {pattern: '[A-Za-z_][A-Za-z0-9_]*'}
	expr: {pattern: '.+'}
	formula: {pattern: '.+'}
	target: {pattern: '.+'}
rules:
	- id: let-assign
		kind: statement
		keyword: let
		pattern: '^let\s+(?P<name>[A-Za-z_][A-Za-z0-9_]*)\s*:=\s*(?P<expr>.+)$'
	- id: else-block
		kind: continuation
		keyword: else
		pattern: '^else:$'
		opensBlock: true
		attachesTo: [if-block]
		mustAlignWithParent: true
	- id: return-expr
		kind: statement
		keyword: return
		pattern: '^return\s+.+$'
validation:
	firstTokenMustBeKeyword: true
	unknownKeywordPolicy: reject
	nonMatchingLinePolicy: reject
	requireContinuationImmediatelyAfterParentBlock: true
	requireConsistentIndentation: true
	requireIndentAfterBlockOpener: true
`)
}
