package web

import (
	"fmt"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// familyRows returns one named family's catalogued rows, so a family's own test
// reads the census rather than a second list of nodes that could drift from it.
func familyRows(t *testing.T, name string) []iconRow {
	t.Helper()
	for _, family := range iconCatalog() {
		if family.name == name {
			return family.rows
		}
	}
	t.Fatalf("the icon catalog walks no %s family", name)
	return nil
}

// assertFamilyDraws holds one family to the two things a drawn family owes, in
// the family's own terms: every rule-bearing row draws a strip of its own, and
// every glyph in that strip names an asset that exists.
//
// It says more than TestIconTotality, which walks only the rule-bearing rows: a
// plumbing row that draws a glyph naming no asset passes there and fails here. A
// missing asset fails nothing at all without it — a glyph with a typo'd stem
// renders as a broken image — so a family that draws without this check ships its
// typos.
//
// A row drawing nothing is not a failure: a silent plumbing row folds into the
// node holding it (engine.Fixed, engine.AlwaysMet), which is the whole point of
// the census's own plumbing classification.
func assertFamilyDraws(t *testing.T, name string) {
	t.Helper()
	for _, row := range familyRows(t, name) {
		gs := row.draw()
		if row.bearsTerm() && len(gs) == 0 {
			t.Errorf("%s %s draws no glyph of its own (ADR 0022)", name, row.name)
			continue
		}
		for _, g := range gs {
			if !assetExists(t, g.asset) {
				t.Errorf("%s %s draws %q, which has no web/assets/%s.svg",
					name, row.name, g.asset, g.asset)
			}
		}
	}
}

// assertFamilyRowsDrawApart holds one family to the other half of the drawing
// rule: two rows that teach a player different words must not draw the same
// picture. A family whose members all flatten into one summary noun passes
// assertFamilyDraws — every row draws something real — while saying nothing about
// which member it is, which is how spreadTargetGlyph drew all five spreads as one
// chosen creature before they were given shapes of their own (ADR 0047).
//
// Plumbing rows are skipped: they fold into the node holding them, and two of them
// folding into nothing is the point rather than a collision.
func assertFamilyRowsDrawApart(t *testing.T, name string) {
	t.Helper()
	seen := map[string]string{}
	for _, row := range familyRows(t, name) {
		if !row.bearsTerm() {
			continue
		}
		gs := row.draw()
		strip := fmt.Sprintf("%+v", gs)
		if prev, dup := seen[strip]; dup {
			t.Errorf(
				"%s rows %s and %s both draw %s; give one a shape of its own (ADR 0047)",
				name, prev, row.name, strip,
			)
		}
		seen[strip] = row.name
	}
}

// TestIconCatalogCoversCensus holds the harness to the engine's census in both
// directions: every family engine.Families() lists is either walked by
// iconCatalog or named in iconCatalogSkipped with its reason, and every family
// iconCatalog walks still exists in the census. A sixteenth catalog accessor
// added to the engine therefore fails here rather than going undrawn and
// unnoticed (ADR 0022).
func TestIconCatalogCoversCensus(t *testing.T) {
	walked := map[string]bool{}
	for _, family := range iconCatalog() {
		walked[family.name] = true
	}
	known := map[string]bool{}
	for _, family := range engine.Families() {
		known[family.Name] = true
		if walked[family.Name] || iconCatalogSkipped[family.Name] != "" {
			continue
		}
		t.Errorf(
			"census family %s is drawn by nothing: add it to iconCatalog, or to "+
				"iconCatalogSkipped with the reason (ADR 0022)",
			family.Name,
		)
	}
	for name := range walked {
		if !known[name] {
			t.Errorf("iconCatalog walks %s, which the census no longer declares", name)
		}
	}
}

// TestIconCatalogPairsWithCensus checks that every catalogued node found its
// census row. The harness pairs the typed accessor with engine.Families() by
// position, which holds because both come from one catalog call; a row with no
// name is that assumption breaking, and it would silently leave the style
// gallery showing a nameless, textless row.
func TestIconCatalogPairsWithCensus(t *testing.T) {
	for _, family := range iconCatalog() {
		for i, row := range family.rows {
			if row.name == "" {
				t.Errorf("%s row %d paired with no census row", family.name, i)
			}
		}
	}
}

// TestNoSilentRuleBearingRow pins the assumption the drawing rule rests on: a
// row that teaches a player a word is never one the census declares silent. The
// rule is "a rule-bearing row must draw, a plumbing row may draw nothing", and a
// silent rule-bearing row would be a third case — a node that owes a term but
// renders no text to draw from. None exists, so the rule needs no exception
// clause, and this test is what stops one appearing unremarked (ADR 0022).
func TestNoSilentRuleBearingRow(t *testing.T) {
	for _, family := range iconCatalog() {
		for _, row := range family.rows {
			if row.silent && row.bearsTerm() {
				t.Errorf(
					"%s row %s is silent yet owes the term %q; decide whether it draws "+
						"(ADR 0022)",
					family.name,
					row.name,
					row.rules.Term,
				)
			}
		}
	}
}
