// Package census scans a Go package's own source for the types that implement a
// node family, so a catalog of that family can be proved complete against the
// source rather than against a hand-kept list (ADR 0046).
//
// A family is identified by the method its interface declares — an Effect is a
// type with Resolve(*EffectContext), a Refinement a type with
// refine(*EffectContext, []LocalID) — which the source states plainly enough for
// go/ast to read. Two consumers share this scanner: the engine's catalog
// totality tests, which fail the build when a node has no census row, and
// `mage tool:census`, which reports the gaps while the census is half-filled.
package census

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
)

// Implementations returns every type declared in dir's non-test Go source that
// declares a method named method whose parameter types satisfy params, mapped to
// the file that declares the method. A pointer receiver is reported under the
// named type it points at.
//
// The parameter predicate is what separates two families whose methods share a
// name: Effect.Text takes none and LogEntry.Text takes a Namer. The file a type
// is found in is what lets a report group a family's gaps the way its source is
// grouped, one line per effect_<mechanic>.go.
func Implementations(
	dir, method string,
	params func(paramTypes []string) bool,
) (map[string]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}
	fset := token.NewFileSet()
	found := map[string]string{}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") ||
			strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", name, err)
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Name.Name != method {
				continue
			}
			if !params(paramTypes(fn.Type.Params)) {
				continue
			}
			if recv := receiverTypeName(fn.Recv); recv != "" {
				found[recv] = name
			}
		}
	}
	return found, nil
}

// Params returns a predicate matching a parameter list whose types are exactly
// the ones named, written as they appear in source ("*EffectContext",
// "[]LocalID"). Params() with no names matches a method that takes none.
func Params(want ...string) func([]string) bool {
	return func(got []string) bool {
		if len(got) != len(want) {
			return false
		}
		for i, t := range want {
			if got[i] != t {
				return false
			}
		}
		return true
	}
}

// paramTypes renders a parameter list as one type string per parameter, so a
// grouped list (a, b int) reports its type once per name the way a caller sees
// it.
func paramTypes(params *ast.FieldList) []string {
	if params == nil {
		return nil
	}
	var out []string
	for _, field := range params.List {
		rendered := types.ExprString(field.Type)
		for range max(len(field.Names), 1) {
			out = append(out, rendered)
		}
	}
	return out
}

// receiverTypeName returns the named type a method is declared on, unwrapping a
// pointer receiver and a generic type's parameter list.
func receiverTypeName(recv *ast.FieldList) string {
	if recv == nil || len(recv.List) != 1 {
		return ""
	}
	expr := recv.List[0].Type
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.IndexExpr:
		if id, ok := e.X.(*ast.Ident); ok {
			return id.Name
		}
	case *ast.IndexListExpr:
		if id, ok := e.X.(*ast.Ident); ok {
			return id.Name
		}
	}
	return ""
}
