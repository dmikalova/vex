package cardtest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestNoCardNameStringLiterals guards the typed card-name reference convention: a
// card that names another card (or itself) must do so through a typed reference —
// another card's exported .Name, a const declared beside the lead card, or the
// {card} placeholder resolved through card.Target.GrantingCard — never a bare
// string literal. A literal desyncs silently when the referenced card is renamed,
// so this test fails the build if one appears in a card-name field.
//
// The guarded surface is any composite-literal field named Name or Host whose
// value is a string literal, plus any .Named("...") call. Two literal shapes are
// allowed because they are a card's own identity, not a reference to another
// card: card.Cluster{Name: "..."} (a deck-generation cluster names the pod it
// leads) and card.New's positional own-name argument (a call, not a Name: field,
// so it never matches).
func TestNoCardNameStringLiterals(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate this test file")
	}
	setsDir := filepath.Join(filepath.Dir(thisFile), "..", "sets")

	fset := token.NewFileSet()
	err := filepath.WalkDir(setsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.CompositeLit:
				if compositeTypeName(node.Type) == "card.Cluster" {
					return false // skip the whole cluster literal
				}
				for _, elt := range node.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					key, ok := kv.Key.(*ast.Ident)
					if !ok || (key.Name != "Name" && key.Name != "Host") {
						continue
					}
					if lit, ok := kv.Value.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						pos := fset.Position(lit.Pos())
						t.Errorf(
							"%s:%d: %s field carries a string literal %s — reference the card by its typed .Name, a const, or card.Target.GrantingCard",
							filepath.Base(pos.Filename),
							pos.Line,
							key.Name,
							lit.Value,
						)
					}
				}
			case *ast.CallExpr:
				sel, ok := node.Fun.(*ast.SelectorExpr)
				if !ok || sel.Sel.Name != "Named" || len(node.Args) != 1 {
					return true
				}
				if lit, ok := node.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
					pos := fset.Position(lit.Pos())
					t.Errorf(
						"%s:%d: .Named(%s) carries a string literal — reference the card by its typed .Name or a const",
						filepath.Base(pos.Filename),
						pos.Line,
						lit.Value,
					)
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("walking %s: %v", setsDir, err)
	}
}

// compositeTypeName renders a composite literal's type as pkg.Type, or "" when it
// is not a package-qualified selector (an anonymous or local type).
func compositeTypeName(expr ast.Expr) string {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return ""
	}
	return pkg.Name + "." + sel.Sel.Name
}
