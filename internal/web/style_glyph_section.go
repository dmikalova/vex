package web

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/engine"
)

// This file is the Glyph vocabulary section of the Style gallery: the glyph
// assets, then the whole node census with the strip the Iconography pass draws
// for each row, then the worked composition examples.
//
// The census listing is the review surface for drawing glyphs (ADR 0022). Adding
// a family's drawings changes what roughly sixty rows render, which cannot be
// judged from a diff — the diff shows which asset a case returns, never whether
// the strip reads as the mechanic. So the section shows every catalogued row
// side by side with the text its node prints, and a row that draws nothing is
// shown blank rather than omitted, because during the work the point is to see
// what is still missing.
//
// It is developer-facing. The player-facing legend lives on /rulebook.

// glyphRow is one catalogued node as the gallery shows it: what the census says
// it is, and what the pass draws for it.
type glyphRow struct {
	// name is the node's name in the engine's source, and text what it prints on
	// a card.
	name, text string
	// note is the census's classification in one phrase: the rulebook term the
	// node owes, or the reason it owes none. A plumbing row carries its reason so
	// its blank strip reads as a decision rather than a hole.
	note string
	// plumbing marks a row that owes no term, which is what licenses it to draw
	// nothing and fold into its children.
	plumbing bool
	// glyphs is the strip the pass drew, empty when the row draws nothing — either
	// because its family has no renderer yet or because it folds away.
	glyphs []glyph
}

// glyphFamilyView is one catalogued family's rows, with the state of its gate.
type glyphFamilyView struct {
	name  string
	gated bool
	rows  []glyphRow
}

// glyphCatalogViews renders the whole census through the Iconography pass, one
// view per family. It draws every row, gated or not: an ungated family's blank
// rows are exactly the work still to do.
func glyphCatalogViews() []glyphFamilyView {
	families := iconCatalog()
	out := make([]glyphFamilyView, 0, len(families))
	for _, family := range families {
		view := glyphFamilyView{name: family.name, gated: family.gated}
		view.rows = make([]glyphRow, 0, len(family.rows))
		for _, row := range family.rows {
			view.rows = append(view.rows, newGlyphRow(row))
		}
		out = append(out, view)
	}
	return out
}

// newGlyphRow renders one catalogued node for the gallery. A row whose family has
// no renderer yet keeps its empty strip, which is what the blank cell shows.
func newGlyphRow(row iconRow) glyphRow {
	out := glyphRow{
		name:     row.name,
		text:     row.text,
		note:     row.rules.Term,
		plumbing: !row.bearsTerm(),
	}
	if out.plumbing {
		out.note = row.rules.NoTerm
	}
	if row.draw != nil {
		out.glyphs, _ = row.draw()
	}
	return out
}

// glyphCatalogRows returns the cached census views, rendering them on first use.
func (s *style) glyphCatalogRows() []glyphFamilyView {
	if s.glyphCatalog == nil {
		s.glyphCatalog = glyphCatalogViews()
	}
	return s.glyphCatalog
}

// glyphVocabSection is the full glyph vocabulary: every glyph asset with its
// name, then the node census with the strip each row draws, then worked
// composition examples run through the real effectGlyphs pass, so nothing here
// can drift from what cards actually render.
func (s *style) glyphVocabSection() app.UI {
	body := []app.UI{
		app.H3().Class("style-h3").Text("Every glyph"),
		app.Div().Class("style-row style-row--icons").Body(glyphAssetGrid()...),
		app.H3().Class("style-h3").Text("The node census"),
		app.P().Class("style-note").Text(
			"Every node the engine catalogues, with the text it prints and the strip " +
				"the Iconography pass draws for it. A row that owes a rulebook term must " +
				"draw a real glyph and a distinct one; a plumbing row owes none and may " +
				"fold into its children. A blank strip on a term-bearing row is a mechanic " +
				"still to draw.",
		),
	}
	for _, family := range s.glyphCatalogRows() {
		body = append(body, glyphFamilyBlock(family))
	}
	body = append(body,
		app.H3().Class("style-h3").Text("Compositions"),
		app.P().Class("style-note").Text(
			"Nested effects as a real card forms them. A catalogued node stands alone, "+
				"so it never shows how a trigger, a wrapper and a target compose into one "+
				"line.",
		),
		app.Div().Class("style-row").Body(glyphCompositionRows()...),
	)
	return app.Div().Body(body...)
}

// glyphAssetGrid draws every glyph asset with its stem, the vocabulary the
// composed strips are built from.
func glyphAssetGrid() []app.UI {
	grid := make([]app.UI, 0, len(galleryIcons))
	for _, name := range galleryIcons {
		grid = append(grid, app.Div().Class("style-swatch").Body(
			icon(name, "icon-outline"),
			app.Span().Class("style-mono").Text(name),
		))
	}
	return grid
}

