package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/yy/len/internal/diag"
)

// Lex tokenizes the provided source and retains comments and docstrings.
func Lex(filePath string, src string) ([]Token, error) {
	l := &state{
		filePath: filePath,
		src:      src,
		line:     1,
		column:   1,
		bol:      true,
	}

	tokens := make([]Token, 0, len(src)/4)
	for {
		token, err := l.nextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
		if token.Kind == KindEOF {
			return tokens, nil
		}
	}
}

// LexExpression tokenizes an inline expression and drops trivia tokens.
func LexExpression(filePath string, src string, start diag.Position) ([]Token, error) {
	tokens, err := Lex(filePath, src)
	if err != nil {
		return nil, err
	}
	filtered := make([]Token, 0, len(tokens))
	for _, token := range tokens {
		if token.Kind == KindComment || token.Kind == KindDocstring {
			continue
		}
		token.Span.Start.Line += start.Line - 1
		token.Span.End.Line += start.Line - 1
		if token.Span.Start.Line == start.Line {
			token.Span.Start.Column += start.Column - 1
		}
		if token.Span.End.Line == start.Line {
			token.Span.End.Column += start.Column - 1
		}
		filtered = append(filtered, token)
	}
	return filtered, nil
}

type state struct {
	filePath string
	src      string
	offset   int
	line     int
	column   int
	bol      bool
}

func (s *state) nextToken() (Token, error) {
	for {
		if s.offset >= len(s.src) {
			span := diag.Span{Start: s.position(), End: s.position()}
			return Token{Kind: KindEOF, Span: span}, nil
		}

		if strings.HasPrefix(s.src[s.offset:], "<#") {
			return s.scanBlockComment()
		}

		r, size := utf8.DecodeRuneInString(s.src[s.offset:])
		if r == utf8.RuneError && size == 1 {
			return Token{}, fmt.Errorf("%s:%d:%d: invalid utf-8", s.filePath, s.line, s.column)
		}

		if r == '\n' {
			s.advanceRune(r, size)
			continue
		}
		if unicode.IsSpace(r) {
			s.advanceRune(r, size)
			continue
		}
		if r == '#' && s.bol {
			return s.scanLineComment()
		}
		if strings.HasPrefix(s.src[s.offset:], "\"\"\"") {
			return s.scanTripleDocstring()
		}
		if r == '"' {
			if s.bol {
				return s.scanString(KindDocstring)
			}
			return s.scanString(KindString)
		}
		if isIdentifierStart(r) {
			return s.scanIdentifier()
		}
		if unicode.IsDigit(r) {
			return s.scanNumber()
		}
		return s.scanSymbol()
	}
}

func (s *state) scanIdentifier() (Token, error) {
	start := s.position()
	begin := s.offset
	for s.offset < len(s.src) {
		r, size := utf8.DecodeRuneInString(s.src[s.offset:])
		if !isIdentifierContinue(r) {
			break
		}
		s.advanceRune(r, size)
	}
	lexeme := s.src[begin:s.offset]
	kind := KindIdentifier
	if IsKeyword(lexeme) {
		kind = KindKeyword
	}
	return Token{Kind: kind, Lexeme: lexeme, Span: diag.Span{Start: start, End: s.position()}}, nil
}

func (s *state) scanNumber() (Token, error) {
	start := s.position()
	begin := s.offset
	for s.offset < len(s.src) {
		r, size := utf8.DecodeRuneInString(s.src[s.offset:])
		if !unicode.IsDigit(r) {
			break
		}
		s.advanceRune(r, size)
	}
	return Token{Kind: KindNumber, Lexeme: s.src[begin:s.offset], Span: diag.Span{Start: start, End: s.position()}}, nil
}

func (s *state) scanString(kind Kind) (Token, error) {
	start := s.position()
	begin := s.offset
	s.advanceString("\"")
	for s.offset < len(s.src) {
		r, size := utf8.DecodeRuneInString(s.src[s.offset:])
		if r == '\\' {
			s.advanceRune(r, size)
			if s.offset < len(s.src) {
				r, size = utf8.DecodeRuneInString(s.src[s.offset:])
				s.advanceRune(r, size)
			}
			continue
		}
		if r == '"' {
			s.advanceRune(r, size)
			return Token{Kind: kind, Lexeme: s.src[begin:s.offset], Span: diag.Span{Start: start, End: s.position()}}, nil
		}
		s.advanceRune(r, size)
	}
	return Token{}, fmt.Errorf("%s:%d:%d: unterminated string", s.filePath, start.Line, start.Column)
}

