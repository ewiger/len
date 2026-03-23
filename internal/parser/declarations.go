package parser

import (
	"fmt"
	"strings"

	"github.com/yy/len/internal/ast"
	"github.com/yy/len/internal/diag"
)

func (p *parser) parseTopLevel() ast.Decl {
	line := p.lines[p.index]
	keyword := leadingWord(line.trimmed)
	switch keyword {
	case "import":
		return p.parseImportDecl(line)
	case "type":
		return p.parseTypeDecl(line)
	case "struct":
		return p.parseStructDecl(line)
	case "rel":
		return p.parseRelDecl(line)
	case "const":
		return p.parseConstDecl(line)
	case "keyword":
		return p.parseKeywordDecl(line)
	case "symbol":
		return p.parseSymbolDecl(line)
	case "trait":
		return p.parseTraitDecl(line)
	case "impl":
		return p.parseImplDecl(line)
	case "syntax":
		return p.parseSyntaxDecl(line)
	case "spec":
		return p.parseSpecDecl(line)
	case "contract":
		return p.parseContractDecl(line)
	case "fn":
		return p.parseFnDecl(line)
	default:
		p.error(lineSpan(p.filePath, line), "parser.decl.unknown", fmt.Sprintf("unknown top-level declaration %q", keyword))
		p.index++
		return nil
	}
}

func (p *parser) parseImportDecl(line sourceLine) ast.Decl {
	p.index++
	module := strings.TrimSpace(strings.TrimPrefix(line.trimmed, "import"))
	parts := splitModulePath(module)
	if len(parts) == 0 {
		p.error(lineSpan(p.filePath, line), "parser.import.invalid", "import declaration requires a module path")
		return nil
	}
	return &ast.ImportDecl{ModulePath: parts, Span: lineSpan(p.filePath, line)}
}

func (p *parser) parseTypeDecl(line sourceLine) ast.Decl {
	p.index++
	name := trimKeyword(line.trimmed, "type")
	if name == "" {
		p.error(lineSpan(p.filePath, line), "parser.type.invalid", "type declaration requires a name")
		return nil
	}
	return &ast.TypeDecl{Name: name, Span: lineSpan(p.filePath, line)}
}

func (p *parser) parseStructDecl(header sourceLine) ast.Decl {
	name := trimKeyword(header.trimmed, "struct")
	if name == "" {
		p.error(lineSpan(p.filePath, header), "parser.struct.invalid", "struct declaration requires a name")
		p.index++
		return nil
	}

	p.index++
	decl := &ast.StructDecl{Name: name}
	bodyIndent, ok := p.hasIndentedBlock(header)
	if !ok {
		decl.Span = lineSpan(p.filePath, header)
		return decl
	}

	end := header
	for p.index < len(p.lines) {
		line := p.lines[p.index]
		if line.blank {
			p.index++
			continue
		}
		if line.indent < bodyIndent {
			break
		}
		if line.indent > bodyIndent {
			p.error(lineSpan(p.filePath, line), "parser.struct.indent", "struct fields must align with the first field")
			end = line
			p.index++
			continue
		}

		nameText, typeText, ok := splitOnce(line.trimmed, ":")
		if !ok {
			p.error(lineSpan(p.filePath, line), "parser.struct.field", "struct field must use `name: Type`")
			end = line
			p.index++
			continue
		}
		fieldName := strings.TrimSpace(nameText)
		if fieldName == "" {
			p.error(lineSpan(p.filePath, line), "parser.struct.field", "struct field requires a name")
			end = line
			p.index++
			continue
		}
		typeExpr, parsed := p.parseInlineExpr(strings.TrimSpace(typeText), diag.Position{File: p.filePath, Line: line.lineNo, Column: line.indent + strings.Index(line.raw, ":") + 2})
		if parsed {
			decl.Fields = append(decl.Fields, ast.FieldDecl{Name: fieldName, Type: typeExpr, Span: lineSpan(p.filePath, line)})
		}
		end = line
		p.index++
	}
	decl.Span = spanFrom(p.filePath, header, end)
	return decl
}

