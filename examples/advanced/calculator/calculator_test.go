package calculator

import (
	"reflect"
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	testCases := []struct {
		name      string
		source    string
		want      int
		wantError string
	}{
		{name: "precedence", source: "1 + 2 * 3", want: 7},
		{name: "parentheses", source: "(1 + 2) * 3", want: 9},
		{name: "division then addition", source: "8 / 2 + 1", want: 5},
		{name: "nested parentheses", source: "8 / (2 + 2)", want: 2},
		{name: "division by zero", source: "1 / 0", wantError: "division by zero"},
		{name: "parse error", source: "1 + )", wantError: "expected literal or '('"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := Calculate(testCase.source)
			if testCase.wantError != "" {
				if err == nil {
					t.Fatalf("Calculate(%q) error = nil, want %q", testCase.source, testCase.wantError)
				}
				if !strings.Contains(err.Error(), testCase.wantError) {
					t.Fatalf("Calculate(%q) error = %q, want substring %q", testCase.source, err.Error(), testCase.wantError)
				}
				return
			}

			if err != nil {
				t.Fatalf("Calculate(%q) error = %v", testCase.source, err)
			}
			if got != testCase.want {
				t.Fatalf("Calculate(%q) = %d, want %d", testCase.source, got, testCase.want)
			}
		})
	}
}

func TestParsePrecedence(t *testing.T) {
	got, err := Parse("1 + 2 * 3")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	want := BinaryExpr{
		Op: "+",
		Left: IntLit{Value: 1},
		Right: BinaryExpr{
			Op:    "*",
			Left:  IntLit{Value: 2},
			Right: IntLit{Value: 3},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Parse AST = %#v, want %#v", got, want)
	}
}

func TestParseRejectsInvalidCharacter(t *testing.T) {
	_, err := Parse("1 + a")
	if err == nil {
		t.Fatal("Parse returned nil error for invalid character")
	}
	if !strings.Contains(err.Error(), "invalid character") {
		t.Fatalf("Parse error = %q, want invalid character message", err.Error())
	}
}