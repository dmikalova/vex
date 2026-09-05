package web

import (
	"strings"
	"unicode"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// This file renders the rulebook's Markdown into go-app nodes. The rule text is
// authored as Markdown in the engine registry, so a body's headings, lists, and
// emphasis have to become real elements rather than the raw text a plain
// paragraph split would leave on screen. The subset here — ATX headings,
// unordered lists, paragraphs, and inline bold/italic/code — is exactly what the
// authored rules use.

// renderMarkdown turns a Markdown body into block-level nodes: headings, unordered
// lists, and paragraphs, each with inline formatting applied.
func renderMarkdown(body string) []app.UI {
	var out []app.UI
	lines := strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n")
	for i := 0; i < len(lines); {
		line := strings.TrimSpace(lines[i])
		switch {
		case line == "":
			i++
		case headingLevel(line) > 0:
			lvl := headingLevel(line)
			out = append(out, mdHeading(lvl, strings.TrimSpace(line[lvl:])))
			i++
		case strings.HasPrefix(line, "- "):
			var items []app.UI
			var cur strings.Builder
			flush := func() {
				if cur.Len() > 0 {
					items = append(items, app.Li().Body(inlineMarkdown(cur.String())...))
					cur.Reset()
				}
			}
			for i < len(lines) {
				t := strings.TrimSpace(lines[i])
				if t == "" {
					break
				}
				if item, ok := strings.CutPrefix(t, "- "); ok {
					flush()
					cur.WriteString(item)
				} else {
					cur.WriteByte(' ')
					cur.WriteString(t)
				}
				i++
			}
			flush()
			out = append(out, app.Ul().Class("doc-ul").Body(items...))
		default:
			var para strings.Builder
			for i < len(lines) {
				t := strings.TrimSpace(lines[i])
				if t == "" || headingLevel(t) > 0 || strings.HasPrefix(t, "- ") {
					break
				}
				if para.Len() > 0 {
					para.WriteByte(' ')
				}
				para.WriteString(t)
				i++
			}
			out = append(out, app.P().Class("doc-p").Body(inlineMarkdown(para.String())...))
		}
	}
	return out
}

// headingLevel reports the ATX heading level of a line (1–3), or 0 when the line
// is not a heading.
func headingLevel(line string) int {
	n := 0
	for n < len(line) && line[n] == '#' {
		n++
	}
	if n == 0 || n > 3 || n >= len(line) || line[n] != ' ' {
		return 0
	}
	return n
}

// mdHeading builds a heading element for the given level with the reference pages'
// heading classes.
func mdHeading(level int, text string) app.UI {
	switch level {
	case 1:
		return app.H1().Class("doc-h1").Body(inlineMarkdown(text)...)
	case 2:
		return app.H2().Class("doc-h2").Body(inlineMarkdown(text)...)
	default:
		return app.H3().Class("doc-h3").Body(inlineMarkdown(text)...)
	}
}

// inlineMarkdown parses inline **bold**, *italic*, _italic_, and `code`, returning
// the pieces as nodes. Underscore emphasis follows Markdown's word-boundary rule,
// so an identifier such as effect_restrict.go is left as literal text.
func inlineMarkdown(s string) []app.UI {
	var out []app.UI
	var lit strings.Builder
	flush := func() {
		if lit.Len() > 0 {
			out = append(out, app.Text(lit.String()))
			lit.Reset()
		}
	}
	r := []rune(s)
	for i := 0; i < len(r); {
		switch {
		case runesHavePrefix(r, i, "**"):
			if j := runesIndex(r, i+2, "**"); j >= 0 {
				flush()
				out = append(out, app.Strong().Body(inlineMarkdown(string(r[i+2:j]))...))
				i = j + 2
				continue
			}
		case r[i] == '*':
			if j := runeIndex(r, i+1, '*'); j > i+1 {
				flush()
				out = append(out, app.Em().Body(inlineMarkdown(string(r[i+1:j]))...))
				i = j + 1
				continue
			}
		case r[i] == '`':
			if j := runeIndex(r, i+1, '`'); j >= 0 {
				flush()
				out = append(out, app.Code().Text(string(r[i+1:j])))
				i = j + 1
				continue
			}
		case r[i] == '_' && underscoreOpens(r, i):
			if j := underscoreCloses(r, i+1); j >= 0 {
				flush()
				out = append(out, app.Em().Body(inlineMarkdown(string(r[i+1:j]))...))
				i = j + 1
				continue
			}
		}
		lit.WriteRune(r[i])
		i++
	}
	flush()
	return out
}

// runesHavePrefix reports whether the runes at i begin with prefix p.
func runesHavePrefix(r []rune, i int, p string) bool {
	pr := []rune(p)
	if i+len(pr) > len(r) {
		return false
	}
	for k, c := range pr {
		if r[i+k] != c {
			return false
		}
	}
	return true
}

// runesIndex returns the index of the next occurrence of p in r at or after from,
// or -1.
func runesIndex(r []rune, from int, p string) int {
	for k := from; k+len([]rune(p)) <= len(r); k++ {
		if runesHavePrefix(r, k, p) {
			return k
		}
	}
	return -1
}

// runeIndex returns the index of the next c in r at or after from, or -1.
func runeIndex(r []rune, from int, c rune) int {
	for k := from; k < len(r); k++ {
		if r[k] == c {
			return k
		}
	}
	return -1
}

// wordRune reports whether c can be part of a word, for the underscore-emphasis
// boundary rule.
func wordRune(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_'
}

// underscoreOpens reports whether the underscore at i can open emphasis: it must
// not sit after a word character (so effect_restrict.go stays literal) and must be
// followed by a non-space.
func underscoreOpens(r []rune, i int) bool {
	if i+1 >= len(r) || r[i+1] == ' ' {
		return false
	}
	return i == 0 || !wordRune(r[i-1])
}

// underscoreCloses returns the index of the underscore that closes emphasis opened
// at from-1, or -1: the closer must not sit before a word character and must not
// follow a space, and the emphasis may not be empty.
func underscoreCloses(r []rune, from int) int {
	for k := from; k < len(r); k++ {
		if r[k] != '_' || k == from || r[k-1] == ' ' {
			continue
		}
		if k+1 < len(r) && wordRune(r[k+1]) {
			continue
		}
		return k
	}
	return -1
}
