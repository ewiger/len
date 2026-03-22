package parser

import (
	"fmt"
	"strings"

	"github.com/yy/len/internal/ast"
	"github.com/yy/len/internal/diag"
	"github.com/yy/len/internal/lexer"
)

func parseExpression(filePath string, text string, start diag.Position) (ast.Expr, []diag.Diagnostic, bool) {
	tokens, err := lexer.LexExpression(filePath, text, start)
	if err != nil {
		return nil, []diag.Diagnostic{{Code: "parser.expr.lex", Message: err.Error(), Severity: diag.SeverityError, Span: diag.Span{Start: start, End: start}}}, false
	}
	ep := &exprParser{tokens: tokens, filePath: filePath}
	expr, ok := ep.parseExpr(1)
	if !ok {
		return nil, ep.diags, false
	}
	if ep.current().Kind != lexer.KindEOF {
		tok := ep.current()
		ep.error(tok.Span, "parser.expr.trailing", fmt.Sprintf("unexpected token %q", tok.Lexeme))
		return expr, ep.diags, false
	}
	return expr, ep.diags, len(ep.diags) == 0
}

type exprParser struct {
	tokens   []lexer.Token
	index    int
	filePath string
	diags    []diag.Diagnostic
}

func (p *exprParser) parseExpr(minPrec int) (ast.Expr, bool) {
	left, ok := p.parsePrefix()
	if !ok {
		return nil, false
	}
	for {
		token := p.current()
		op, ok := p.binaryOperator(token)
		if !ok {
			break
		}
		prec, rightAssoc := precedence(op)
		if prec < minPrec {
			break
		}
		p.index++
		nextMin := prec + 1
		if rightAssoc {
			nextMin = prec
		}
		right, ok := p.parseExpr(nextMin)
		if !ok {
			return nil, false
		}
		left = &ast.BinaryExpr{Left: left, Op: op, Right: right, Span: spanJoin(left.GetSpan(), right.GetSpan())}
	}
	return left, true
}

func (p *exprParser) parsePrefix() (ast.Expr, bool) {
	token := p.current()
	if token.Kind == lexer.KindKeyword && (token.Lexeme == "forall" || token.Lexeme == "exists") {
		return p.parseQuantified()
	}
	if token.Kind == lexer.KindKeyword && token.Lexeme == "not" {
		p.index++
		expr, ok := p.parseExpr(6)
		if !ok {
			return nil, false
		}
		return &ast.UnaryExpr{Op: token.Lexeme, Expr: expr, Span: spanJoin(token.Span, expr.GetSpan())}, true
	}
	if token.Kind == lexer.KindOperator && (token.Lexeme == "!" || token.Lexeme == "-") {
		p.index++
		expr, ok := p.parseExpr(6)
		if !ok {
			return nil, false
		}
		return &ast.UnaryExpr{Op: token.Lexeme, Expr: expr, Span: spanJoin(token.Span, expr.GetSpan())}, true
	}
	return p.parsePrimary()
}

func (p *exprParser) parseQuantified() (ast.Expr, bool) {
	quantifier := p.current()
	p.index++
	binders, ok := p.parseBindersUntilDot()
	if !ok {
		return nil, false
	}
	if p.current().Kind != lexer.KindDot {
		p.error(p.current().Span, "parser.quantifier.dot", "quantified expression requires `.` before the body")
		return nil, false
	}
	dot := p.current()
	p.index++
	body, ok := p.parseExpr(1)
	if !ok {
		return nil, false
	}
	return &ast.QuantifiedExpr{Quantifier: quantifier.Lexeme, Binders: binders, Body: body, Span: diag.Span{Start: quantifier.Span.Start, End: body.GetSpan().End}}, dot.Kind == lexer.KindDot
}

func (p *exprParser) parsePrimary() (ast.Expr, bool) {
	token := p.current()
	switch token.Kind {
	case lexer.KindIdentifier, lexer.KindKeyword:
		if token.Lexeme == "true" || token.Lexeme == "false" {
			p.index++
			return &ast.BoolExpr{Value: token.Lexeme == "true", Span: token.Span}, true
		}
		return p.parseNameOrCall()
	case lexer.KindNumber:
		p.index++
		return &ast.NumberExpr{Value: token.Lexeme, Span: token.Span}, true
	case lexer.KindString:
		p.index++
		return &ast.StringExpr{Value: token.Lexeme, Span: token.Span}, true
	case lexer.KindLParen:
		open := token
		p.index++
		inner, ok := p.parseExpr(1)
		if !ok {
			return nil, false
		}
		if p.current().Kind != lexer.KindRParen {
			p.error(p.current().Span, "parser.group.close", "expected closing `)`")
			return nil, false
		}
		close := p.current()
		p.index++
		return &ast.GroupExpr{Inner: inner, Span: diag.Span{Start: open.Span.Start, End: close.Span.End}}, true
	default:
		p.error(token.Span, "parser.expr.primary", fmt.Sprintf("unexpected token %q", token.Lexeme))
		return nil, false
	}
}

