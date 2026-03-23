package quasi_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yy/len/internal/quasi"
)

func TestLoadProfileAndValidateSurface(t *testing.T) {
	root := t.TempDir()
	profilePath := filepath.Join(root, "procedural-algorithm.quasi-style.yaml")
	writeFile(t, profilePath, `version: 1
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
  - id: if-block
    kind: block
    keyword: if
    pattern: '^if\s+(?P<formula>.+):$'
    opensBlock: true
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

	profile, err := quasi.LoadProfile(profilePath)
	if err != nil {
		t.Fatalf("LoadProfile returned error: %v", err)
	}
	result := quasi.Validator{Profile: profile}.Validate(quasi.Block{Lines: []quasi.RawLine{
		{TrimmedText: "if ready:", IndentColumn: 4},
		{TrimmedText: "return xs", IndentColumn: 8},
	}})
	if !result.OK() {
		t.Fatalf("Validate returned diagnostics: %#v", result.Diagnostics)
	}
	bad := quasi.Validator{Profile: profile}.Validate(quasi.Block{Lines: []quasi.RawLine{{TrimmedText: "else:", IndentColumn: 4}}})
	if len(bad.Diagnostics) == 0 {
		t.Fatal("expected continuation diagnostic")
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func TestProceduralAlgorithmProfileAcceptsCalculatorForms(t *testing.T) {
  profile, err := quasi.LoadProfile(filepath.Join("..", "..", "doc", "proposals", "accepted", "lip-0001-cli-parser-n-validation", "procedural-algorithm.quasi-style.yaml"))
  if err != nil {
    t.Fatalf("LoadProfile returned error: %v", err)
  }

  result := quasi.Validator{Profile: profile}.Validate(quasi.Block{Lines: []quasi.RawLine{
    {TrimmedText: "let out := []", IndentColumn: 8},
    {TrimmedText: "append(out, token(\"int\", slice(source, start, i)))", IndentColumn: 8},
    {TrimmedText: "fn parse_factor():", IndentColumn: 8},
    {TrimmedText: "if tokens[pos].kind = \"int\":", IndentColumn: 12},
    {TrimmedText: "return IntLit(n)", IndentColumn: 16},
    {TrimmedText: "else:", IndentColumn: 12},
    {TrimmedText: "fail ParseError(\"expected literal or '('\")", IndentColumn: 16},
    {TrimmedText: "case expr of", IndentColumn: 8},
    {TrimmedText: "IntLit(n):", IndentColumn: 12},
    {TrimmedText: "return n", IndentColumn: 16},
    {TrimmedText: "BinaryExpr(\"+\", left, right):", IndentColumn: 12},
    {TrimmedText: "return eval(left) + eval(right)", IndentColumn: 16},
    {TrimmedText: "expect(tokens[pos], \")\")", IndentColumn: 8},
  }})
  if !result.OK() {
    t.Fatalf("Validate returned diagnostics: %#v", result.Diagnostics)
  }
}
