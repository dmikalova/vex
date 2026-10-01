package web

import (
	"sync"

	"github.com/dmikalova/vex/internal/engine"
)

// This file is the Iconography pass's population (ADR 0022): the engine's node
// census, read as the rows the pass owes a glyph. It answers "which mechanics
// must be drawn", and it replaces the card walk that used to answer it — a walk
// over the registered cards reaches only the mechanics an implemented card
// happens to carry, roughly half the pool, so an effect node with no glyph
// passed the gate until someone wrote the card that used it.
//
// The census is complete by construction (ADR 0018): a node internal/engine
// declares and the census does not list fails the engine's own build. Drawing
// the population from it makes the icon gate as complete as the engine is.
//
// The harness walks the typed catalog accessors — EffectCatalog(),
// ConditionCatalog(), and the thirteen beside them — rather than
// engine.Families(), because Family.Rows erases a node to its type name and the
// pass needs the node itself to draw it. It reads Families() alongside for the
// one thing the typed rows cannot give: the text a node renders, whose renderer
// is unexported for every family but Effect, Condition, Count, CreatureVerb and
// Target.
//
// The text-bearing value enums (engine.Enums()) belong to the same contract but
// not to this harness: an Enum's rows are already erased to FamilyRow, so no
// value survives for the pass to draw.

// iconRow is one catalogued node as the Iconography pass sees it: what the
// census says the node owes a player, and the strip the pass draws for it.
type iconRow struct {
	// name is the member's name in the engine's source — a node family's concrete
	// type, a Target filter builder's method name.
	name string
	// text is what the node prints on a card, taken from its census row.
	text string
	// rules is the census's classification: the rulebook term the node owes, or
	// the reason it owes none. It decides what the row owes the strip — a
	// rule-bearing row must draw a real glyph and a distinct one, while a
	// plumbing row may draw nothing and fold into its children (ADR 0022). The
	// judgement lives in internal/engine beside the reason the node teaches no
	// rule, so the client keeps no allowlist of its own.
	rules engine.RulesBearing
	// silent carries the census's declaration that the node prints no text by
	// design. No row is both silent and rule-bearing (TestNoSilentRuleBearingRow),
	// so the drawing rule needs no exception clause for one.
	silent bool
	// node is the catalogued value itself, for a consumer that needs more than the
	// strip: TestGlyphFamiliesAreDisjoint runs each effect down glyphFamilies.
	node any
	// draw renders the node to its strip, reporting whether a real mapping covered
	// it. It is nil for a family with no renderer yet, which is how a not-yet-drawn
	// family reads as blank rather than as a wrong picture.
	draw func() ([]glyph, bool)
}

// bearsTerm reports whether the census says this node teaches a player a word,
// which is what obliges it to draw a real glyph of its own.
func (r iconRow) bearsTerm() bool { return r.rules.Term != "" }

// iconFamily is one catalogued node family as the Iconography pass is measured
// against it.
type iconFamily struct {
	// name is the family's name in the census, e.g. "Refinement".
	name string
	// rows is the family's census, in catalog order.
	rows []iconRow
	// gated reports whether the icon totality gate fails on this family, mirroring
	// engine's Family.Gated: a family whose glyphs are not drawn yet is shown by
	// the style gallery but fails nothing, and the change that draws the family
	// flips this switch.
	gated bool
}

// gate switches the icon totality gate on for the family, which a family does
// once every rule-bearing row in it draws.
func (f iconFamily) gate() iconFamily {
	f.gated = true
	return f
}

// censusFamilies memoizes engine.Families(), which rebuilds every family's rows
// on each call. The census is fixed for the life of the program, and the style
// gallery reads it once per family per render, so building it once keeps a
// hover-triggered re-render from rebuilding all sixteen families fifteen times.
var censusFamilies = sync.OnceValue(engine.Families)

// censusRows returns the named family's census rows, which carry each node's
// name, classification and rendered text. The engine's own text renderers are
// unexported for most families — a Refinement renders through clause, a Spread
// through spreadText — so the census is the only way the client reads what a
// node prints.
func censusRows(name string) []engine.FamilyRow {
	for _, f := range censusFamilies() {
		if f.Name == name {
			return f.Rows
		}
	}
	return nil
}