func (p *parser) parseRelDecl(line sourceLine) ast.Decl {
	p.index++
	head := trimKeyword(line.trimmed, "rel")
	open := strings.Index(head, "(")
	close := strings.LastIndex(head, ")")
	if open <= 0 || close < open {
		p.error(lineSpan(p.filePath, line), "parser.rel.invalid", "rel declaration requires a name and parameter list")
		return nil
	}
	name := strings.TrimSpace(head[:open])
	binderText := head[open+1 : close]
	params, ok := p.parseBindersText(binderText, diag.Position{File: p.filePath, Line: line.lineNo, Column: line.indent + open + 2})
	if !ok {
		return nil
	}
	return &ast.RelDecl{Name: name, Params: params, Span: lineSpan(p.filePath, line)}
}

func (p *parser) parseConstDecl(line sourceLine) ast.Decl {
	p.index++
	content := trimKeyword(line.trimmed, "const")
	name, typeText, ok := splitOnce(content, ":")
	if !ok {
		p.error(lineSpan(p.filePath, line), "parser.const.invalid", "const declaration must use `const Name: Type`")
		return nil
	}
	typeExpr, parsed := p.parseInlineExpr(strings.TrimSpace(typeText), diag.Position{File: p.filePath, Line: line.lineNo, Column: line.indent + strings.Index(line.raw, ":") + 2})
	if !parsed {
		return nil
	}
	return &ast.ConstDecl{Name: strings.TrimSpace(name), Type: typeExpr, Span: lineSpan(p.filePath, line)}
}

func (p *parser) parseKeywordDecl(line sourceLine) ast.Decl {
	p.index++
	name := trimKeyword(line.trimmed, "keyword")
	if name == "" {
		p.error(lineSpan(p.filePath, line), "parser.keyword.invalid", "keyword declaration requires a name")
		return nil
	}
	return &ast.KeywordDecl{Name: name, Span: lineSpan(p.filePath, line)}
}

func (p *parser) parseSymbolDecl(line sourceLine) ast.Decl {
	p.index++
	content := trimKeyword(line.trimmed, "symbol")
	left, right, ok := splitOnce(content, " as ")
	if !ok {
		p.error(lineSpan(p.filePath, line), "parser.symbol.invalid", "symbol declaration must use `symbol Name as \"value\"`")
		return nil
	}
	return &ast.SymbolDecl{Name: strings.TrimSpace(left), Value: strings.TrimSpace(right), Span: lineSpan(p.filePath, line)}
}

func (p *parser) parseTraitDecl(line sourceLine) ast.Decl {
	p.index++
	name := trimKeyword(line.trimmed, "trait")
	if name == "" {
		p.error(lineSpan(p.filePath, line), "parser.trait.invalid", "trait declaration requires a name")
		return nil
	}
	return &ast.TraitDecl{Name: name, Span: lineSpan(p.filePath, line)}
}

func (p *parser) parseImplDecl(line sourceLine) ast.Decl {
	p.index++
	content := trimKeyword(line.trimmed, "impl")
	left, right, ok := splitOnce(content, ":")
	if !ok {
		p.error(lineSpan(p.filePath, line), "parser.impl.invalid", "impl declaration must use `impl TypeExpr : TraitName`")
		return nil
	}
	typeExpr, parsed := p.parseInlineExpr(strings.TrimSpace(left), diag.Position{File: p.filePath, Line: line.lineNo, Column: line.indent + 6})
	if !parsed {
		return nil
	}
	return &ast.ImplDecl{Type: typeExpr, TraitName: strings.TrimSpace(right), Span: lineSpan(p.filePath, line)}
}

func (p *parser) parseSyntaxDecl(line sourceLine) ast.Decl {
	p.index++
	content := trimKeyword(line.trimmed, "syntax")
	surfaceText, rest, ok := splitOnce(content, " where ")
	if !ok {
		p.error(lineSpan(p.filePath, line), "parser.syntax.invalid", "syntax declaration requires `where` and `implies`")
		return nil
	}
	binderText, canonicalText, ok := splitOnce(rest, " implies ")
	if !ok {
		p.error(lineSpan(p.filePath, line), "parser.syntax.invalid", "syntax declaration requires `implies`")
		return nil
	}
	surface, okSurface := p.parseInlineExpr(strings.TrimSpace(surfaceText), diag.Position{File: p.filePath, Line: line.lineNo, Column: line.indent + 8})
	binders, okBinders := p.parseBindersText(strings.TrimSpace(binderText), diag.Position{File: p.filePath, Line: line.lineNo, Column: line.indent + 8 + len(surfaceText) + len(" where ")})
	canonical, okCanonical := p.parseInlineExpr(strings.TrimSpace(canonicalText), diag.Position{File: p.filePath, Line: line.lineNo, Column: line.indent + 8 + len(surfaceText) + len(" where ") + len(binderText) + len(" implies ")})
	if !okSurface || !okBinders || !okCanonical {
		return nil
	}
	return &ast.SyntaxDecl{Surface: surface, Binders: binders, Canonical: canonical, Span: lineSpan(p.filePath, line)}
}