func (p *exprParser) parseNameOrCall() (ast.Expr, bool) {
	start := p.current()
	parts := []string{start.Lexeme}
	end := start.Span
	p.index++
	for p.current().Kind == lexer.KindDot {
		p.index++
		next := p.current()
		if next.Kind != lexer.KindIdentifier && next.Kind != lexer.KindKeyword {
			p.error(next.Span, "parser.qualified.name", "qualified name requires an identifier after `.`")
			return nil, false
		}
		parts = append(parts, next.Lexeme)
		end = next.Span
		p.index++
	}
	var expr ast.Expr
	if len(parts) == 1 {
		expr = &ast.IdentExpr{Name: parts[0], Span: diag.Span{Start: start.Span.Start, End: end.End}}
	} else {
		expr = &ast.QualifiedExpr{Parts: parts, Span: diag.Span{Start: start.Span.Start, End: end.End}}
	}

	for p.current().Kind == lexer.KindLParen {
		open := p.current()
		p.index++
		args := make([]ast.Expr, 0)
		if p.current().Kind != lexer.KindRParen {
			for {
				arg, ok := p.parseExpr(1)
				if !ok {
					return nil, false
				}
				args = append(args, arg)
				if p.current().Kind != lexer.KindComma {
					break
				}
				p.index++
			}
		}
		if p.current().Kind != lexer.KindRParen {
			p.error(p.current().Span, "parser.call.close", "function application requires closing `)`")
			return nil, false
		}
		close := p.current()
		p.index++
		expr = &ast.ApplyExpr{Callee: expr, Args: args, Span: diag.Span{Start: expr.GetSpan().Start, End: close.Span.End}}
		_ = open
	}
	return expr, true
}

func (p *exprParser) parseBindersUntilDot() ([]ast.Binder, bool) {
	segments := make([]lexer.Token, 0)
	depth := 0
	for {
		tok := p.current()
		if tok.Kind == lexer.KindEOF {
			break
		}
		if tok.Kind == lexer.KindLParen {
			depth++
		}
		if tok.Kind == lexer.KindRParen && depth > 0 {
			depth--
		}
		if tok.Kind == lexer.KindDot && depth == 0 {
			break
		}
		segments = append(segments, tok)
		p.index++
	}
	bp := &binderParser{owner: nil, tokens: segments, filePath: p.filePath}
	binders, ok := bp.parseTokens()
	p.diags = append(p.diags, bp.diags...)
	return binders, ok
}

func (p *exprParser) current() lexer.Token {
	if p.index >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.index]
}

func (p *exprParser) error(span diag.Span, code string, message string) {
	p.diags = append(p.diags, diag.Diagnostic{Code: code, Message: message, Severity: diag.SeverityError, Span: span})
}

func (p *exprParser) binaryOperator(token lexer.Token) (string, bool) {
	if token.Kind == lexer.KindKeyword {
		switch token.Lexeme {
		case "and", "or", "implies", "iff", "in", "subsetof":
			return token.Lexeme, true
		}
	}
	if token.Kind == lexer.KindOperator {
		switch token.Lexeme {
		case "=", "<", ">", "<=", ">=", "!=", "+", "-", "*", "/", "/\\", "\\/", "=>", "<=>":
			return token.Lexeme, true
		}
	}
	return "", false
}

func precedence(op string) (int, bool) {
	switch op {
	case "iff", "<=>":
		return 1, true
	case "implies", "=>":
		return 2, true
	case "or", "\\/":
		return 3, false
	case "and", "/\\":
		return 4, false
	case "=", "in", "subsetof", "<", ">", "<=", ">=", "!=":
		return 5, false
	case "+", "-":
		return 6, false
	case "*", "/":
		return 7, false
	default:
		return 0, false
	}
}

type binderParser struct {
	owner    *parser
	tokens   []lexer.Token
	index    int
	filePath string
	diags    []diag.Diagnostic
}

