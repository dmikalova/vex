// Package main contains the repo-local formatter for multiline keyed composite
// literals.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type violation struct {
	path string
	line int
	col  int
}

type rewrite struct {
	start int
	end   int
	text  string
}

func main() {
	fsys := flag.NewFlagSet("fmtmklv", flag.ContinueOnError)
	fsys.SetOutput(os.Stderr)
	fix := fsys.Bool("fix", false, "rewrite in place instead of checking")
	if err := fsys.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	paths := fsys.Args()
	if len(paths) == 0 {
		paths = []string{"./..."}
	}

	violations, err := run(paths, *fix)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(violations) > 0 {
		for _, v := range violations {
			fmt.Printf(
				"%s:%d:%d: multiline keyed composite literal required\n",
				v.path,
				v.line,
				v.col,
			)
		}
		if !*fix {
			os.Exit(1)
		}
	}
}

func run(paths []string, fix bool) ([]violation, error) {
	files, err := expandGoFiles(paths)
	if err != nil {
		return nil, err
	}

	var violations []violation
	for _, path := range files {
		src, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		out, changed, fileViolations, err := formatSource(path, src)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if len(fileViolations) > 0 {
			violations = append(violations, fileViolations...)
		}
		if fix && changed {
			if err := os.WriteFile(path, out, 0o644); err != nil {
				return nil, err
			}
		}
	}
	return violations, nil
}

func expandGoFiles(paths []string) ([]string, error) {
	seen := make(map[string]bool)
	var out []string
	for _, path := range paths {
		if path == "./..." || strings.Contains(path, "...") {
			pkgs, err := goListPackages(path)
			if err != nil {
				return nil, err
			}
			for _, pkgDir := range pkgs {
				if err := walkGoFiles(pkgDir, &out, seen); err != nil {
					return nil, err
				}
			}
			continue
		}
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			if err := walkGoFiles(path, &out, seen); err != nil {
				return nil, err
			}
			continue
		}
		if strings.HasSuffix(path, ".go") {
			if !seen[path] {
				seen[path] = true
				out = append(out, path)
			}
			continue
		}
		return nil, fmt.Errorf("unsupported path %q: expected .go file or package pattern", path)
	}
	return out, nil
}

func goListPackages(pattern string) ([]string, error) {
	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", pattern)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Fields(string(out))
	for i, s := range lines {
		lines[i] = strings.TrimSpace(s)
	}
	return lines, nil
}

func walkGoFiles(root string, out *[]string, seen map[string]bool) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		if !seen[path] {
			seen[path] = true
			*out = append(*out, path)
		}
		return nil
	})
}

func formatSource(path string, src []byte) ([]byte, bool, []violation, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, false, nil, err
	}

	var edits []rewrite
	var violations []violation
	ast.Inspect(file, func(node ast.Node) bool {
		lit, ok := node.(*ast.CompositeLit)
		if !ok || !shouldRewrite(lit) || !singleLine(fset, lit) {
			return true
		}
		rewriteText, err := rewriteLiteral(fset, lit, src)
		if err != nil {
			return false
		}
		if rewriteText == "" {
			return true
		}
		start := fset.Position(lit.Pos()).Offset
		end := fset.Position(lit.End()).Offset
		edits = append(edits, rewrite{
			start: start,
			end:   end,
			text:  rewriteText,
		})
		pos := fset.Position(lit.Pos())
		violations = append(violations, violation{
			path: path,
			line: pos.Line,
			col:  pos.Column,
		})
		return true
	})

	if len(edits) == 0 {
		return src, false, violations, nil
	}
	ordered := append([]rewrite(nil), edits...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].start < ordered[j].start })
	var out bytes.Buffer
	cursor := 0
	for _, edit := range ordered {
		if edit.start < cursor {
			continue
		}
		out.Write(src[cursor:edit.start])
		out.WriteString(edit.text)
		cursor = edit.end
	}
	out.Write(src[cursor:])
	result := out.Bytes()
	if formatted, err := format.Source(result); err == nil {
		result = formatted
	}
	return result, true, violations, nil
}

func shouldRewrite(lit *ast.CompositeLit) bool {
	if lit == nil || len(lit.Elts) < 2 || lit.Type == nil {
		return false
	}
	for _, elt := range lit.Elts {
		if _, ok := elt.(*ast.KeyValueExpr); !ok {
			return false
		}
	}
	return true
}

func singleLine(fset *token.FileSet, lit *ast.CompositeLit) bool {
	start := fset.Position(lit.Pos())
	end := fset.Position(lit.End())
	return start.Line == end.Line
}

func rewriteLiteral(fset *token.FileSet, lit *ast.CompositeLit, _ []byte) (string, error) {
	if lit.Type == nil {
		return "", nil
	}
	var typeBuf bytes.Buffer
	if err := printer.Fprint(&typeBuf, fset, lit.Type); err != nil {
		return "", err
	}
	parts := make([]string, 0, len(lit.Elts))
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			return "", nil
		}
		var kvBuf bytes.Buffer
		if err := printer.Fprint(&kvBuf, fset, kv); err != nil {
			return "", err
		}
		parts = append(parts, "\t"+strings.TrimSpace(kvBuf.String())+",")
	}
	return fmt.Sprintf("%s{\n%s\n}", typeBuf.String(), strings.Join(parts, "\n")), nil
}
