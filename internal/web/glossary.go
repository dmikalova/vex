package web

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file is the /glossary page: every rule term as an alphabetical
// Title→Definition list. It draws engine.Glossary() — the glossary facet of the
// same registry the rulebook draws from — so the two pages cannot list different
// terms. Definitions are one-line Rules-voice glosses; a term whose Definition is
// not yet written shows a placeholder until it is.

// glossaryPending marks a term whose one-line Definition has not been written yet.
const glossaryPending = "—"

// NewGlossary returns the root component for the /glossary page.
func NewGlossary() app.Composer { return &glossary{} }

// glossary is the /glossary page component. It holds no state: every render walks
// the engine's glossary afresh.
type glossary struct {
	app.Compo
}

// Render draws the glossary as a definition list, one term per row.
func (g *glossary) Render() app.UI {
	entries := engine.Glossary()
	rows := make([]app.UI, 0, len(entries)*2)
	for _, e := range entries {
		def := app.Dd().Class("glossary-def")
		if e.Definition == "" {
			def = def.Text(glossaryPending)
		} else {
			def = def.Body(inlineMarkdown(e.Definition)...)
		}
		rows = append(rows,
			app.Dt().Class("glossary-term").Text(e.Title),
			def,
		)
	}
	return app.Div().Class("doc-page").Body(
		docHeader("/glossary"),
		app.Section().Class("doc-section").Body(
			app.H2().Class("doc-h2").Text("Glossary"),
			app.Dl().Class("glossary").Body(rows...),
		),
	)
}