func newBinderParser(owner *parser, text string, start diag.Position) *binderParser {
	tokens, err := lexer.LexExpression(start.File, text, start)
	bp := &binderParser{owner: owner, tokens: tokens, filePath: start.File}
	if err != nil {
		bp.diags = append(bp.diags, diag.Diagnostic{Code: "parser.binders.lex", Message: err.Error(), Severity: diag.SeverityError, Span: diag.Span{Start: start, End: start}})
	}
	return bp
}

func (p *binderParser) parse() ([]ast.Binder, bool) {
	return p.parseTokens()
}

func (p *binderParser) parseTokens() ([]ast.Binder, bool) {
	if len(p.diags) > 0 {
		if p.owner != nil {
			p.owner.diags = append(p.owner.diags, p.diags...)
		}
		return nil, false
	}
	binders := make([]ast.Binder, 0)
	for p.current().Kind != lexer.KindEOF {
		names := make([]lexer.Token, 0)
		for {
			tok := p.current()
			if tok.Kind != lexer.KindIdentifier {
				p.error(tok.Span, "parser.binders.name", "binder list requires an identifier")
				return nil, p.finish(false)
			}
			names = append(names, tok)
			p.index++
			if p.current().Kind != lexer.KindComma {
				break
			}
			p.index++
			if p.current().Kind == lexer.KindIdentifier && p.peek().Kind == lexer.KindColon {
				tok := p.current()
				names = append(names, tok)
				p.index++
				break
			}
		}
		if p.current().Kind != lexer.KindColon {
			p.error(p.current().Span, "parser.binders.colon", "binder list requires `:` after binder names")
			return nil, p.finish(false)
		}
		p.index++
		typeTokens := make([]lexer.Token, 0)
		depth := 0
		for p.current().Kind != lexer.KindEOF {
			if p.current().Kind == lexer.KindLParen {
				depth++
			}
			if p.current().Kind == lexer.KindRParen && depth > 0 {
				depth--
			}
			if depth == 0 && p.current().Kind == lexer.KindComma && p.looksLikeBinderBoundary() {
				break
			}
			typeTokens = append(typeTokens, p.current())
			p.index++
		}
		typeExpr, ok := parseExprFromTokens(p.filePath, typeTokens)
		if !ok {
			return nil, p.finish(false)
		}
		for _, name := range names {
			binders = append(binders, ast.Binder{Name: name.Lexeme, Type: typeExpr, Span: diag.Span{Start: name.Span.Start, End: typeExpr.GetSpan().End}})
		}
		if p.current().Kind == lexer.KindComma {
			p.index++
		}
	}
	return binders, p.finish(true)
}

func parseExprFromTokens(filePath string, tokens []lexer.Token) (ast.Expr, bool) {
	if len(tokens) == 0 {
		return nil, false
	}
	cloned := append([]lexer.Token{}, tokens...)
	end := tokens[len(tokens)-1].Span.End
	cloned = append(cloned, lexer.Token{Kind: lexer.KindEOF, Span: diag.Span{Start: end, End: end}})
	ep := &exprParser{tokens: cloned, filePath: filePath}
	expr, ok := ep.parseExpr(1)
	return expr, ok && len(ep.diags) == 0
}

func (p *binderParser) current() lexer.Token {
	if p.index >= len(p.tokens) {
		if len(p.tokens) == 0 {
			return lexer.Token{Kind: lexer.KindEOF}
		}
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.index]
}

func (p *binderParser) peek() lexer.Token {
	if p.index+1 >= len(p.tokens) {
		return lexer.Token{Kind: lexer.KindEOF}
	}
	return p.tokens[p.index+1]
}

func (p *binderParser) looksLikeBinderBoundary() bool {
	if p.index+2 >= len(p.tokens) {
		return false
	}
	return p.tokens[p.index+1].Kind == lexer.KindIdentifier && p.tokens[p.index+2].Kind == lexer.KindColon
}

func (p *binderParser) error(span diag.Span, code string, message string) {
	p.diags = append(p.diags, diag.Diagnostic{Code: code, Message: message, Severity: diag.SeverityError, Span: span})
}

func (p *binderParser) finish(ok bool) bool {
	if p.owner != nil {
		p.owner.diags = append(p.owner.diags, p.diags...)
	}
	return ok && len(p.diags) == 0
}

func spanJoin(left diag.Span, right diag.Span) diag.Span {
	return diag.Span{Start: left.Start, End: right.End}
}

func renderTokens(tokens []lexer.Token) string {
	parts := make([]string, 0, len(tokens))
	for _, token := range tokens {
		parts = append(parts, token.Lexeme)
	}
	return strings.Join(parts, " ")
}
