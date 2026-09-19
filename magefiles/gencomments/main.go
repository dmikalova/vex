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
//
// The same comment is written above each card's `func Test<Card>` in its
// `_test.go` (the card-authoring guide says that block is generated, not
// hand-written), matched by mapping the test name Test<X> back to the card var X
// declared in the same set package.
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
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/dmikalova/vactrol/internal/card"
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
	templates := templateNames()
	srcFiles, err := cardFiles(setsRoot)
	if err != nil {
		return err
	}
	testFiles, err := testCardFiles(setsRoot)
	if err != nil {
		return err
	}
	names, err := varNamesByDir(srcFiles)
	if err != nil {
		return err
	}
	changed := 0
	for _, path := range srcFiles {
		ok, err := rewriteFile(path, defs, templates)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if ok {
			changed++
		}
	}
	for _, path := range testFiles {
		ok, err := rewriteTestFile(path, defs, names[filepath.Dir(path)], templates)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if ok {
			changed++
		}
	}
	fmt.Printf("gencomments: %d card files scanned, %d comments rewritten\n",
		len(srcFiles)+len(testFiles), changed)
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

// templateNames is the set of printed names whose card is a generative template
// (it carries a Materializer, see card.Template): the card in the file is a
// placeholder face, and the concrete card is chosen per deck at generation. The
// comment box flags this so a reader does not mistake the sparse face for a
// finished card.
func templateNames() map[string]bool {
	out := map[string]bool{}
	for _, rc := range card.Cards() {
		if rc.Materializer != nil {
			out[rc.Def.Name] = true
		}
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

// testCardFiles lists the _test.go files under root, in a stable order — the card
// tests whose doc comment above func Test<Card> mirrors the card box.
func testCardFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		switch {
		case err != nil:
			return err
		case d.IsDir(), !strings.HasSuffix(path, "_test.go"):
			return nil
		}
		files = append(files, path)
		return nil
	})
	sort.Strings(files)
	return files, err
}

// varNamesByDir maps each set directory to its card vars — the identifier a
// `var X = card.New("Name", …)` declares, mapped to the printed Name. It links a
// test (func Test<X>) back to the card it exercises so the test's doc comment is
// generated from the same definition as the card's own file.
func varNamesByDir(files []string) (map[string]map[string]string, error) {
	out := map[string]map[string]string{}
	for _, path := range files {
		src, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, src, 0)
		if err != nil {
			return nil, err
		}
		dir := filepath.Dir(path)
		wrappers := wrapperTemplates(f)
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			ident, name, ok := cardVarOf(gd, wrappers)
			if !ok {
				continue
			}
			if out[dir] == nil {
				out[dir] = map[string]string{}
			}
			out[dir][ident] = name
		}
	}
	return out, nil
}

// edit is a byte-range replacement in a file's source: [start,end) becomes text.
// An empty range (start==end) is an insertion point.
type edit struct {
	start, end int
	text       string
}

// applyEdits rewrites src by applying edits from the end backwards, so earlier
// offsets stay valid and the rest of the file — including hand-written formatting
// — is untouched.
func applyEdits(src []byte, edits []edit) []byte {
	out := src
	for _, e := range slices.Backward(edits) {

		out = append(append(append([]byte{}, out[:e.start]...), e.text...), out[e.end:]...)
	}
	return out
}

// rewriteFile replaces the doc comment of every `var X = card.New("Name", …)` in
// one file, reporting whether anything changed. Edits are applied to the raw
// source bytes from the end backwards so earlier offsets stay valid, which keeps
// the rest of the file — including hand-written formatting — untouched.
func rewriteFile(
	path string,
	defs map[string]engine.CardDefinition,
	templates map[string]bool,
) (bool, error) {
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
	wrappers := wrapperTemplates(f)

	var edits []edit
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		name, ok := cardNameOf(gd, wrappers)
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
		edits = append(edits, edit{start, end, renderComment(&def, templates[name])})
	}
	out := applyEdits(src, edits)
	if bytes.Equal(out, src) {
		return false, nil
	}
	return true, os.WriteFile(path, out, 0o644)
}

