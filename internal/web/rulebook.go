package web

import (
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file is the /rulebook page: the whole engine rulebook rendered as text. It
// draws engine.RuleBook() — the assembled, ordered sections built from the
// engine's rule-term registry — so the page is the rulebook, straight from the
// code that enforces it. The engine owns the rules; this page only draws them,
// rendering each term's Markdown body (see markdown.go) and topping the page with
// a table of contents.

// NewRulebook returns the root component for the /rulebook page.
func NewRulebook() app.Composer { return &rulebook{} }

// rulebook is the /rulebook page component. It holds no state: every render walks
// the engine's assembled rulebook afresh.
type rulebook struct {
	app.Compo
}

// Render draws the overview, a table of contents, then each rulebook section with
// its intro and term groups, in the engine's canonical order.
func (r *rulebook) Render() app.UI {
	book := engine.RuleBook()
	body := []app.UI{docHeader("/rulebook")}
	if overview := engine.RuleOverview(); overview != "" {
		body = append(body, app.Section().Class("doc-section").
			Body(renderMarkdown(overview)...))
	}
	body = append(body, r.toc(book))
	for _, sec := range book {
		body = append(body, r.section(sec))
	}
	return app.Div().Class("doc-page").Body(body...)
}

// toc builds the table of contents: each section links to its heading, and the
// term groups that carry their own heading nest beneath it. The section-titled
// overview group has no heading of its own, so it folds into the section link.
func (r *rulebook) toc(book []engine.RuleSection) app.UI {
	sections := make([]app.UI, 0, len(book))
	for _, sec := range book {
		var subs []app.UI
		for _, gr := range sec.Groups {
			if strings.EqualFold(gr.Title, sec.Title) {
				continue
			}
			subs = append(subs, app.Li().Body(
				app.A().Href("#"+docSlug(gr.Title)).Text(gr.Title),
			))
		}
		item := []app.UI{app.A().Href("#" + docSlug(sec.Title)).Text(sec.Title)}
		if len(subs) > 0 {
			item = append(item, app.Ul().Class("doc-toc-sub").Body(subs...))
		}
		sections = append(sections, app.Li().Body(item...))
	}
	return app.Nav().Class("doc-toc").Body(
		app.H2().Class("doc-h2").Text("Contents"),
		app.Ul().Class("doc-toc-list").Body(sections...),
	)
}

// section draws one rulebook section: its heading (an anchor target for the table
// of contents), its intro, and its groups.
func (r *rulebook) section(sec engine.RuleSection) app.UI {
	items := []app.UI{app.H2().Class("doc-h2").ID(docSlug(sec.Title)).Text(sec.Title)}
	if sec.Intro != "" {
		items = append(items, renderMarkdown(sec.Intro)...)
	}
	for _, gr := range sec.Groups {
		items = append(items, r.group(sec.Title, gr)...)
	}
	return app.Section().Class("doc-section").Body(items...)
}

// group draws one term group: a heading (an anchor target, suppressed when it
// repeats the section's own), then each part's subheading and Markdown body.
func (r *rulebook) group(sectionTitle string, gr engine.RuleGroup) []app.UI {
	var out []app.UI
	if !strings.EqualFold(gr.Title, sectionTitle) {
		out = append(out, app.H3().Class("doc-h3").ID(docSlug(gr.Title)).Text(gr.Title))
	}
	for _, p := range gr.Parts {
		if p.Subtitle != "" {
			out = append(out, app.H4().Class("doc-h4").Text(p.Subtitle))
		}
		out = append(out, renderMarkdown(p.Body)...)
	}
	return out
}