func (p *parser) parseSpecDecl(header sourceLine) ast.Decl {
	name := trimKeyword(header.trimmed, "spec")
	if name == "" {
		p.error(lineSpan(p.filePath, header), "parser.spec.invalid", "spec declaration requires a name")
		p.index++
		return nil
	}
	p.index++
	bodyIndent, ok := p.expectIndentedBlock(header)
	if !ok {
		return &ast.SpecDecl{Name: name, Span: lineSpan(p.filePath, header)}
	}

	decl := &ast.SpecDecl{Name: name}
	end := header
	for p.index < len(p.lines) {
		line := p.lines[p.index]
		if line.blank {
			p.index++
			continue
		}
		if line.indent < bodyIndent {
			break
		}
		if line.indent > bodyIndent {
			p.error(lineSpan(p.filePath, line), "parser.spec.indent", "spec clauses must align with the first clause")
			p.index++
			continue
		}

		switch leadingWord(line.trimmed) {
		case "given":
			binders, ok := p.parseBindersText(trimKeyword(line.trimmed, "given"), diag.Position{File: p.filePath, Line: line.lineNo, Column: line.indent + 7})
			if ok {
				decl.Given = append(decl.Given, binders...)
			}
			end = line
			p.index++
		case "must":
			exprText, lastLine := p.collectContinuationText(line, "must")
			expr, parsed := p.parseInlineExpr(exprText, diag.Position{File: p.filePath, Line: line.lineNo, Column: line.indent + 6})
			if parsed {
				decl.Must = expr
			}
			end = lastLine
		default:
			p.error(lineSpan(p.filePath, line), "parser.spec.clause", "spec body only supports `given` and `must` clauses")
			end = line
			p.index++
		}
	}
	decl.Span = spanFrom(p.filePath, header, end)
	return decl
}

func (p *parser) parseContractDecl(header sourceLine) ast.Decl {
	content := trimKeyword(header.trimmed, "contract")
	name := strings.TrimSpace(content)
	params := []ast.Binder(nil)
	if open := strings.Index(content, "("); open >= 0 {
		close := strings.LastIndex(content, ")")
		if close < open {
			p.error(lineSpan(p.filePath, header), "parser.contract.invalid", "contract declaration has an unclosed parameter list")
			p.index++
			return nil
		}
		name = strings.TrimSpace(content[:open])
		parsedParams, ok := p.parseBindersText(content[open+1:close], diag.Position{File: p.filePath, Line: header.lineNo, Column: header.indent + open + len("contract ") + 1})
		if !ok {
			p.index++
			return nil
		}
		params = parsedParams
		if strings.TrimSpace(content[close+1:]) != "" {
			p.error(lineSpan(p.filePath, header), "parser.contract.invalid", "unexpected trailing content after contract parameters")
		}
	}
	if name == "" {
		p.error(lineSpan(p.filePath, header), "parser.contract.invalid", "contract declaration requires a name")
		p.index++
		return nil
	}

	p.index++
	bodyIndent, ok := p.expectIndentedBlock(header)
	decl := &ast.ContractDecl{Name: name, Params: params}
	if !ok {
		decl.Span = lineSpan(p.filePath, header)
		return decl
	}

	end := header
	for p.index < len(p.lines) {
		line := p.lines[p.index]
		if line.blank {
			p.index++
			continue
		}
		if line.indent < bodyIndent {
			break
		}
		if line.indent > bodyIndent {
			p.error(lineSpan(p.filePath, line), "parser.contract.indent", "contract members must align with the first member")
			end = line
			p.index++
			continue
		}

		var member ast.Decl
		switch leadingWord(line.trimmed) {
		case "rel":
			member = p.parseRelDecl(line)
		case "spec":
			member = p.parseSpecDecl(line)
		case "fn":
			member = p.parseFnDecl(line)
		default:
			p.error(lineSpan(p.filePath, line), "parser.contract.member", "contract body only supports rel, spec, and fn declarations")
			end = line
			p.index++
			continue
		}
		if member != nil {
			decl.Members = append(decl.Members, member)
			end = spanEndLine(member.GetSpan(), line)
		}
	}
	decl.Span = spanFrom(p.filePath, header, end)
	return decl
}

