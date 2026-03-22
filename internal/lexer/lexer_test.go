package lexer_test

import (
	"testing"

	"github.com/yy/len/internal/lexer"
)

func TestLexRecognizesTriviaKeywordsAndOperators(t *testing.T) {
	source := "# comment\n\"doc\"\nfn hello(input: Seq) -> output: Seq\n    ensures output = input /\\ true <=> false\n"
	tokens, err := lexer.Lex("sample.l1", source)
	if err != nil {
		t.Fatalf("Lex returned error: %v", err)
	}

	var kinds []lexer.Kind
	var lexemes []string
	for _, token := range tokens {
		kinds = append(kinds, token.Kind)
		lexemes = append(lexemes, token.Lexeme)
	}

	wantKinds := []lexer.Kind{
		lexer.KindComment,
		lexer.KindDocstring,
		lexer.KindKeyword,
		lexer.KindIdentifier,
		lexer.KindLParen,
		lexer.KindIdentifier,
		lexer.KindColon,
		lexer.KindIdentifier,
		lexer.KindRParen,
		lexer.KindOperator,
		lexer.KindIdentifier,
		lexer.KindColon,
		lexer.KindIdentifier,
		lexer.KindKeyword,
		lexer.KindIdentifier,
		lexer.KindOperator,
		lexer.KindIdentifier,
		lexer.KindOperator,
		lexer.KindKeyword,
		lexer.KindOperator,
		lexer.KindKeyword,
		lexer.KindEOF,
	}
	if len(kinds) != len(wantKinds) {
		t.Fatalf("got %d tokens, want %d\nlexemes=%v", len(kinds), len(wantKinds), lexemes)
	}
	for i := range wantKinds {
		if kinds[i] != wantKinds[i] {
			t.Fatalf("token %d: got %q (%q), want %q", i, kinds[i], lexemes[i], wantKinds[i])
		}
	}
	if lexemes[9] != "->" || lexemes[17] != "/\\" || lexemes[19] != "<=>" {
		t.Fatalf("operator lexemes not preserved: %v", lexemes)
	}
}

func TestLexSupportsBlockCommentsAndTripleDocstrings(t *testing.T) {
	source := "<# hidden\nblock #>\n\"\"\"multi\nline\"\"\"\nimport core.math.set\n"
	tokens, err := lexer.Lex("sample.l1", source)
	if err != nil {
		t.Fatalf("Lex returned error: %v", err)
	}

	if tokens[0].Kind != lexer.KindComment {
		t.Fatalf("first token = %q, want comment", tokens[0].Kind)
	}
	if tokens[1].Kind != lexer.KindDocstring {
		t.Fatalf("second token = %q, want docstring", tokens[1].Kind)
	}
	if tokens[2].Kind != lexer.KindKeyword || tokens[2].Lexeme != "import" {
		t.Fatalf("third token = %#v, want import keyword", tokens[2])
	}
}
