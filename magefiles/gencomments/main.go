// Command gencomments rewrites the doc comment above every card definition from
// the definition itself. It is the source of `mage generateComments` (and half of
// `mage gen`).
//
// A card's comment is its printed card: the name, the labeled House/Type/Rarity/
// stats block, and the rules text produced by the effect AST's Text() methods
// (engine.RenderCardText). Generating it means the comment can never drift from
// behavior — the way to change printed text is to change the effect, not the
// comment.
//
// The card definitions come from the live registry (importing package cards
// enrolls every set), keyed by the name literal in each file's card.New call, so
// no file needs to declare which card it holds. Files excluded from the build
// (the `//go:build todo` stubs) register nothing and are left alone.
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
)

// setsRoot is the tree of card files whose comments are generated.
const setsRoot = "internal/cards/sets"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gencomments:", err)
		os.Exit(1)
	}
}

func run() error {
	defs := definitionsByName()
	files, err := cardFiles(setsRoot)
	if err != nil {
		return err
	}
	changed := 0
	for _, path := range files {
		ok, err := rewriteFile(path, defs)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if ok {
			changed++
		}
	}
	fmt.Printf("gencomments: %d card files scanned, %d comments rewritten\n", len(files), changed)
	return nil
}

// definitionsByName indexes the live card registry by printed name.
func definitionsByName() map[string]engine.CardDefinition {
	out := make(map[string]engine.CardDefinition)
	for _, d := range cards.All() {
		out[d.Name] = d
	}
	return out
}

// cardFiles lists the non-test Go files under root, in a stable order.
func cardFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		switch {
		case err != nil:
			return err
		case d.IsDir(), !strings.HasSuffix(path, ".go"), strings.HasSuffix(path, "_test.go"):
			return nil
		}
		files = append(files, path)
		return nil
	})
	sort.Strings(files)
	return files, err
}

// rewriteFile replaces the doc comment of every `var X = card.New("Name", …)` in
// one file, reporting whether anything changed. Edits are applied to the raw
// source bytes from the end backwards so earlier offsets stay valid, which keeps
// the rest of the file — including hand-written formatting — untouched.
func rewriteFile(path string, defs map[string]engine.CardDefinition) (bool, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return false, err
	}
	base := fset.File(f.Pos()).Base()

	type edit struct {
		start, end int // byte range of the existing comment (empty range = insert)
		text       string
	}
	var edits []edit
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		name, ok := cardNameOf(gd)
		if !ok {
			continue
		}
		def, ok := defs[name]
		if !ok {
			// The card is not in the registry: a build-excluded stub, or a name that
			// does not match its card.New literal. Either way there is nothing to
			// render from, so leave the file as the author wrote it.
			continue
		}
		start, end := int(gd.Pos())-base, int(gd.Pos())-base
		if gd.Doc != nil {
			start = int(gd.Doc.Pos()) - base
		}
		edits = append(edits, edit{start, end, renderComment(&def)})
	}
	out := src
	for i := len(edits) - 1; i >= 0; i-- {
		e := edits[i]
		out = append(append(append([]byte{}, out[:e.start]...), e.text...), out[e.end:]...)
	}
	if bytes.Equal(out, src) {
		return false, nil
	}
	return true, os.WriteFile(path, out, 0o644)
}

// cardNameOf returns the printed name a `var X = card.New("Name", …)` declaration
// builds, or false if the declaration is not a card.
func cardNameOf(gd *ast.GenDecl) (string, bool) {
	if len(gd.Specs) != 1 {
		return "", false
	}
	vs, ok := gd.Specs[0].(*ast.ValueSpec)
	if !ok || len(vs.Values) != 1 {
		return "", false
	}
	call, ok := vs.Values[0].(*ast.CallExpr)
	if !ok || len(call.Args) == 0 {
		return "", false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "New" {
		return "", false
	}
	if pkg, ok := sel.X.(*ast.Ident); !ok || pkg.Name != "card" {
		return "", false
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	name, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}
	return name, true
}

// renderComment builds a card's whole doc comment, ending in the newline that
// butts it against the var declaration. The detail block is tab-indented so godoc
// renders it preformatted; gofmt requires the blank `//` line after the title.
func renderComment(def *engine.CardDefinition) string {
	var b strings.Builder
	b.WriteString("// " + def.Name + "\n//\n")
	for _, line := range strings.Split(engine.RenderCardText(def), "\n") {
		if line == "" {
			b.WriteString("//\n")
			continue
		}
		b.WriteString("//\t" + line + "\n")
	}
	return b.String()
}
