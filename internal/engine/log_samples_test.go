package engine

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
	"testing"
)

// TestLogEntrySamplesTotality binds the sample catalog to the set of LogEntry
// variants the source actually declares, so a new variant fails the build until
// it is catalogued (ADR 0046). It reads the package's own source rather than a
// hand-kept list: a variant is any type with a Text(Namer) method, which is the
// LogEntry interface. Record — the frame-plus-entry wrapper — implements Text the
// same way but is not itself a narrated variant, so it is excluded.
func TestLogEntrySamplesTotality(t *testing.T) {
	variants := logEntryVariantTypes(t)
	sampled := map[string]bool{}
	for _, e := range LogEntrySamples() {
		sampled[reflect.TypeOf(e).Name()] = true
	}
	for name := range variants {
		if !sampled[name] {
			t.Errorf("log entry %s has no sample in LogEntrySamples; add one (ADR 0046)", name)
		}
	}
	for name := range sampled {
		if !variants[name] {
			t.Errorf("LogEntrySamples has %s, which is not a LogEntry variant", name)
		}
	}
}

// logEntryVariantTypes parses the engine package's non-test source and returns
// the name of every type that implements LogEntry — a type with a Text(Namer)
// method — with Record excluded.
func logEntryVariantTypes(t *testing.T) map[string]bool {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("reading engine package directory: %v", err)
	}
	fset := token.NewFileSet()
	types := map[string]bool{}
	for _, entry := range entries {
		source := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(source, ".go") ||
			strings.HasSuffix(source, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, source, nil, 0)
		if err != nil {
			t.Fatalf("parsing %s: %v", source, err)
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Name.Name != "Text" {
				continue
			}
			if !hasSingleNamerParam(fn.Type.Params) {
				continue
			}
			if name := receiverTypeName(fn.Recv); name != "" && name != "Record" {
				types[name] = true
			}
		}
	}
	if len(types) == 0 {
		t.Fatal("found no LogEntry variants; AST parsing is broken")
	}
	return types
}

// hasSingleNamerParam reports whether a parameter list is exactly one Namer,
// which is the LogEntry.Text signature and distinguishes it from Effect.Text,
// which takes none.
func hasSingleNamerParam(params *ast.FieldList) bool {
	if params == nil || len(params.List) != 1 {
		return false
	}
	field := params.List[0]
	if len(field.Names) > 1 {
		return false
	}
	id, ok := field.Type.(*ast.Ident)
	return ok && id.Name == "Namer"
}

// receiverTypeName returns the named type a method is declared on, unwrapping a
// pointer receiver.
func receiverTypeName(recv *ast.FieldList) string {
	if recv == nil || len(recv.List) != 1 {
		return ""
	}
	switch e := recv.List[0].Type.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		if id, ok := e.X.(*ast.Ident); ok {
			return id.Name
		}
	}
	return ""
}