func (p *parser) parseFnDecl(header sourceLine) ast.Decl {
	content := trimKeyword(header.trimmed, "fn")
	open := strings.Index(content, "(")
	close := strings.LastIndex(content, ")")
	if open <= 0 || close < open {
		p.error(lineSpan(p.filePath, header), "parser.fn.invalid", "fn declaration requires a name and parameter list")
		p.index++
		return nil
	}
	name := strings.TrimSpace(content[:open])
	params, ok := p.parseBindersText(content[open+1:close], diag.Position{File: p.filePath, Line: header.lineNo, Column: header.indent + open + 4})
	if !ok {
		p.index++
		return nil
	}

	var result *ast.Binder
	rest := strings.TrimSpace(content[close+1:])
	if rest != "" {
		if !strings.HasPrefix(rest, "->") {
			p.error(lineSpan(p.filePath, header), "parser.fn.result", "fn result binder must follow `->`")
		} else {
			binders, ok := p.parseBindersText(strings.TrimSpace(strings.TrimPrefix(rest, "->")), diag.Position{File: p.filePath, Line: header.lineNo, Column: header.indent + close + 4})
			if ok && len(binders) == 1 {
				result = &binders[0]
			} else if ok {
				p.error(lineSpan(p.filePath, header), "parser.fn.result", "fn result must declare exactly one binder")
			}
		}
	}

	p.index++
	bodyIndent, bodyOK := p.expectIndentedBlock(header)
	decl := &ast.FnDecl{Name: name, Params: params, Result: result}
	if !bodyOK {
		decl.Span = lineSpan(p.filePath, header)
		return decl
	}

	end := header
	sawQuasi := false
	for p.index < len(p.lines) {
		line := p.lines[p.index]
		if line.blank {
			p.index++
			continue
		}
		if line.indent < bodyIndent {
			break
		}
		if line.indent > bodyIndent {
			p.error(lineSpan(p.filePath, line), "parser.fn.indent", "fn clauses must align with the first clause")
			p.index++
			continue
		}
		if sawQuasi {
			p.error(lineSpan(p.filePath, line), "parser.fn.quasi.order", "quasi clause must be the final clause in an fn body")
			end = line
			p.index++
			continue
		}

		switch leadingWord(line.trimmed) {
		case "requires":
			exprText, lastLine := p.collectContinuationText(line, "requires")
			expr, parsed := p.parseInlineExpr(exprText, diag.Position{File: p.filePath, Line: line.lineNo, Column: line.indent + 10})
			if parsed {
				decl.Clauses = append(decl.Clauses, &ast.RequiresClause{Formula: expr, Span: spanFrom(p.filePath, line, lastLine)})
			}
			end = lastLine
		case "ensures":
			exprText, lastLine := p.collectContinuationText(line, "ensures")
			expr, parsed := p.parseInlineExpr(exprText, diag.Position{File: p.filePath, Line: line.lineNo, Column: line.indent + 9})
			if parsed {
				decl.Clauses = append(decl.Clauses, &ast.EnsuresClause{Formula: expr, Span: spanFrom(p.filePath, line, lastLine)})
			}
			end = lastLine
		case "implements":
			exprText, lastLine := p.collectContinuationText(line, "implements")
			expr, parsed := p.parseInlineExpr(exprText, diag.Position{File: p.filePath, Line: line.lineNo, Column: line.indent + 12})
			if parsed {
				decl.Clauses = append(decl.Clauses, &ast.ImplementsClause{Formula: expr, Span: spanFrom(p.filePath, line, lastLine)})
			}
			end = lastLine
		case "quasi":
			quasi, lastLine := p.parseQuasiClause(line)
			decl.Quasi = quasi
			sawQuasi = true
			end = lastLine
		default:
			p.error(lineSpan(p.filePath, line), "parser.fn.clause", "fn body only supports requires, ensures, implements, and quasi clauses")
			end = line
			p.index++
		}
	}
	if decl.Quasi == nil {
		p.error(lineSpan(p.filePath, header), "parser.fn.quasi.required", "fn declaration requires a quasi clause")
	}

	decl.Span = spanFrom(p.filePath, header, end)
	return decl
}

