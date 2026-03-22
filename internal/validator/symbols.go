package validator

import "github.com/yy/len/internal/ast"

type namespace string

const (
	namespaceType    namespace = "type"
	namespaceRel     namespace = "rel"
	namespaceFn      namespace = "fn"
	namespaceConst   namespace = "const"
	namespaceTrait   namespace = "trait"
	namespaceKeyword namespace = "keyword"
	namespaceSymbol  namespace = "symbol"
	namespaceSpec    namespace = "spec"
)

type symbol struct {
	Name  string
	Arity int
	Decl  ast.Decl
	Space namespace
}

type symbolTable struct {
	bySpace map[namespace]map[string]symbol
}

func newSymbolTable() symbolTable {
	return symbolTable{bySpace: map[namespace]map[string]symbol{}}
}

func (t *symbolTable) put(space namespace, name string, value symbol) (symbol, bool) {
	if t.bySpace[space] == nil {
		t.bySpace[space] = map[string]symbol{}
	}
	existing, ok := t.bySpace[space][name]
	if ok {
		return existing, false
	}
	t.bySpace[space][name] = value
	return value, true
}

func (t *symbolTable) get(space namespace, name string) (symbol, bool) {
	spaceMap := t.bySpace[space]
	if spaceMap == nil {
		return symbol{}, false
	}
	value, ok := spaceMap[name]
	return value, ok
}
