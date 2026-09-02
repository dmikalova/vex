// Command genrules assembles the Vactrol rulebook (docs/rulebook.md) from the
// doc comments on the engine's implementation. It is the source of
// `mage generateRules` (and half of `mage gen`).
//
// The idea mirrors gencomments: keep the rules text next to the code that
// enforces them so the two can never drift. Code opts into the rulebook by
// carrying a `//rulebook:<section> <Title>` directive in a comment; the rest of
// that comment group is the entry's body. The directive can sit in a
// declaration's doc comment, or in a standalone comment block just above it, so a
// function can keep an ordinary implementation doc comment and still contribute a
// player-facing rule. A title may add a ` / <subheading>` suffix to gather
// several code sites under one heading. Entries are grouped by section and sorted
// alphabetically by title, so adding a keyword or effect just drops it into its
// section without renumbering anything.
//
// Foundational prose that isn't tied to any one declaration — the overview and
// each section's intro — lives as Markdown under docs/rulebook/ and is spliced in
// around the generated entries.
package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const (
	// defaultRoot is the package tree scanned for rulebook directives.
	defaultRoot = "internal/engine"
	// proseDir holds the hand-written Markdown fragments (overview + section
	// intros) that frame the generated entries.
	proseDir = "docs/rulebook"
	// outputFile is the assembled rulebook.
	outputFile = "docs/rulebook.md"
	// directive is the doc-comment marker that files a declaration into the
	// rulebook. gofmt treats it as a directive and keeps it verbatim.
	directive = "rulebook:"
)

// section is one top-level rulebook section. Sections print in this fixed order;
// the entries within each are sorted alphabetically by title, so the order of
// terms is data-driven and never has to be maintained by hand.
type section struct {
	key   string // directive key that files entries here, e.g. "keyword"
	title string // heading printed for the section, e.g. "Keywords"
}

// sections is the ordered spine of the rulebook. Add a section by adding a row
// here (and, optionally, a docs/rulebook/<key>.md intro).
var sections = []section{
	{"turn", "The Turn"},
	{"combat", "Combat"},
	{"cardtype", "Card Types"},
	{"keyword", "Keywords"},
	{"ability", "Abilities"},
	{"effect", "Effects"},
}

// entry is one harvested rulebook term: its title, an optional subheading that
// groups several entries under one title, and its body Markdown.
type entry struct {
	title    string
	subtitle string
	body     string
}

func main() {
	roots := os.Args[1:]
	if len(roots) == 0 {
		roots = []string{defaultRoot}
	}

	byKey := map[string][]entry{}
	for _, root := range roots {
		if err := harvest(root, byKey); err != nil {
			fmt.Fprintln(os.Stderr, "genrules:", err)
			os.Exit(1)
		}
	}

	doc, err := render(byKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, "genrules:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outputFile, []byte(doc), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "genrules:", err)
		os.Exit(1)
	}

	var n int
	for _, es := range byKey {
		n += len(es)
	}
	fmt.Printf("genrules: wrote %s (%d entries)\n", outputFile, n)
}

// harvest walks root, parsing every non-test .go file and collecting the
// rulebook entry from each comment group that carries a directive. It scans all
// comments, not just doc comments, so a rule can live in its own comment block
// sitting above a declaration's ordinary (implementation) doc comment.
func harvest(root string, byKey map[string][]entry) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() ||
			filepath.Ext(path) != ".go" ||
			isTestFile(path) {
			return nil
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		for _, cg := range file.Comments {
			collect(cg, byKey)
		}
		return nil
	})
}

// collect files the entry described by doc, if any, under its section key. A
// title may carry a " / <subheading>" suffix, which groups it with other entries
// of the same title under one heading.
func collect(doc *ast.CommentGroup, byKey map[string][]entry) {
	key, rawTitle, body, ok := parseDirective(doc)
	if !ok {
		return
	}
	title, subtitle, _ := strings.Cut(rawTitle, " / ")
	byKey[key] = append(byKey[key], entry{
		title:    strings.TrimSpace(title),
		subtitle: strings.TrimSpace(subtitle),
		body:     body,
	})
}

// parseDirective extracts the section key, title, and body from a doc comment.
// The body is every non-directive comment line, so it is unaffected by where
// gofmt places the directive within the comment. ok is false when the comment is
// absent or carries no rulebook directive.
func parseDirective(doc *ast.CommentGroup) (key, title, body string, ok bool) {
	if doc == nil {
		return "", "", "", false
	}
	var bodyLines []string
	for _, c := range doc.List {
		line := commentText(c.Text)
		if k, t, isDir := directiveOf(line); isDir {
			key, title, ok = k, t, true
			continue
		}
		bodyLines = append(bodyLines, strings.TrimPrefix(line, " "))
	}
	if !ok {
		return "", "", "", false
	}
	return key, title, joinBody(bodyLines), true
}

// commentText strips the comment markers from a single comment token, leaving the
// raw line (leading space, if any, preserved so paragraph text keeps its shape).
func commentText(raw string) string {
	if s, ok := strings.CutPrefix(raw, "//"); ok {
		return s
	}
	s := strings.TrimPrefix(raw, "/*")
	s = strings.TrimSuffix(s, "*/")
	return s
}

// directiveOf reports whether line is a rulebook directive and, if so, returns
// its section key and title.
func directiveOf(line string) (key, title string, ok bool) {
	trimmed := strings.TrimSpace(line)
	rest, ok := strings.CutPrefix(trimmed, directive)
	if !ok {
		return "", "", false
	}
	rest = strings.TrimSpace(rest)
	key, title, _ = strings.Cut(rest, " ")
	return key, strings.TrimSpace(title), true
}

