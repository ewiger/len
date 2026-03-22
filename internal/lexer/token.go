package lexer
package lexer

import "github.com/yy/len/internal/diag"

// Kind identifies one lexical token category.
type Kind string

const (
	KindEOF       Kind = "eof"
	KindIdentifier Kind = "identifier"
	KindKeyword   Kind = "keyword"
	KindNumber    Kind = "number"
	KindString    Kind = "string"
	KindComment   Kind = "comment"
	KindDocstring Kind = "docstring"
	KindOperator  Kind = "operator"
	KindDot       Kind = "."
	KindColon     Kind = ":"
	KindComma     Kind = ","
	KindLParen    Kind = "("
	KindRParen    Kind = ")"
)

// Token is one lexeme plus its source span.
type Token struct {
	Kind   Kind
	Lexeme string
	Span   diag.Span
}

var keywords = map[string]struct{}{
	"and": {}, "as": {}, "const": {}, "else": {}, "ensures": {},
	"exists": {}, "false": {}, "fn": {}, "forall": {}, "from": {},
	"given": {}, "iff": {}, "impl": {}, "implements": {}, "implies": {},
	"import": {}, "in": {}, "keyword": {}, "must": {}, "not": {},
	"or": {}, "quasi": {}, "rel": {}, "requires": {}, "spec": {},
	"struct": {}, "symbol": {}, "syntax": {}, "trait": {}, "true": {},
	"type": {}, "where": {}, "while": {}, "for": {},
	"subsetof": {},
}

// IsKeyword reports whether the identifier is reserved.
func IsKeyword(value string) bool {
	_, ok := keywords[value]
	return ok
}