// rewriteTestFile replaces the doc comment of every `func Test<Card>(t *testing.T)`
// in one test file with the card's card-box comment — the same block the card's own
// source file carries. The card is found by mapping the test name Test<X> to the
// card var X declared in the same package (varToName). Test funcs whose name maps
// to no card, or a card not in the registry, are left alone.
func rewriteTestFile(
	path string,
	defs map[string]engine.CardDefinition,
	varToName map[string]string,
	templates map[string]bool,
) (bool, error) {
	if len(varToName) == 0 {
		return false, nil
	}
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

	var edits []edit
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv != nil || !strings.HasPrefix(fd.Name.Name, "Test") {
			continue
		}
		name, ok := varToName[strings.TrimPrefix(fd.Name.Name, "Test")]
		if !ok {
			continue
		}
		def, ok := defs[name]
		if !ok {
			continue
		}
		start, end := int(fd.Pos())-base, int(fd.Pos())-base
		if fd.Doc != nil {
			start = int(fd.Doc.Pos()) - base
		}
		edits = append(edits, edit{start, end, renderComment(&def, templates[name])})
	}
	out := applyEdits(src, edits)
	if bytes.Equal(out, src) {
		return false, nil
	}
	return true, os.WriteFile(path, out, 0o644)
}

// cardVarOf returns the identifier and printed name a card-declaring var builds,
// or ok=false if the declaration is not one. It recognizes any single-var
// initializer whose value is a call to card.New or a set-local family wrapper.
// The name comes from a string-literal first argument (`card.New("Name", …)` or
// `master("Master of 1", …)`); when the first argument is not a literal, the name
// is resolved through the wrapper (`master(1, …)` → "Master of 1", see callName).
// The callee is not checked — a name absent from the built registry is filtered
// downstream (see rewriteFile), so only real cards are documented.
func cardVarOf(gd *ast.GenDecl, wrappers map[string]nameTemplate) (ident, name string, ok bool) {
	if gd.Tok != token.VAR || len(gd.Specs) != 1 {
		return "", "", false
	}
	vs, ok := gd.Specs[0].(*ast.ValueSpec)
	if !ok || len(vs.Names) != 1 || len(vs.Values) != 1 || !vs.Names[0].IsExported() {
		return "", "", false
	}
	call, ok := vs.Values[0].(*ast.CallExpr)
	if !ok || len(call.Args) == 0 {
		return "", "", false
	}
	name, ok = callName(call, wrappers)
	if !ok {
		return "", "", false
	}
	return vs.Names[0].Name, name, true
}

// callName returns the printed card name a card-building call yields: its
// string-literal first argument when it has one, or — when the first argument is
// not a literal — the name a set-local wrapper builds from the call's literal
// arguments via fmt.Sprintf (master(1, …) → "Master of 1").
func callName(call *ast.CallExpr, wrappers map[string]nameTemplate) (string, bool) {
	if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
		s, err := strconv.Unquote(lit.Value)
		return s, err == nil
	}
	id, ok := call.Fun.(*ast.Ident)
	if !ok {
		return "", false
	}
	tmpl, ok := wrappers[id.Name]
	if !ok {
		return "", false
	}
	return tmpl.eval(call.Args)
}

// nameTemplate describes how a set-local wrapper builds its card.New name from
// the wrapper's parameters: a fmt.Sprintf format and, in verb order, the wrapper
// parameter index feeding each verb.
type nameTemplate struct {
	format string
	args   []int
}

// eval renders the name from a wrapper call's arguments, or false if a needed
// argument is missing or not an int/string literal.
func (t nameTemplate) eval(callArgs []ast.Expr) (string, bool) {
	vals := make([]any, len(t.args))
	for i, idx := range t.args {
		if idx >= len(callArgs) {
			return "", false
		}
		v, ok := literalValue(callArgs[idx])
		if !ok {
			return "", false
		}
		vals[i] = v
	}
	return fmt.Sprintf(t.format, vals...), true
}