func (p *parser) parseQuasiClause(header sourceLine) (*ast.QuasiClause, sourceLine) {
	content := trimKeyword(header.trimmed, "quasi")
	styleName := ""
	if content != ":" {
		styleText := strings.TrimSpace(strings.TrimSuffix(content, ":"))
		if styleText != "" {
			styleText = strings.TrimSpace(strings.TrimPrefix(styleText, "using style"))
			styleName = strings.TrimSpace(styleText)
		}
	}
	p.index++
	bodyIndent, ok := p.expectIndentedBlock(header)
	if !ok {
		return &ast.QuasiClause{StyleName: styleName, HeaderSpan: lineSpan(p.filePath, header), Span: lineSpan(p.filePath, header)}, header
	}
	lines := make([]ast.RawQuasiLine, 0)
	end := header
	lastContent := header
	for p.index < len(p.lines) {
		line := p.lines[p.index]
		if !line.blank && line.indent < bodyIndent {
			break
		}
		if line.blank {
			lines = append(lines, ast.RawQuasiLine{
				Text:         line.raw,
				TrimmedText:  "",
				IndentColumn: line.indent,
				LineSpan:     lineSpan(p.filePath, line),
			})
			end = line
			p.index++
			continue
		}
		lines = append(lines, ast.RawQuasiLine{
			Text:         line.raw,
			TrimmedText:  line.trimmed,
			IndentColumn: line.indent,
			LineSpan:     lineSpan(p.filePath, line),
		})
		end = line
		lastContent = line
		p.index++
	}

	for len(lines) > 0 && lines[len(lines)-1].TrimmedText == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) > 0 {
		end = lastContent
	}

	blockSpan := lineSpan(p.filePath, header)
	if len(lines) > 0 {
		blockSpan = diag.Span{Start: lines[0].LineSpan.Start, End: lines[len(lines)-1].LineSpan.End}
	}
	clause := &ast.QuasiClause{
		StyleName:  styleName,
		HeaderSpan: lineSpan(p.filePath, header),
		Block: ast.RawQuasiBlock{
			Lines: lines,
			Span:  blockSpan,
		},
		Span: spanFrom(p.filePath, header, end),
	}
	return clause, end
}

func (p *parser) parseBindersText(text string, start diag.Position) ([]ast.Binder, bool) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, true
	}
	bp := newBinderParser(p, text, start)
	binders, ok := bp.parse()
	return binders, ok
}

func (p *parser) parseInlineExpr(text string, start diag.Position) (ast.Expr, bool) {
	text = strings.TrimSpace(text)
	if text == "" {
		p.error(diag.Span{Start: start, End: start}, "parser.expr.empty", "expected expression")
		return nil, false
	}
	expr, diags, ok := parseExpression(p.filePath, text, start)
	p.diags = append(p.diags, diags...)
	return expr, ok
}

func (p *parser) expectIndentedBlock(header sourceLine) (int, bool) {
	peek := p.index
	for peek < len(p.lines) && p.lines[peek].blank {
		peek++
	}
	if peek >= len(p.lines) || p.lines[peek].indent <= header.indent {
		p.error(lineSpan(p.filePath, header), "parser.block.missing", "expected an indented block")
		return 0, false
	}
	return p.lines[peek].indent, true
}

func (p *parser) hasIndentedBlock(header sourceLine) (int, bool) {
	peek := p.index
	for peek < len(p.lines) && p.lines[peek].blank {
		peek++
	}
	if peek >= len(p.lines) || p.lines[peek].indent <= header.indent {
		return 0, false
	}
	return p.lines[peek].indent, true
}

func (p *parser) collectContinuationText(line sourceLine, keyword string) (string, sourceLine) {
	parts := []string{strings.TrimSpace(strings.TrimPrefix(line.trimmed, keyword))}
	startIndent := line.indent
	last := line
	p.index++
	for p.index < len(p.lines) {
		next := p.lines[p.index]
		if next.blank {
			p.index++
			continue
		}
		if next.indent <= startIndent {
			break
		}
		parts = append(parts, next.trimmed)
		last = next
		p.index++
	}
	return strings.Join(parts, " "), last
}

func leadingWord(text string) string {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return ""
	}
	if len(parts) >= 2 && parts[0] == "else" && parts[1] == "if" {
		return "else if"
	}
	return parts[0]
}

func splitModulePath(path string) []string {
	parts := strings.Split(strings.TrimSpace(path), ".")
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			filtered = append(filtered, part)
		}
	}
	return filtered
}
