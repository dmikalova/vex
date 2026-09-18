package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// facadeDir is the authoring surface whose exported names this report inventories.
const facadeDir = "internal/card"

// setsDir holds the implemented card definitions that consume the facade.
const setsDir = "internal/cards/sets"

// node is one exported name on the card facade and how widely cards use it.
type node struct {
	Name     string
	Category string
	Cards    int // card definition files naming it
}

// nodeUsage reports every exported name on the card facade grouped by the
// category its declaration block documents, with how many card definitions use
// it. It answers "is this node carrying its weight, and what else lives on its
// axis" without a throwaway grep. Low-usage names are not defects on their own —
// half the card pool is unimplemented — but a node no card reaches for is the
// place to check that it decomposes into reusable atoms rather than hard-coding
// one card (see the composability rule in docs/style-guide.md).
func nodeUsage(args []string) error {
	maxUses := -1
	category := ""
	for _, a := range args {
		switch {
		case strings.HasPrefix(a, "-max="):
			n, err := strconv.Atoi(strings.TrimPrefix(a, "-max="))
			if err != nil {
				return fmt.Errorf("cardlookup node-usage: bad -max: %w", err)
			}
			maxUses = n
		case strings.HasPrefix(a, "-category="):
			category = strings.ToLower(strings.TrimPrefix(a, "-category="))
		default:
			return fmt.Errorf(
				"usage: cardlookup node-usage [-max=<n>] [-category=<substring>]")
		}
	}

	nodes, err := facadeNodes()
	if err != nil {
		return err
	}
	if err := countNodeUses(nodes); err != nil {
		return err
	}
	printNodeUsage(nodes, maxUses, category)
	return nil
}

// facadeNodes parses internal/card and returns every exported name it declares,
// tagged with the category its declaration block's doc comment names ("Æmber
// effects", "Damage and combat", ...). A block with no doc comment falls back to
// its file stem, so a name is never uncategorized.
func facadeNodes() ([]*node, error) {
	entries, err := os.ReadDir(facadeDir)
	if err != nil {
		return nil, fmt.Errorf("cardlookup node-usage: %w", err)
	}
	fset := token.NewFileSet()
	var nodes []*node
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := filepath.Join(facadeDir, name)
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("cardlookup node-usage: %w", err)
		}
		stem := strings.TrimSuffix(name, ".go")
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok == token.IMPORT {
				continue
			}
			cat := blockCategory(gen, stem)
			for _, spec := range gen.Specs {
				for _, id := range specNames(spec) {
					if !id.IsExported() {
						continue
					}
					nodes = append(nodes, &node{Name: id.Name, Category: cat})
				}
			}
		}
	}
	return nodes, nil
}

// specNames returns the identifiers a type or value spec declares.
func specNames(spec ast.Spec) []*ast.Ident {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		return []*ast.Ident{s.Name}
	case *ast.ValueSpec:
		return s.Names
	}
	return nil
}

// blockCategory reads a declaration block's heading — the first sentence of its
// doc comment, e.g. "// Damage and combat." — and falls back to the file stem.
// Only a parenthesized block has a heading; a lone declaration's doc comment
// describes that one name, not a group, so it groups under its file instead.
func blockCategory(gen *ast.GenDecl, stem string) string {
	if gen.Doc == nil || !gen.Lparen.IsValid() {
		return stem
	}
	head := firstSentence(strings.TrimSpace(gen.Doc.Text()))
	if head == "" {
		return stem
	}
	if len(head) > categoryWidth {
		head = head[:categoryWidth] + "…"
	}
	return head
}

// categoryWidth caps a category heading so one long doc comment cannot widen the
// report past a terminal line.
const categoryWidth = 56

// firstSentence returns text up to the first sentence-ending period — one
// followed by whitespace or end of text — or up to the first line break.
func firstSentence(text string) string {
	if i := strings.IndexByte(text, '\n'); i >= 0 {
		text = text[:i]
	}
	for i, r := range text {
		if r != '.' {
			continue
		}
		if i+1 == len(text) || text[i+1] == ' ' {
			return strings.TrimSpace(text[:i])
		}
	}
	return strings.TrimSpace(text)
}

// countNodeUses walks the implemented card definitions and counts, per node, the
// card files that name it. Tests are skipped: a node used only by its own test is
// unused by the card pool, which is the fact this report is for.
func countNodeUses(nodes []*node) error {
	patterns := make(map[*node]*regexp.Regexp, len(nodes))
	for _, n := range nodes {
		patterns[n] = regexp.MustCompile(`\bcard\.` + regexp.QuoteMeta(n.Name) + `\b`)
	}
	err := filepath.WalkDir(setsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		body := string(data)
		for _, n := range nodes {
			if patterns[n].MatchString(body) {
				n.Cards++
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("cardlookup node-usage: %w", err)
	}
	return nil
}

// printNodeUsage prints the nodes grouped by category, rarest first inside each
// group so the names to scrutinize lead, then a summary of the whole facade.
func printNodeUsage(nodes []*node, maxUses int, category string) {
	byCategory := map[string][]*node{}
	for _, n := range nodes {
		if maxUses >= 0 && n.Cards > maxUses {
			continue
		}
		if category != "" && !strings.Contains(strings.ToLower(n.Category), category) {
			continue
		}
		byCategory[n.Category] = append(byCategory[n.Category], n)
	}
	cats := make([]string, 0, len(byCategory))
	for c := range byCategory {
		cats = append(cats, c)
	}
	sort.Strings(cats)

	// Widen the name column to the longest name printed, so one long name shifts
	// the count column for every row rather than only for its own.
	width := len("NODE")
	for _, group := range byCategory {
		for _, n := range group {
			width = max(width, len(n.Name))
		}
	}

	fmt.Printf("  %-*s %6s\n", width, "NODE", "CARDS")
	for _, c := range cats {
		group := byCategory[c]
		sort.Slice(group, func(i, j int) bool {
			if group[i].Cards != group[j].Cards {
				return group[i].Cards < group[j].Cards
			}
			return group[i].Name < group[j].Name
		})
		fmt.Printf("\n%s (%d)\n", c, len(group))
		for _, n := range group {
			fmt.Printf("  %-*s %6d\n", width, n.Name, n.Cards)
		}
	}

	var zero, single int
	for _, n := range nodes {
		switch n.Cards {
		case 0:
			zero++
		case 1:
			single++
		}
	}
	fmt.Printf("\n%d exported names: %d unused, %d used by one card (%d total at most one)\n",
		len(nodes), zero, single, zero+single)
}