// wrapperTemplates finds set-local wrapper funcs in f whose card.New name is a
// fmt.Sprintf over the wrapper's parameters, keyed by wrapper name, so a call
// like master(1, …) resolves to a printed name without a string-literal argument.
func wrapperTemplates(f *ast.File) map[string]nameTemplate {
	out := map[string]nameTemplate{}
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv != nil || fd.Body == nil || fd.Type.Params == nil {
			continue
		}
		arg, ok := cardNewNameArg(fd.Body)
		if !ok {
			continue
		}
		if tmpl, ok := sprintfTemplate(arg, paramIndex(fd.Type)); ok {
			out[fd.Name.Name] = tmpl
		}
	}
	return out
}

// cardNewNameArg returns the first argument of the card-building New call in a
// wrapper body — the expression that builds the printed name. The call is either
// card.New (the facade) or set.New (a set package's registrar).
func cardNewNameArg(body *ast.BlockStmt) (ast.Expr, bool) {
	var arg ast.Expr
	ast.Inspect(body, func(n ast.Node) bool {
		if arg != nil {
			return false
		}
		if call, ok := n.(*ast.CallExpr); ok && len(call.Args) > 0 &&
			(isSelector(call.Fun, "card", "New") || isSelector(call.Fun, "set", "New")) {
			arg = call.Args[0]
			return false
		}
		return true
	})
	return arg, arg != nil
}

// sprintfTemplate reads a `fmt.Sprintf(format, params…)` name expression into a
// nameTemplate, mapping each verb argument back to the wrapper parameter feeding
// it. It fails on any argument that is not one of the wrapper's parameters.
func sprintfTemplate(expr ast.Expr, params map[string]int) (nameTemplate, bool) {
	call, ok := expr.(*ast.CallExpr)
	if !ok || !isSelector(call.Fun, "fmt", "Sprintf") || len(call.Args) < 1 {
		return nameTemplate{}, false
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return nameTemplate{}, false
	}
	format, err := strconv.Unquote(lit.Value)
	if err != nil {
		return nameTemplate{}, false
	}
	args := make([]int, 0, len(call.Args)-1)
	for _, a := range call.Args[1:] {
		id, ok := a.(*ast.Ident)
		if !ok {
			return nameTemplate{}, false
		}
		idx, ok := params[id.Name]
		if !ok {
			return nameTemplate{}, false
		}
		args = append(args, idx)
	}
	return nameTemplate{format: format, args: args}, true
}

// paramIndex maps each named parameter of ft to its positional index.
func paramIndex(ft *ast.FuncType) map[string]int {
	out := map[string]int{}
	i := 0
	for _, field := range ft.Params.List {
		if len(field.Names) == 0 {
			i++
			continue
		}
		for _, n := range field.Names {
			out[n.Name] = i
			i++
		}
	}
	return out
}

// literalValue returns the Go value of an int or string basic literal.
func literalValue(expr ast.Expr) (any, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return nil, false
	}
	switch lit.Kind {
	case token.INT:
		n, err := strconv.Atoi(lit.Value)
		return n, err == nil
	case token.STRING:
		s, err := strconv.Unquote(lit.Value)
		return s, err == nil
	}
	return nil, false
}

// isSelector reports whether e is the selector expression pkg.name.
func isSelector(e ast.Expr, pkg, name string) bool {
	sel, ok := e.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != name {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	return ok && id.Name == pkg
}

// cardNameOf returns the printed name a `var X = card.New("Name", …)` (or a
// set-local wrapper) declaration builds, or false if the declaration is not a card.
func cardNameOf(gd *ast.GenDecl, wrappers map[string]nameTemplate) (string, bool) {
	_, name, ok := cardVarOf(gd, wrappers)
	return name, ok
}

// renderComment builds a card's whole doc comment, ending in the newline that
// butts it against the var declaration. The detail block is tab-indented so godoc
// renders it preformatted; gofmt requires the blank `//` line after the title.
func renderComment(def *engine.CardDefinition, isTemplate bool) string {
	var b strings.Builder
	b.WriteString("// " + def.Name + "\n//\n")
	for line := range strings.SplitSeq(engine.RenderCardText(def), "\n") {
		if line == "" {
			b.WriteString("//\n")
			continue
		}
		b.WriteString("//\t" + line + "\n")
	}
	if isTemplate {
		b.WriteString(
			"//\n//\tTemplate: its concrete card is materialized per deck at generation.\n",
		)
	}
	return b.String()
}
