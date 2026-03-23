package calculator

import (
	"fmt"
	"strconv"
	"unicode"
)

type Expr interface {
	isExpr()
}

type IntLit struct {
	Value int
}

func (IntLit) isExpr() {}

type BinaryExpr struct {
	Op    string
	Left  Expr
	Right Expr
}

func (BinaryExpr) isExpr() {}

type tokenKind int

const (
	tokenInt tokenKind = iota
	tokenPlus
	tokenMinus
	tokenStar
	tokenSlash
	tokenLParen
	tokenRParen
	tokenEOF
)

type token struct {
	kind tokenKind
	text string
	pos  int
}

func lex(source string) ([]token, error) {
	var out []token
	runes := []rune(source)
	i := 0

	for i < len(runes) {
		r := runes[i]

		if unicode.IsSpace(r) {
			i++
			continue
		}

		if unicode.IsDigit(r) {
			start := i
			for i < len(runes) && unicode.IsDigit(runes[i]) {
				i++
			}
			out = append(out, token{kind: tokenInt, text: string(runes[start:i]), pos: start})
			continue
		}

		switch r {
		case '+':
			out = append(out, token{kind: tokenPlus, text: "+", pos: i})
		case '-':
			out = append(out, token{kind: tokenMinus, text: "-", pos: i})
		case '*':
			out = append(out, token{kind: tokenStar, text: "*", pos: i})
		case '/':
			out = append(out, token{kind: tokenSlash, text: "/", pos: i})
		case '(':
			out = append(out, token{kind: tokenLParen, text: "(", pos: i})
		case ')':
			out = append(out, token{kind: tokenRParen, text: ")", pos: i})
		default:
			return nil, fmt.Errorf("invalid character %q at position %d", r, i)
		}
		i++
	}

	out = append(out, token{kind: tokenEOF, pos: len(runes)})
	return out, nil
}

type parser struct {
	tokens []token
	pos    int
}

func Parse(source string) (Expr, error) {
	tokens, err := lex(source)
	if err != nil {
		return nil, err
	}

	p := &parser{tokens: tokens}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if p.peek().kind != tokenEOF {
		tok := p.peek()
		return nil, fmt.Errorf("unexpected trailing input at position %d", tok.pos)
	}

	return expr, nil
}

func (p *parser) peek() token {
	return p.tokens[p.pos]
}

func (p *parser) consume() token {
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}

func (p *parser) parseExpr() (Expr, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for {
		switch p.peek().kind {
		case tokenPlus, tokenMinus:
			op := p.consume().text
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			left = BinaryExpr{Op: op, Left: left, Right: right}
		default:
			return left, nil
		}
	}
}

func (p *parser) parseTerm() (Expr, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for {
		switch p.peek().kind {
		case tokenStar, tokenSlash:
			op := p.consume().text
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			left = BinaryExpr{Op: op, Left: left, Right: right}
		default:
			return left, nil
		}
	}
}

func (p *parser) parseFactor() (Expr, error) {
	tok := p.peek()

	switch tok.kind {
	case tokenInt:
		p.consume()
		value, err := strconv.Atoi(tok.text)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", tok.text, tok.pos)
		}
		return IntLit{Value: value}, nil
	case tokenLParen:
		p.consume()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.peek().kind != tokenRParen {
			return nil, fmt.Errorf("expected ')' at position %d", p.peek().pos)
		}
		p.consume()
		return expr, nil
	default:
		return nil, fmt.Errorf("expected literal or '(' at position %d", tok.pos)
	}
}

func Eval(expr Expr) (int, error) {
	switch value := expr.(type) {
	case IntLit:
		return value.Value, nil
	case BinaryExpr:
		left, err := Eval(value.Left)
		if err != nil {
			return 0, err
		}

		right, err := Eval(value.Right)
		if err != nil {
			return 0, err
		}

		switch value.Op {
		case "+":
			return left + right, nil
		case "-":
			return left - right, nil
		case "*":
			return left * right, nil
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left / right, nil
		default:
			return 0, fmt.Errorf("unknown operator %q", value.Op)
		}
	default:
		return 0, fmt.Errorf("unknown expression node %T", expr)
	}
}

func Calculate(source string) (int, error) {
	expr, err := Parse(source)
	if err != nil {
		return 0, err
	}

	return Eval(expr)
}
