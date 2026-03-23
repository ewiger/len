package validator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yy/len/internal/ast"
	"github.com/yy/len/internal/diag"
	"github.com/yy/len/internal/loader"
	"github.com/yy/len/internal/quasi"
)

// Validator performs semantic checks over a parsed program graph.
type Validator struct {
	ProfileDir        string
	DefaultQuasiStyle string
}

// Validate checks duplicates, binder scope, arity, type references, and quasi blocks.
func (v Validator) Validate(program loader.Program) []diag.Diagnostic {
	table := newSymbolTable()
	diags := append([]diag.Diagnostic{}, program.Diagnostics...)

	for _, unit := range program.Units {
		for _, decl := range unit.File.Decls {
			diags = append(diags, registerDeclSymbols(table, decl)...)
		}
	}

	for _, unit := range program.Units {
		for _, decl := range unit.File.Decls {
			diags = append(diags, v.validateDecl(table, decl)...)
		}
	}
	diag.Sort(diags)
	return diags
}

func registerDeclSymbols(table symbolTable, decl ast.Decl) []diag.Diagnostic {
	switch typed := decl.(type) {
	case *ast.TypeDecl:
		return insertSymbol(table, namespaceType, typed.Name, 0, typed)
	case *ast.StructDecl:
		return insertSymbol(table, namespaceType, typed.Name, len(typed.Fields), typed)
	case *ast.ContractDecl:
		diags := insertSymbol(table, namespaceContract, typed.Name, len(typed.Params), typed)
		for _, member := range typed.Members {
			diags = append(diags, registerDeclSymbols(table, member)...)
		}
		return diags
	case *ast.RelDecl:
		return insertSymbol(table, namespaceRel, typed.Name, len(typed.Params), typed)
	case *ast.FnDecl:
		return insertSymbol(table, namespaceFn, typed.Name, len(typed.Params), typed)
	case *ast.ConstDecl:
		return insertSymbol(table, namespaceConst, typed.Name, 0, typed)
	case *ast.TraitDecl:
		return insertSymbol(table, namespaceTrait, typed.Name, 0, typed)
	case *ast.KeywordDecl:
		return insertSymbol(table, namespaceKeyword, typed.Name, 0, typed)
	case *ast.SymbolDecl:
		return insertSymbol(table, namespaceSymbol, typed.Name, 0, typed)
	case *ast.SpecDecl:
		return insertSymbol(table, namespaceSpec, typed.Name, 0, typed)
	default:
		return nil
	}
}

func insertSymbol(table symbolTable, space namespace, name string, arity int, decl ast.Decl) []diag.Diagnostic {
	_, ok := table.put(space, name, symbol{Name: name, Arity: arity, Decl: decl, Space: space})
	if ok {
		return nil
	}
	return []diag.Diagnostic{{
		Code:     "validator.duplicate.decl",
		Message:  fmt.Sprintf("duplicate %s declaration %q", space, name),
		Severity: diag.SeverityError,
		Span:     decl.GetSpan(),
	}}
}

func (v Validator) validateDecl(table symbolTable, decl ast.Decl) []diag.Diagnostic {
	return v.validateDeclWithScope(table, decl, nil)
}

