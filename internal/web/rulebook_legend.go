package web

import (
	"fmt"
	"sync"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// This file is the player-facing glyph legend: the term→glyph map /rulebook draws
// beside each heading, so a player who meets an unfamiliar picture on a card looks
// it up in the same place they look up the rule.
//
// The map needs no table of its own. A rule-bearing census row names the rulebook
// term it owes (engine.RulesBearing.Term) and the Iconography pass draws that same
// row (icon_catalog.go), so pairing the two gives every drawn glyph its term and
// every rule-bearing term its glyphs. Both halves are complete by construction —
// the census fails the engine's build on an unclassified node (ADR 0018), and the
// icon gate fails on a rule-bearing row that draws nothing (ADR 0022) — so the
// legend cannot drift from the engine for the same reason the rulebook cannot.
//
// Several nodes may name one term (the Archive nodes are all "Archive"), so a term
// carries every distinct strip its nodes draw rather than one picked from them.
// There is no /glyphs page: the developer gallery at /style already lists the whole
// census, and a third surface listing the same thing would be the one to go stale.

// termGlyphs is the legend, built once: each rulebook term mapped to the distinct
// strips the census draws for it, in catalog order. The census is fixed for the
// life of the program and /rulebook redraws it on every render, so it is built
// once rather than per render.
var termGlyphs = sync.OnceValue(buildTermGlyphs)

// buildTermGlyphs walks the catalog harness and files each rule-bearing row's
// strip under the term the row names. A row whose family has no renderer yet, and
// a plumbing row that folds into its children, draw nothing and are skipped: an
// empty strip teaches no picture.
func buildTermGlyphs() map[string][][]glyph {
	out := map[string][][]glyph{}
	seen := map[string]bool{}
	for _, family := range iconCatalog() {
		for _, row := range family.rows {
			gs := rowStrip(row)
			if len(gs) == 0 {
				continue
			}
			key := row.rules.Term + "\x00" + legendKey(gs)
			if seen[key] {
				continue
			}
			seen[key] = true
			out[row.rules.Term] = append(out[row.rules.Term], gs)
		}
	}
	return out
}

// rowStrip is the strip a row contributes to the legend: the picture its term
// teaches, with the catalogued node's own quantities dropped. A node is filled
// only as far as it needs to render, so its numbers — the 1 on "deal 1 damage",
// the 3 on "the most powerful 3" — are the sample's, not the term's, and two nodes
// of one term differing only there teach one picture, not two. A comparison mark
// stays: "power 3 or more" is a threshold the vocabulary draws, not a count.
//
// A row that teaches no term, or whose family has no renderer yet, contributes
// nothing.
func rowStrip(row iconRow) []glyph {
	if !row.bearsTerm() || row.draw == nil {
		return nil
	}
	gs := row.draw()
	for i := range gs {
		gs[i].qty = 0
	}
	return gs
}

// legendKey renders a strip to the comparable key the legend dedupes on: every
// axis a glyph draws, in order.
func legendKey(gs []glyph) string {
	return fmt.Sprintf("%v", gs)
}

// termLegend renders the glyphs the census draws for the term titled title, or
// nothing for a title no node names — a section overview, or a rule the strip
// never has to say. It returns the heading's trailing children, so a term with no
// glyph leaves its heading exactly as it was.
func termLegend(title string) []app.UI {
	strips := termGlyphs()[title]
	if len(strips) == 0 {
		return nil
	}
	items := make([]app.UI, 0, len(strips))
	for _, gs := range strips {
		items = append(items, app.Span().Class("doc-legend-strip").Body(
			iconLine(glyphLine{glyphs: gs}),
		))
	}
	return []app.UI{app.Span().Class("doc-legend").Body(items...)}
}

// legendNote tells a player what the glyphs beside the headings are. It is the
// one piece of prose the legend needs: the pictures themselves are drawn from the
// census, so nothing here names a glyph or a rule.
func legendNote() app.UI {
	return app.Section().Class("doc-section").Body(
		app.P().Class("doc-p").Class("doc-legend-note").Text(
			"A card prints its rules as glyphs. Each heading below shows the glyphs " +
				"its rule draws. A rule drawn more than one way shows each of them.",
		),
	)
}

// headingBody is a rulebook heading's children: its title, then the glyphs its
// term draws. A heading whose title names no catalogued node is left as the bare
// title.
func headingBody(title string) []app.UI {
	return append([]app.UI{app.Text(title)}, termLegend(title)...)
}
