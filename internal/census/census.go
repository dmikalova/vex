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
//
// The scan reads the source at four heights, one per shape a catalogued thing
// takes: Implementations, Builders and Fields find a family's members, Constants
// finds an enum's, and Interfaces finds the families themselves.
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
	return scan(dir, func(file string, f *ast.File, found map[string]string) {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Name.Name != method {
				continue
			}
			if !params(paramTypes(fn.Type.Params)) {
				continue
			}
			if recv := receiverTypeName(fn.Recv); recv != "" {
				found[recv] = file
			}
		}
	})
}

// Builders returns every method in dir's non-test Go source declared on the
// named receiver type whose only result is the named type, mapped to the file
// declaring it. Unexported methods are left out: a builder a card cannot write is
// not a member of the family a card writes.
//
// It is the shape half of the scan, for a family identified by what its members
// look like rather than by a method they all declare. Target is such a family: it
// is a flag struct rather than an interface, so its members are its filter
// builders, `func (t Target) X(...) Target`.
func Builders(dir, recv, result string) (map[string]string, error) {
	return scan(dir, func(file string, f *ast.File, found map[string]string) {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || !fn.Name.IsExported() {
				continue
			}
			if receiverTypeName(fn.Recv) != recv || !returnsOnly(fn.Type.Results, result) {
				continue
			}
			found[fn.Name.Name] = file
		}
	})
}

// Fields returns every exported field of the named struct type declared in dir's
// non-test Go source, mapped to the file declaring the struct. Unexported fields
// are left out: a field a card cannot write is not a member of the family a card
// writes. An embedded field is reported under the type it embeds, which is the
// name a literal writes it by.
//
// It is the second shape half of the scan, for a family whose members are the
// axes of one comparable value rather than methods. Filter is such a family: a
// card names its narrowing in a struct literal (ADR 0005), so an added axis is a
// field, and this is what makes an uncatalogued one a red build.
func Fields(dir, typeName string) (map[string]string, error) {
	return scan(dir, func(file string, f *ast.File, found map[string]string) {
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Name.Name != typeName {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				for _, name := range structFieldNames(st) {
					found[name] = file
				}
			}
		}
	})
}

// structFieldNames returns the exported field names a struct declares, naming an
// embedded field by the type it embeds.
func structFieldNames(st *ast.StructType) []string {
	var out []string
	for _, field := range st.Fields.List {
		if len(field.Names) == 0 {
			if name := embeddedName(field.Type); name != "" {
				out = append(out, name)
			}
			continue
		}
		for _, name := range field.Names {
			if name.IsExported() {
				out = append(out, name.Name)
			}
		}
	}
	return out
}

// embeddedName returns the exported name an embedded field is written by,
// unwrapping a pointer and a package qualifier, and none when the embedded type
// is unexported.
func embeddedName(expr ast.Expr) string {
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		expr = sel.Sel
	}
	ident, ok := expr.(*ast.Ident)
	if !ok || !ident.IsExported() {
		return ""
	}
	return ident.Name
}

// returnsOnly reports whether a result list is exactly one value of the named
// type, so a method returning the type alongside an ok flag is not a builder.
func returnsOnly(results *ast.FieldList, typeName string) bool {
	rendered := paramTypes(results)
	return len(rendered) == 1 && rendered[0] == typeName
}

// Interfaces returns every interface type declared in dir's non-test Go source,
// mapped to the file declaring it. Only a named declaration is reported: an
// anonymous interface written inline in a signature declares nothing a catalog
// could name.
//
// It closes the census one level above the families. A node family is an
// interface, so enumerating the package's interfaces is what lets a
// classification insist that each one is either a catalogued family or a
// recorded non-family — otherwise a new strategy axis and its implementations
// would sit outside every catalog with the build green.
func Interfaces(dir string) (map[string]string, error) {
	return scan(dir, func(file string, f *ast.File, found map[string]string) {
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if _, ok := ts.Type.(*ast.InterfaceType); ok {
					found[ts.Name.Name] = file
				}
			}
		}
	})
}

// Constants returns every constant declared with the named type in dir's
// non-test Go source, mapped to the file declaring it. It is the enum half of the
// census: an enum's members are constants rather than types, so an enumerating
// function like Keywords() is proved complete against this rather than against
// the author's memory.
//
// It reads a const block the way the compiler does. A spec that names neither a
// type nor a value repeats the previous spec's type, which is what makes an iota
// run of one-word lines belong to the type its first line named; a spec that
// gives a value of its own does not, so `NumCounterKinds = int(iota)` closing a
// CounterKind block is an int and not reported.
func Constants(dir, typeName string) (map[string]string, error) {
	return scan(dir, func(file string, f *ast.File, found map[string]string) {
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.CONST {
				continue
			}
			carried := ""
			for _, spec := range gen.Specs {
				value, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				carried = specType(value, carried)
				if carried != typeName {
					continue
				}
				for _, name := range value.Names {
					if name.Name != "_" {
						found[name.Name] = file
					}
				}
			}
		}
	})
}

// specType returns the type a const spec is declared with: its own when it names
// one, the type carried down the block when it names neither a type nor a value,
// and none when it gives a value of its own.
func specType(spec *ast.ValueSpec, carried string) string {
	if spec.Type != nil {
		if ident, ok := spec.Type.(*ast.Ident); ok {
			return ident.Name
		}
		return ""
	}
	if len(spec.Values) > 0 {
		return ""
	}
	return carried
}

// scan parses each non-test Go file in dir and lets collect record what it finds,
// keyed by name and mapped to the file it was found in.
func scan(
	dir string,
	collect func(file string, f *ast.File, found map[string]string),
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
		collect(name, f, found)
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