func (v Validator) validateDeclWithScope(table symbolTable, decl ast.Decl, scope map[string]struct{}) []diag.Diagnostic {
	switch typed := decl.(type) {
	case *ast.RelDecl:
		return v.validateBinders(table, typed.Params, scope)
	case *ast.StructDecl:
		return v.validateStructDecl(table, typed)
	case *ast.ContractDecl:
		return v.validateContractDecl(table, typed, scope)
	case *ast.ConstDecl:
		return v.validateTypeExpr(table, typed.Type, scope)
	case *ast.ImplDecl:
		diags := v.validateTypeExpr(table, typed.Type, scope)
		if _, ok := table.get(namespaceTrait, typed.TraitName); !ok {
			diags = append(diags, diag.Diagnostic{Code: "validator.trait.unresolved", Message: fmt.Sprintf("unknown trait %q", typed.TraitName), Severity: diag.SeverityError, Span: typed.Span})
		}
		return diags
	case *ast.SyntaxDecl:
		diags := v.validateBinders(table, typed.Binders, scope)
		localScope := bindScope(scope, typed.Binders)
		diags = append(diags, v.validateExpr(table, typed.Surface, localScope)...)
		diags = append(diags, v.validateExpr(table, typed.Canonical, localScope)...)
		return diags
	case *ast.SpecDecl:
		diags := v.validateBinders(table, typed.Given, scope)
		if typed.Must == nil {
			diags = append(diags, diag.Diagnostic{Code: "validator.spec.must", Message: "spec requires a must clause", Severity: diag.SeverityError, Span: typed.Span})
			return diags
		}
		return append(diags, v.validateExpr(table, typed.Must, bindScope(scope, typed.Given))...)
	case *ast.FnDecl:
		diags := v.validateBinders(table, typed.Params, scope)
		localScope := bindScope(scope, typed.Params)
		if typed.Result != nil {
			diags = append(diags, v.validateBinders(table, []ast.Binder{*typed.Result}, localScope)...)
			localScope = bindScope(localScope, []ast.Binder{*typed.Result})
		}
		for _, clause := range typed.Clauses {
			switch c := clause.(type) {
			case *ast.RequiresClause:
				diags = append(diags, v.validateExpr(table, c.Formula, localScope)...)
			case *ast.EnsuresClause:
				diags = append(diags, v.validateExpr(table, c.Formula, localScope)...)
			case *ast.ImplementsClause:
				diags = append(diags, v.validateExpr(table, c.Formula, localScope)...)
			}
		}
		if typed.Quasi != nil {
			diags = append(diags, v.validateQuasi(typed.Quasi)...)
		}
		return diags
	}
	return nil
}

func (v Validator) validateStructDecl(table symbolTable, decl *ast.StructDecl) []diag.Diagnostic {
	diags := make([]diag.Diagnostic, 0)
	seen := map[string]struct{}{}
	for _, field := range decl.Fields {
		if _, exists := seen[field.Name]; exists {
			diags = append(diags, diag.Diagnostic{Code: "validator.field.duplicate", Message: fmt.Sprintf("duplicate field %q", field.Name), Severity: diag.SeverityError, Span: field.Span})
		} else {
			seen[field.Name] = struct{}{}
		}
		diags = append(diags, v.validateTypeExpr(table, field.Type, nil)...)
	}
	return diags
}

func (v Validator) validateContractDecl(table symbolTable, decl *ast.ContractDecl, scope map[string]struct{}) []diag.Diagnostic {
	diags := v.validateBinders(table, decl.Params, scope)
	contractScope := bindScope(scope, decl.Params)
	memberSeen := map[string]diag.Span{}
	for _, member := range decl.Members {
		if name := declName(member); name != "" {
			if span, exists := memberSeen[name]; exists {
				diags = append(diags, diag.Diagnostic{Code: "validator.contract.member.duplicate", Message: fmt.Sprintf("duplicate contract member %q", name), Severity: diag.SeverityError, Span: member.GetSpan()})
				diags = append(diags, diag.Diagnostic{Code: "validator.contract.member.duplicate", Message: fmt.Sprintf("duplicate contract member %q", name), Severity: diag.SeverityError, Span: span})
			} else {
				memberSeen[name] = member.GetSpan()
			}
		}
		diags = append(diags, v.validateDeclWithScope(table, member, contractScope)...)
	}
	return diags
}

