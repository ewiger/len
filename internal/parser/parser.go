package parser

import (
	"strings"

	"github.com/yy/len/internal/ast"
	"github.com/yy/len/internal/diag"
)

// Parse reads one len.l1 source file into an AST plus parser diagnostics.
func Parse(filePath string, src string) (*ast.File, []diag.Diagnostic) {
	p := newParser(filePath, src)
	return p.parseFile()
}

type parser struct {
	filePath string
	lines    []sourceLine
	index    int
	diags    []diag.Diagnostic
}

type sourceLine struct {
	raw     string
	trimmed string
	indent  int
	lineNo  int
	blank   bool
}

func newParser(filePath string, src string) *parser {
	return &parser{filePath: filePath, lines: preprocessSource(src)}
}

func (p *parser) parseFile() (*ast.File, []diag.Diagnostic) {
	decls := make([]ast.Decl, 0)
	for p.index < len(p.lines) {
		line, ok := p.currentSignificant()
		if !ok {
			break
		}
		if line.indent != 0 {
			p.error(lineSpan(p.filePath, line), "parser.unexpected-indent", "top-level declaration must start at column 1")
			p.index++
			continue
		}
		decl := p.parseTopLevel()
		if decl != nil {
			decls = append(decls, decl)
		}
	}

	span := diag.Span{}
	if len(p.lines) > 0 {
		span.Start = diag.Position{File: p.filePath, Line: 1, Column: 1}
		last := p.lines[len(p.lines)-1]
		span.End = diag.Position{File: p.filePath, Line: last.lineNo, Column: len(last.raw) + 1}
	}
	return &ast.File{Decls: decls, Span: span}, p.diags
}

func preprocessSource(src string) []sourceLine {
	rawLines := strings.Split(src, "\n")
	result := make([]sourceLine, 0, len(rawLines))
	inBlockComment := false
	inDocstring := false

	for i, raw := range rawLines {
		trimmed := strings.TrimSpace(raw)
		if inBlockComment {
			if strings.Contains(raw, "#>") {
				inBlockComment = false
			}
			continue
		}
		if inDocstring {
			if strings.Contains(raw, "\"\"\"") {
				inDocstring = false
			}
			continue
		}
		if strings.HasPrefix(trimmed, "<#") {
			if !strings.Contains(trimmed, "#>") {
				inBlockComment = true
			}
			continue
		}
		if strings.HasPrefix(trimmed, "\"\"\"") {
			if strings.Count(trimmed, "\"\"\"") < 2 {
				inDocstring = true
			}
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(trimmed, "\"") && strings.HasSuffix(trimmed, "\"") && !strings.Contains(trimmed[1:len(trimmed)-1], "\"") {
			continue
		}

		result = append(result, sourceLine{
			raw:     raw,
			trimmed: trimmed,
			indent:  countIndent(raw),
			lineNo:  i + 1,
			blank:   trimmed == "",
		})
	}

	return result
}

func countIndent(raw string) int {
	count := 0
	for _, r := range raw {
		switch r {
		case ' ':
			count++
		case '\t':
			count += 4
		default:
			return count
		}
	}
	return count
}

func (p *parser) currentSignificant() (sourceLine, bool) {
	for p.index < len(p.lines) {
		if p.lines[p.index].blank {
			p.index++
			continue
		}
		return p.lines[p.index], true
	}
	return sourceLine{}, false
}

func (p *parser) error(span diag.Span, code string, message string) {
	p.diags = append(p.diags, diag.Diagnostic{Code: code, Message: message, Severity: diag.SeverityError, Span: span})
}

func lineSpan(filePath string, line sourceLine) diag.Span {
	return diag.Span{
		Start: diag.Position{File: filePath, Line: line.lineNo, Column: line.indent + 1},
		End:   diag.Position{File: filePath, Line: line.lineNo, Column: len(line.raw) + 1},
	}
}

func spanFrom(filePath string, start sourceLine, end sourceLine) diag.Span {
	return diag.Span{
		Start: diag.Position{File: filePath, Line: start.lineNo, Column: start.indent + 1},
		End:   diag.Position{File: filePath, Line: end.lineNo, Column: len(end.raw) + 1},
	}
}

func spanEndLine(span diag.Span, fallback sourceLine) sourceLine {
	line := fallback
	if span.End.Line > 0 {
		line.lineNo = span.End.Line
		line.indent = 0
		line.raw = strings.Repeat(" ", max(span.End.Column-1, 0))
	}
	return line
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func trimKeyword(text string, keyword string) string {
	return strings.TrimSpace(strings.TrimPrefix(text, keyword))
}

func splitOnce(text string, sep string) (string, string, bool) {
	idx := strings.Index(text, sep)
	if idx < 0 {
		return "", "", false
	}
	return text[:idx], text[idx+len(sep):], true
}
