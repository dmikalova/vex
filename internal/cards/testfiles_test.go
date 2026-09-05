package cards

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// setsDir is the path, relative to this package, to the per-set card packages.
const setsDir = "sets"

// setPackage is one set's source split into its card implementation files and its
// test files, both keyed by the shared snake_case base name (a card in
// "succubus.go" is tested in "succubus_test.go").
type setPackage struct {
	name string
	// cardFiles maps each implementation's base name (filename minus ".go") to that
	// filename. An implementation is a buildable, non-test file that declares an
	// exported `var X = card.New(...)`. Build-excluded stubs (//go:build todo) and
	// the generated 0set.go (reprints only, no card vars) are not implementations.
	cardFiles map[string]string
	// testFiles maps each test's base name (filename minus "_test.go") to that
	// filename.
	testFiles map[string]string
}

// TestEveryCardHasATest enforces the one-file-per-card convention: each card
// implementation "foo.go" must be tested in a sibling "foo_test.go". Grouping
// several cards' tests into one file by shared mechanic is disallowed — the
// mechanic itself belongs to the engine and is covered by engine tests; a card's
// test file proves that one card's wiring.
func TestEveryCardHasATest(t *testing.T) {
	for _, pkg := range setPackages(t) {
		var missing []string
		for base, file := range pkg.cardFiles {
			if _, ok := pkg.testFiles[base]; !ok {
				missing = append(missing, file)
			}
		}
		sort.Strings(missing)
		for _, file := range missing {
			t.Errorf("%s: %s has no matching %s_test.go", pkg.name, file, stem(file))
		}
	}
}

// TestNoOrphanedTestFiles enforces the converse: each "foo_test.go" must test a
// card implemented in a sibling "foo.go". An orphaned test file — one whose card
// was renamed or removed, or a mechanic-grouped file that tests several cards at
// once — fails here.
func TestNoOrphanedTestFiles(t *testing.T) {
	for _, pkg := range setPackages(t) {
		var orphans []string
		for base, file := range pkg.testFiles {
			if _, ok := pkg.cardFiles[base]; !ok {
				orphans = append(orphans, file)
			}
		}
		sort.Strings(orphans)
		for _, file := range orphans {
			t.Errorf("%s: %s has no matching card implementation %s.go", pkg.name, file, stem(file))
		}
	}
}

// stem returns a filename's shared base name, dropping its ".go" or "_test.go"
// suffix so a card file and its test file map to the same key.
func stem(file string) string {
	return strings.TrimSuffix(strings.TrimSuffix(file, ".go"), "_test")
}

// setPackages parses every set package directory and splits each into its card
// implementation files and its test files. It uses go/build so the same build
// constraints the compiler applies decide which files count: //go:build todo
// stubs land in IgnoredGoFiles and are not treated as implementations.
func setPackages(t *testing.T) []setPackage {
	t.Helper()

	entries, err := filepath.Glob(filepath.Join(setsDir, "*"))
	if err != nil {
		t.Fatalf("globbing set packages: %v", err)
	}

	var pkgs []setPackage
	fset := token.NewFileSet()
	for _, dir := range entries {
		bp, err := build.ImportDir(dir, 0)
		if err != nil {
			t.Fatalf("importing %s: %v", dir, err)
		}

		pkg := setPackage{
			name:      filepath.Base(dir),
			cardFiles: make(map[string]string),
			testFiles: make(map[string]string),
		}
		for _, file := range bp.GoFiles {
			if declaresCard(parseFile(t, fset, filepath.Join(dir, file))) {
				pkg.cardFiles[stem(file)] = file
			}
		}
		for _, file := range append(bp.TestGoFiles, bp.XTestGoFiles...) {
			pkg.testFiles[stem(file)] = file
		}

		// Guard against a parsing regression silently finding nothing, which would
		// turn both checks into a false green: every set has cards and tests.
		if len(pkg.cardFiles) == 0 {
			t.Fatalf("%s: found no card implementations; parsing is broken", pkg.name)
		}
		if len(pkg.testFiles) == 0 {
			t.Fatalf("%s: found no test files; parsing is broken", pkg.name)
		}
		pkgs = append(pkgs, pkg)
	}

	if len(pkgs) == 0 {
		t.Fatalf("no set packages found under %s", setsDir)
	}
	return pkgs
}

// parseFile parses one Go source file or fails the test.
func parseFile(t *testing.T, fset *token.FileSet, path string) *ast.File {
	t.Helper()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parsing %s: %v", path, err)
	}
	return f
}

// declaresCard reports whether f declares an exported top-level var assigned from
// card.New — i.e. whether the file implements a card.
func declaresCard(f *ast.File) bool {
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.VAR {
			continue
		}
		for _, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || len(vs.Values) != 1 {
				continue
			}
			if vs.Names[0].IsExported() && isCardNewCall(vs.Values[0]) {
				return true
			}
		}
	}
	return false
}

// isCardNewCall reports whether expr is a call to card.New.
func isCardNewCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "New" {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && pkg.Name == "card"
}