func (v Validator) validateQuasi(quasiClause *ast.QuasiClause) []diag.Diagnostic {
	styleName := quasiClause.StyleName
	if styleName == "" {
		styleName = v.DefaultQuasiStyle
		if styleName == "" {
			styleName = "ProceduralAlgorithm"
		}
	}
	profilePath := filepath.Join(v.ProfileDir, styleFileName(styleName))
	profile, err := quasi.LoadProfile(profilePath)
	if err != nil {
		return []diag.Diagnostic{{
			Code:     "validator.quasi.profile",
			Message:  err.Error(),
			Severity: diag.SeverityError,
			Span:     quasiClause.HeaderSpan,
		}}
	}
	lines := make([]quasi.RawLine, 0, len(quasiClause.Block.Lines))
	for _, line := range quasiClause.Block.Lines {
		lines = append(lines, quasi.RawLine{
			Text:         line.Text,
			TrimmedText:  line.TrimmedText,
			IndentColumn: line.IndentColumn,
			Span: quasi.SourceSpan{
				Start: quasi.SourcePosition{File: line.LineSpan.Start.File, Line: line.LineSpan.Start.Line, Column: line.LineSpan.Start.Column},
				End:   quasi.SourcePosition{File: line.LineSpan.End.File, Line: line.LineSpan.End.Line, Column: line.LineSpan.End.Column},
			},
		})
	}
	result := quasi.Validator{Profile: profile}.Validate(quasi.Block{
		StyleName: styleName,
		Lines:     lines,
		Span: quasi.SourceSpan{
			Start: quasi.SourcePosition{File: quasiClause.Block.Span.Start.File, Line: quasiClause.Block.Span.Start.Line, Column: quasiClause.Block.Span.Start.Column},
			End:   quasi.SourcePosition{File: quasiClause.Block.Span.End.File, Line: quasiClause.Block.Span.End.Line, Column: quasiClause.Block.Span.End.Column},
		},
	})
	diags := make([]diag.Diagnostic, 0, len(result.Diagnostics))
	for _, finding := range result.Diagnostics {
		diags = append(diags, diag.Diagnostic{
			Code:     finding.Code,
			Message:  finding.Message,
			Severity: diag.SeverityError,
			Span: diag.Span{
				Start: diag.Position{File: finding.Span.Start.File, Line: finding.Span.Start.Line, Column: finding.Span.Start.Column},
				End:   diag.Position{File: finding.Span.End.File, Line: finding.Span.End.Line, Column: finding.Span.End.Column},
			},
		})
	}
	return diags
}

func styleFileName(styleName string) string {
	var parts []rune
	for i, r := range styleName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			parts = append(parts, '-')
		}
		parts = append(parts, rune(strings.ToLower(string(r))[0]))
	}
	return string(parts) + ".quasi-style.yaml"
}

func (v Validator) validateBinders(table symbolTable, binders []ast.Binder, prior map[string]struct{}) []diag.Diagnostic {
	diags := make([]diag.Diagnostic, 0)
	seen := map[string]struct{}{}
	for key := range prior {
		seen[key] = struct{}{}
	}
	for _, binder := range binders {
		if _, exists := seen[binder.Name]; exists {
			diags = append(diags, diag.Diagnostic{Code: "validator.binder.duplicate", Message: fmt.Sprintf("duplicate binder %q", binder.Name), Severity: diag.SeverityError, Span: binder.Span})
			continue
		}
		seen[binder.Name] = struct{}{}
		diags = append(diags, v.validateTypeExpr(table, binder.Type, seen)...)
	}
	return diags
}

func bindScope(base map[string]struct{}, binders []ast.Binder) map[string]struct{} {
	scope := map[string]struct{}{}
	for key := range base {
		scope[key] = struct{}{}
	}
	for _, binder := range binders {
		scope[binder.Name] = struct{}{}
	}
	return scope
}

func (v Validator) validateTypeExpr(table symbolTable, expr ast.Expr, scope map[string]struct{}) []diag.Diagnostic {
	if expr == nil {
		return nil
	}
	name := exprName(expr)
	if name == "" {
		return []diag.Diagnostic{{Code: "validator.type.invalid", Message: "type expression must resolve to a name", Severity: diag.SeverityError, Span: expr.GetSpan()}}
	}
	if _, ok := scope[name]; ok {
		return nil
	}
	if _, ok := table.get(namespaceType, name); !ok {
		return []diag.Diagnostic{{Code: "validator.type.unresolved", Message: fmt.Sprintf("unknown type %q", name), Severity: diag.SeverityError, Span: expr.GetSpan()}}
	}
	return nil
}

