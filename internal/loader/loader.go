package loader
package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yy/len/internal/ast"
	"github.com/yy/len/internal/diag"
	"github.com/yy/len/internal/parser"
)

// Unit is one parsed source file plus loader metadata.
type Unit struct {
	Path       string
	ModuleName string
	Source     string
	File       *ast.File
}

// Program is the transitive closure of requested files and imports.
type Program struct {
	Units       []*Unit
	ByPath      map[string]*Unit
	Diagnostics []diag.Diagnostic
}

// Loader resolves and parses source files plus transitive imports.
type Loader struct {
	Root string
}

// LoadPaths parses the provided source paths and their transitive imports.
func (l Loader) LoadPaths(paths []string) Program {
	program := Program{ByPath: map[string]*Unit{}}
	visiting := map[string]bool{}
	for _, path := range paths {
		abs, err := filepath.Abs(path)
		if err != nil {
			program.Diagnostics = append(program.Diagnostics, diag.Diagnostic{
				Code:     "loader.path.invalid",
				Message:  err.Error(),
				Severity: diag.SeverityError,
			})
			continue
		}
		l.loadPath(abs, &program, visiting, nil)
	}
	return program
}

func (l Loader) loadPath(path string, program *Program, visiting map[string]bool, importSpan *diag.Span) {
	if program.ByPath[path] != nil {
		return
	}
	if visiting[path] {
		span := diag.Span{}
		if importSpan != nil {
			span = *importSpan
		}
		program.Diagnostics = append(program.Diagnostics, diag.Diagnostic{
			Code:     "loader.import.cycle",
			Message:  "import cycle detected",
			Severity: diag.SeverityError,
			Span:     span,
		})
		return
	}
	visiting[path] = true
	defer delete(visiting, path)

	data, err := os.ReadFile(path)
	if err != nil {
		span := diag.Span{}
		if importSpan != nil {
			span = *importSpan
		}
		program.Diagnostics = append(program.Diagnostics, diag.Diagnostic{
			Code:     "loader.read.failed",
			Message:  err.Error(),
			Severity: diag.SeverityError,
			Span:     span,
		})
		return
	}

	file, diags := parser.Parse(path, string(data))
	program.Diagnostics = append(program.Diagnostics, diags...)
	unit := &Unit{Path: path, ModuleName: l.moduleName(path), Source: string(data), File: file}
	program.ByPath[path] = unit
	program.Units = append(program.Units, unit)

	if file == nil {
		return
	}
	for _, decl := range file.Decls {
		imp, ok := decl.(*ast.ImportDecl)
		if !ok {
			continue
		}
		resolved := filepath.Join(l.Root, "lang", "l1", filepath.Join(imp.ModulePath...)+".l1")
		if _, err := os.Stat(resolved); err != nil {
			program.Diagnostics = append(program.Diagnostics, diag.Diagnostic{
				Code:     "loader.import.missing",
				Message:  fmt.Sprintf("import %s could not be resolved", strings.Join(imp.ModulePath, ".")),
				Severity: diag.SeverityError,
				Span:     imp.Span,
			})
			continue
		}
		span := imp.Span
		l.loadPath(resolved, program, visiting, &span)
	}
}

func (l Loader) moduleName(path string) string {
	rel, err := filepath.Rel(filepath.Join(l.Root, "lang", "l1"), path)
	if err != nil || strings.HasPrefix(rel, "..") {
		base := filepath.Base(path)
		return strings.TrimSuffix(base, filepath.Ext(base))
	}
	module := strings.TrimSuffix(rel, filepath.Ext(rel))
	return filepath.ToSlash(module)
}