func (s *state) scanTripleDocstring() (Token, error) {
	start := s.position()
	begin := s.offset
	s.advanceString("\"\"\"")
	for s.offset < len(s.src) {
		if strings.HasPrefix(s.src[s.offset:], "\"\"\"") {
			s.advanceString("\"\"\"")
			return Token{Kind: KindDocstring, Lexeme: s.src[begin:s.offset], Span: diag.Span{Start: start, End: s.position()}}, nil
		}
		r, size := utf8.DecodeRuneInString(s.src[s.offset:])
		s.advanceRune(r, size)
	}
	return Token{}, fmt.Errorf("%s:%d:%d: unterminated triple docstring", s.filePath, start.Line, start.Column)
}

func (s *state) scanLineComment() (Token, error) {
	start := s.position()
	begin := s.offset
	for s.offset < len(s.src) {
		r, size := utf8.DecodeRuneInString(s.src[s.offset:])
		if r == '\n' {
			break
		}
		s.advanceRune(r, size)
	}
	return Token{Kind: KindComment, Lexeme: s.src[begin:s.offset], Span: diag.Span{Start: start, End: s.position()}}, nil
}

func (s *state) scanBlockComment() (Token, error) {
	start := s.position()
	begin := s.offset
	s.advanceString("<#")
	for s.offset < len(s.src) {
		if strings.HasPrefix(s.src[s.offset:], "#>") {
			s.advanceString("#>")
			return Token{Kind: KindComment, Lexeme: s.src[begin:s.offset], Span: diag.Span{Start: start, End: s.position()}}, nil
		}
		r, size := utf8.DecodeRuneInString(s.src[s.offset:])
		s.advanceRune(r, size)
	}
	return Token{}, fmt.Errorf("%s:%d:%d: unterminated block comment", s.filePath, start.Line, start.Column)
}

func (s *state) scanSymbol() (Token, error) {
	start := s.position()
	for _, candidate := range []string{"<=>", "->", ":=", "=>", "/\\", "\\/", "<=", ">=", "!=", ".."} {
		if strings.HasPrefix(s.src[s.offset:], candidate) {
			s.advanceString(candidate)
			return Token{Kind: KindOperator, Lexeme: candidate, Span: diag.Span{Start: start, End: s.position()}}, nil
		}
	}

	r, size := utf8.DecodeRuneInString(s.src[s.offset:])
	s.advanceRune(r, size)
	lexeme := string(r)
	switch r {
	case '.':
		return Token{Kind: KindDot, Lexeme: lexeme, Span: diag.Span{Start: start, End: s.position()}}, nil
	case ':':
		return Token{Kind: KindColon, Lexeme: lexeme, Span: diag.Span{Start: start, End: s.position()}}, nil
	case ',':
		return Token{Kind: KindComma, Lexeme: lexeme, Span: diag.Span{Start: start, End: s.position()}}, nil
	case '(':
		return Token{Kind: KindLParen, Lexeme: lexeme, Span: diag.Span{Start: start, End: s.position()}}, nil
	case ')':
		return Token{Kind: KindRParen, Lexeme: lexeme, Span: diag.Span{Start: start, End: s.position()}}, nil
	case '=', '<', '>', '+', '-', '*', '/', '!', '|':
		return Token{Kind: KindOperator, Lexeme: lexeme, Span: diag.Span{Start: start, End: s.position()}}, nil
	default:
		return Token{}, fmt.Errorf("%s:%d:%d: unexpected character %q", s.filePath, start.Line, start.Column, r)
	}
}

func (s *state) advanceString(value string) {
	for _, r := range value {
		_, size := utf8.DecodeRuneInString(s.src[s.offset:])
		s.advanceRune(r, size)
	}
}

func (s *state) advanceRune(r rune, size int) {
	s.offset += size
	if r == '\n' {
		s.line++
		s.column = 1
		s.bol = true
		return
	}
	s.column++
	if !unicode.IsSpace(r) {
		s.bol = false
	}
}

func (s *state) position() diag.Position {
	return diag.Position{File: s.filePath, Line: s.line, Column: s.column}
}

func isIdentifierStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdentifierContinue(r rune) bool {
	return isIdentifierStart(r) || unicode.IsDigit(r)
}