func (v Validator) validateExpr(table symbolTable, expr ast.Expr, scope map[string]struct{}) []diag.Diagnostic {
	switch typed := expr.(type) {
	case *ast.IdentExpr:
		if _, ok := scope[typed.Name]; ok {
			return nil
		}
		if _, ok := table.get(namespaceConst, typed.Name); ok {
			return nil
		}
		if _, ok := table.get(namespaceFn, typed.Name); ok {
			return nil
		}
		if _, ok := table.get(namespaceRel, typed.Name); ok {
			return nil
		}
		return []diag.Diagnostic{{Code: "validator.name.unresolved", Message: fmt.Sprintf("unknown name %q", typed.Name), Severity: diag.SeverityError, Span: typed.Span}}
	case *ast.QualifiedExpr:
		name := typed.Parts[len(typed.Parts)-1]
		if _, ok := table.get(namespaceConst, name); ok {
			return nil
		}
		if _, ok := table.get(namespaceFn, name); ok {
			return nil
		}
		if _, ok := table.get(namespaceRel, name); ok {
			return nil
		}
		if _, ok := table.get(namespaceType, name); ok {
			return nil
		}
		return []diag.Diagnostic{{Code: "validator.name.unresolved", Message: fmt.Sprintf("unknown qualified name %q", strings.Join(typed.Parts, ".")), Severity: diag.SeverityError, Span: typed.Span}}
	case *ast.ApplyExpr:
		diags := make([]diag.Diagnostic, 0)
		name := exprName(typed.Callee)
		if sym, ok := table.get(namespaceRel, name); ok {
			if len(typed.Args) != sym.Arity {
				diags = append(diags, diag.Diagnostic{Code: "validator.arity.rel", Message: fmt.Sprintf("relation %q expects %d arguments, got %d", name, sym.Arity, len(typed.Args)), Severity: diag.SeverityError, Span: typed.Span})
			}
		} else if sym, ok := table.get(namespaceFn, name); ok {
			if len(typed.Args) != sym.Arity {
				diags = append(diags, diag.Diagnostic{Code: "validator.arity.fn", Message: fmt.Sprintf("function %q expects %d arguments, got %d", name, sym.Arity, len(typed.Args)), Severity: diag.SeverityError, Span: typed.Span})
			}
		} else if sym, ok := table.get(namespaceType, name); ok {
			if _, isStruct := sym.Decl.(*ast.StructDecl); isStruct {
				if len(typed.Args) != sym.Arity {
					diags = append(diags, diag.Diagnostic{Code: "validator.arity.struct", Message: fmt.Sprintf("struct %q expects %d constructor arguments, got %d", name, sym.Arity, len(typed.Args)), Severity: diag.SeverityError, Span: typed.Span})
				}
			} else {
				diags = append(diags, diag.Diagnostic{Code: "validator.name.unresolved", Message: fmt.Sprintf("unknown callable %q", name), Severity: diag.SeverityError, Span: typed.Callee.GetSpan()})
			}
		} else {
			diags = append(diags, diag.Diagnostic{Code: "validator.name.unresolved", Message: fmt.Sprintf("unknown callable %q", name), Severity: diag.SeverityError, Span: typed.Callee.GetSpan()})
		}
		for _, arg := range typed.Args {
			diags = append(diags, v.validateExpr(table, arg, scope)...)
		}
		return diags
	case *ast.BinaryExpr:
		diags := v.validateExpr(table, typed.Left, scope)
		return append(diags, v.validateExpr(table, typed.Right, scope)...)
	case *ast.UnaryExpr:
		return v.validateExpr(table, typed.Expr, scope)
	case *ast.GroupExpr:
		return v.validateExpr(table, typed.Inner, scope)
	case *ast.QuantifiedExpr:
		diags := v.validateBinders(table, typed.Binders, scope)
		inner := bindScope(scope, typed.Binders)
		return append(diags, v.validateExpr(table, typed.Body, inner)...)
	case *ast.NumberExpr, *ast.StringExpr, *ast.BoolExpr:
		return nil
	default:
		return nil
	}
}

func declName(decl ast.Decl) string {
	switch typed := decl.(type) {
	case *ast.TypeDecl:
		return typed.Name
	case *ast.StructDecl:
		return typed.Name
	case *ast.ContractDecl:
		return typed.Name
	case *ast.RelDecl:
		return typed.Name
	case *ast.ConstDecl:
		return typed.Name
	case *ast.KeywordDecl:
		return typed.Name
	case *ast.SymbolDecl:
		return typed.Name
	case *ast.TraitDecl:
		return typed.Name
	case *ast.SpecDecl:
		return typed.Name
	case *ast.FnDecl:
		return typed.Name
	default:
		return ""
	}
}

func exprName(expr ast.Expr) string {
	switch typed := expr.(type) {
	case *ast.IdentExpr:
		return typed.Name
	case *ast.QualifiedExpr:
		if len(typed.Parts) == 0 {
			return ""
		}
		return typed.Parts[len(typed.Parts)-1]
	default:
		return ""
	}
}
