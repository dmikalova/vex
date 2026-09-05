package web

import (
	"strings"
	"unicode"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// This file is the shared frame for the reference pages — /rulebook and
// /glossary. Both are read-only facets of the engine's one rule-term registry:
// the rulebook draws engine.RuleBook() (the assembled sections), the glossary
// draws engine.Glossary() (the same terms as an alphabetical list). The engine
// owns the ordering and the term set; these pages only draw them. The page frame
// (brand bar and cross-links) and the anchor-slug helper live here so the two
// pages stay one look; Markdown rendering lives in markdown.go.

// docLink is one entry in the reference pages' nav row.
type docLink struct {
	href  string
	label string
}

// docLinks is the nav shared by every reference page: back to the game and across
// to each sibling reference page.
var docLinks = []docLink{
	{"/", "Game"},
	{"/rulebook", "Rulebook"},
	{"/glossary", "Glossary"},
}

// docHeader is the reference pages' brand bar: the title and the cross-page nav.
// active is the href of the current page, rendered as inert text rather than a
// link to itself.
func docHeader(active string) app.UI {
	nav := make([]app.UI, 0, len(docLinks))
	for _, l := range docLinks {
		if l.href == active {
			nav = append(nav, app.Span().
				Class("doc-nav-link").Class("doc-nav-link--active").
				Text(l.label))
			continue
		}
		nav = append(nav, app.A().Class("doc-nav-link").Href(l.href).Text(l.label))
	}
	return app.Header().Class("doc-header").Body(
		app.Span().Class("doc-brand").Text("Vactrol"),
		app.Nav().Class("doc-nav").Body(nav...),
	)
}

// docSlug renders a heading title as a GitHub-style anchor slug — lowercased,
// with spaces and hyphens collapsed to hyphens — so the table of contents can link
// to the heading that carries the same slug as its id.
func docSlug(title string) string {
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