// joinBody joins body lines, trims surrounding blank lines, and collapses runs of
// blank lines to a single paragraph break.
func joinBody(lines []string) string {
	var out []string
	prevBlank := true // drop leading blanks
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			if prevBlank {
				continue
			}
			prevBlank = true
			out = append(out, "")
			continue
		}
		prevBlank = false
		out = append(out, l)
	}
	return strings.TrimRight(strings.Join(out, "\n"), "\n")
}

// render assembles the full rulebook Markdown from the harvested entries and the
// prose fragments.
func render(byKey map[string][]entry) (string, error) {
	var b strings.Builder
	b.WriteString("<!-- Code generated by magefiles/genrules; DO NOT EDIT. -->\n")
	b.WriteString("<!-- Term entries come from doc comments under " + defaultRoot + "; -->\n")
	b.WriteString("<!-- prose lives in " + proseDir + "/. Run `mage gen` to regenerate. -->\n\n")

	overview, err := readProse("overview")
	if err != nil {
		return "", err
	}
	if overview != "" {
		b.WriteString(overview)
		b.WriteString("\n\n")
	}

	var index []indexEntry
	for _, sec := range sections {
		entries := byKey[sec.key]
		intro, err := readProse(sec.key)
		if err != nil {
			return "", err
		}
		if len(entries) == 0 && intro == "" {
			continue
		}
		fmt.Fprintf(&b, "## %s\n\n", sec.title)
		if intro != "" {
			b.WriteString(intro)
			b.WriteString("\n\n")
		}
		for _, gr := range sectionLead(groupByTitle(entries), sec.title) {
			// An entry whose title matches its section (e.g. Combat) would render a
			// heading identical to the section's, so emit its body directly under
			// the section instead of repeating the heading.
			if !strings.EqualFold(gr.title, sec.title) {
				fmt.Fprintf(&b, "### %s\n\n", gr.title)
			}
			index = append(index, indexEntry{gr.title, anchor(gr.title)})
			// Untitled parts render first as the heading's own body; subheaded
			// parts follow, ordered by subheading (number them to control order).
			var subbed []entry
			for _, p := range gr.parts {
				if p.subtitle == "" {
					writeBody(&b, p.body)
					continue
				}
				subbed = append(subbed, p)
			}
			sort.SliceStable(subbed, func(i, j int) bool {
				return strings.ToLower(subbed[i].subtitle) < strings.ToLower(subbed[j].subtitle)
			})
			for _, p := range subbed {
				fmt.Fprintf(&b, "#### %s\n\n", p.subtitle)
				index = append(index, indexEntry{p.subtitle, anchor(p.subtitle)})
				writeBody(&b, p.body)
			}
		}
	}

	writeIndex(&b, index)

	return strings.TrimRight(b.String(), "\n") + "\n", nil
}

// indexEntry is one term in the trailing alphabetical index: the heading text and
// the GitHub-style anchor slug that links to it.
type indexEntry struct {
	title  string
	anchor string
}

// writeIndex appends the alphabetical index of every term heading, each linking
// back to its section. It is emitted once, after all sections.
func writeIndex(b *strings.Builder, index []indexEntry) {
	if len(index) == 0 {
		return
	}
	sort.SliceStable(index, func(i, j int) bool {
		return strings.ToLower(index[i].title) < strings.ToLower(index[j].title)
	})
	b.WriteString("## Index\n\n")
	for _, e := range index {
		fmt.Fprintf(b, "- [%s](#%s)\n", e.title, e.anchor)
	}
	b.WriteString("\n")
}

// anchor renders a term title as a GitHub-style heading slug: lowercased, with
// spaces and hyphens collapsed to hyphens and every other character dropped.
// Letters outside ASCII are kept — GitHub slugs "Æmber" as "æmber".
func anchor(title string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(title) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		case r == ' ' || r == '-':
			b.WriteByte('-')
		}
	}
	return b.String()
}

// titleGroup gathers every entry sharing a title so several code sites can feed
// one rulebook heading.
type titleGroup struct {
	title string
	parts []entry
}

// groupByTitle merges entries that share a title and returns the groups sorted
// alphabetically by title, each keeping its parts in harvest order.
func groupByTitle(entries []entry) []titleGroup {
	var groups []titleGroup
	index := map[string]int{}
	for _, e := range entries {
		i, ok := index[e.title]
		if !ok {
			i = len(groups)
			index[e.title] = i
			groups = append(groups, titleGroup{title: e.title})
		}
		groups[i].parts = append(groups[i].parts, e)
	}
	sort.SliceStable(groups, func(i, j int) bool {
		return strings.ToLower(groups[i].title) < strings.ToLower(groups[j].title)
	})
	return groups
}

// sectionLead reorders groups so the entry whose title equals the section (its
// overview, rendered without a repeated heading) leads, ahead of the ### terms.
func sectionLead(groups []titleGroup, section string) []titleGroup {
	sort.SliceStable(groups, func(i, j int) bool {
		return strings.EqualFold(groups[i].title, section) &&
			!strings.EqualFold(groups[j].title, section)
	})
	return groups
}

// writeBody appends a non-empty body paragraph followed by a blank line.
func writeBody(b *strings.Builder, body string) {
	if body != "" {
		b.WriteString(body)
		b.WriteString("\n\n")
	}
}

// readProse returns the trimmed contents of docs/rulebook/<name>.md, or "" when
// that fragment does not exist.
func readProse(name string) (string, error) {
	data, err := os.ReadFile(filepath.Join(proseDir, name+".md"))
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// isTestFile reports whether path is a _test.go file.
func isTestFile(path string) bool {
	return strings.HasSuffix(filepath.Base(path), "_test.go")
}