// newIconFamily pairs a family's typed catalog with its census rows — the typed
// side carries the node to draw, the census side its name, classification and
// rendered text — and binds draw to each node. draw is nil for a family whose
// glyphs are not drawn yet.
//
// The two sides pair by position: the engine's newFamily builds Family.Rows from
// the very slice the typed accessor returns, in order, so row i names node i.
// TestIconCatalogPairsWithCensus holds that, failing on a row the census did not
// name.
func newIconFamily[T any](
	name string,
	catalog []engine.Catalogued[T],
	draw func(T) ([]glyph, bool),
) iconFamily {
	census := censusRows(name)
	rows := make([]iconRow, len(catalog))
	// Indexed rather than ranged over the values, for the same reason the engine's
	// newFamily is: ruleguard's largeloopcopy rule panics on a ranged value whose
	// type still carries a type parameter.
	for i := range catalog {
		node := catalog[i].Node
		row := iconRow{rules: catalog[i].Rules, silent: catalog[i].Silent, node: node}
		if i < len(census) {
			row.name, row.text = census[i].Type, census[i].Text()
		}
		if draw != nil {
			row.draw = func() ([]glyph, bool) { return draw(node) }
		}
		rows[i] = row
	}
	return iconFamily{name: name, rows: rows}
}

// iconCatalogSkipped names the census families this harness does not walk, with
// the reason, so a family's absence reads as a decision rather than an omission.
// TestIconCatalogCoversCensus holds every other family to having an entry in
// iconCatalog.
var iconCatalogSkipped = map[string]string{
	"Destination": "the engine exports no DestinationCatalog; a destination is " +
		"drawn as part of the zone-move effect carrying it",
}

// iconCatalog returns every catalogued node family the Iconography pass is
// measured against: the fifteen typed accessors the engine exports, in the order
// engine.Families() lists them.
//
// Effect, Condition, Count, BonusIconSubject, Selection and Gather are gated
// today. The rest are shown by the style gallery and gate nothing until their
// glyphs are drawn, and the change that draws one flips its gate.
//
// Two families wait on more than drawing, both for the same reason — a value
// internal/web cannot read. Target's rows are one Target per filter builder, and
// those filters are unexported fields, so the pass cannot tell which filter a row
// applied until the replace-target-s-filter-methods goal makes them a comparable
// exported Filter. Refinement's rows are values of unexported types built by
// constructor functions — engine.PowerAtLeast(6) builds an unexported
// powerAtLeast{minPower: 6} — so the pass can compare a refinement against the
// exported singletons but cannot recognise a parameterised one at all.
//
// A family added to the engine must be added here too — TestIconCatalogCovers-
// Census fails otherwise — so a sixteenth accessor cannot arrive without an icon
// decision.
func iconCatalog() []iconFamily {
	return []iconFamily{
		newIconFamily("Effect", engine.EffectCatalog(), effectGlyphs).gate(),
		newIconFamily("Refinement", engine.RefinementCatalog(), nil),
		newIconFamily("Condition", engine.ConditionCatalog(), conditionGlyphs).gate(),
		newIconFamily("Count", engine.CountCatalog(), countGlyphs).gate(),
		newIconFamily(
			"BonusIconSubject",
			engine.BonusIconSubjectCatalog(),
			bonusIconSubjectGlyphs,
		).gate(),
		newIconFamily("CreatureVerb", engine.CreatureVerbCatalog(), nil),
		newIconFamily("Selection", engine.SelectionCatalog(), selectionGlyphs).gate(),
		newIconFamily("Spread", engine.SpreadCatalog(), nil),
		newIconFamily("TopAct", engine.TopActCatalog(), nil),
		newIconFamily("PerTarget", engine.PerTargetCatalog(), nil),
		newIconFamily("Loss", engine.LossCatalog(), nil),
		newIconFamily("RepeatGate", engine.RepeatGateCatalog(), nil),
		newIconFamily("Gather", engine.GatherCatalog(), gatherGlyphs).gate(),
		newIconFamily("Quantity", engine.QuantityCatalog(), nil),
		newIconFamily("Target", engine.TargetFilterCatalog(), nil),
	}
}

// forEachCatalogedEffect calls fn for every catalogued node that is an Effect,
// plumbing rows included, since the glyph family chain claims those too. It is
// the view for a consumer that needs the node rather than its strip — the
// disjointness test runs each effect down glyphFamilies, which a bound draw
// closure cannot do.
func forEachCatalogedEffect(fn func(row iconRow, e engine.Effect)) {
	for _, family := range iconCatalog() {
		for _, row := range family.rows {
			if e, ok := row.node.(engine.Effect); ok {
				fn(row, e)
			}
		}
	}
}

// forEachGatedIconRow calls fn for every rule-bearing row of every gated family:
// exactly the rows the icon gate holds to drawing a real, distinct glyph. A
// plumbing row is skipped — it may fold into its children and draw nothing — and
// so is a family whose glyphs are not drawn yet.
func forEachGatedIconRow(fn func(family iconFamily, row iconRow)) {
	for _, family := range iconCatalog() {
		if !family.gated {
			continue
		}
		for _, row := range family.rows {
			if row.bearsTerm() {
				fn(family, row)
			}
		}
	}
}