// glyphFamilyBlock draws one family's census: a heading carrying its row count
// and gate state, then a row per catalogued node.
func glyphFamilyBlock(family glyphFamilyView) app.UI {
	gate := "not gated"
	if family.gated {
		gate = "gated"
	}
	rows := make([]app.UI, 0, len(family.rows))
	for _, row := range family.rows {
		rows = append(rows, glyphCatalogRow(row))
	}
	return app.Div().Body(
		app.H4().Class("style-h4").Body(
			app.Text(family.name),
			app.Span().Class("style-glyph-count").Text(
				strconv.Itoa(len(family.rows))+" nodes · "+gate,
			),
		),
		app.Div().Class("style-glyph-rows").Body(rows...),
	)
}

// glyphCatalogRow draws one catalogued node: its strip, its name, the text it
// prints, and what the census says it owes. A row that draws nothing gets an
// empty cell marked as blank, so a gap is visible rather than absent.
func glyphCatalogRow(row glyphRow) app.UI {
	return app.Div().Class("style-glyph-row").Body(
		app.Div().Class(cx("style-glyph-strip", blankClass(len(row.glyphs) == 0))).Body(
			app.If(len(row.glyphs) > 0, func() app.UI {
				return iconLine(glyphLine{glyphs: row.glyphs})
			}),
		),
		app.Div().Class("style-glyph-meta").Body(
			app.Span().Class("style-mono").Text(row.name),
			app.Span().Class("style-glyph-text").Text(row.text),
			app.Span().Class(cx("style-glyph-note", plumbingClass(row.plumbing))).
				Text(glyphNote(row)),
		),
	)
}

// blankClass returns the modifier marking a strip cell that drew nothing.
func blankClass(blank bool) string {
	if blank {
		return "style-glyph-strip--blank"
	}
	return ""
}

// plumbingClass returns the modifier marking a note that records a plumbing
// reason rather than a rulebook term.
func plumbingClass(plumbing bool) string {
	if plumbing {
		return "style-glyph-note--plumbing"
	}
	return ""
}

// glyphNote is the row's classification as one phrase, labelling a plumbing
// reason so it does not read as a term title.
func glyphNote(row glyphRow) string {
	if row.plumbing {
		return "plumbing — " + row.note
	}
	return row.note
}

// glyphComposition is one hand-written worked example: a whole ability, trigger
// and all, as a card would carry it.
type glyphComposition struct {
	label    string
	triggers []string
	effect   engine.Effect
}

// glyphCompositions are the worked examples. They are kept beside the census
// listing rather than replaced by it: a catalogued node is filled only as far as
// it needs to render, so it never forms the nested shapes a real card does — a
// trigger head, a target, a wrapper around another effect — and those are what
// the strip has to stay readable under.
func glyphCompositions() []glyphComposition {
	return []glyphComposition{
		{"Play/Fight/Reap: destroy the least powerful creature",
			[]string{"glyph-play", "glyph-fight", "glyph-reap"}, engine.Destroy{
				Target: engine.Target{Kind: engine.TargetChosenCreature},
			}},
		{
			"Deal 3 damage to an enemy creature",
			[]string{"glyph-play"},
			engine.DealDamage{
				Amount: 3,
				Target: engine.Target{Kind: engine.TargetChosenEnemyCreature},
			},
		},
		{"Gain 1 Æmber", []string{"glyph-reap"}, engine.GainAember{Amount: 1}},
		{"Stun an enemy creature", []string{"glyph-fight"}, engine.Stun{
			Target: engine.Target{Kind: engine.TargetChosenEnemyCreature},
		}},
		{"Draw 2 cards", []string{"glyph-reap"}, engine.Draw{Amount: 2}},
		{"A creature cannot fight", []string{"glyph-action"}, engine.Restrict{
			Player:   engine.Opponent,
			Action:   engine.RestrictFighting,
			Duration: engine.OpponentNextTurn,
		}},
		{
			"Look at the top card of your deck",
			[]string{"glyph-play"},
			engine.LookAtTopOfDeck{Amount: 1},
		},
		{
			"Take control of an enemy creature",
			[]string{"glyph-play"},
			engine.TakeControl{
				Target: engine.Target{Kind: engine.TargetChosenEnemyCreature},
			},
		},
	}
}

// glyphCompositionRows draws the worked examples, ending with the abstract
// fallback glyph a mechanic gets before it is drawn.
func glyphCompositionRows() []app.UI {
	examples := glyphCompositions()
	rows := make([]app.UI, 0, len(examples)+1)
	for _, ex := range examples {
		gs, _ := effectGlyphs(ex.effect)
		rows = append(rows, app.Div().Class("style-swatch").Body(
			iconLine(glyphLine{
				triggers: ex.triggers,
				glyphs:   gs,
			}),
			app.Span().Class("style-mono").Text(ex.label),
		))
	}
	return append(rows, app.Div().Class("style-swatch").Body(
		iconLine(glyphLine{glyphs: []glyph{{asset: "glyph-unknown"}}}),
		app.Span().Class("style-mono").Text("a mechanic not yet drawn"),
	))
}
