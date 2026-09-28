package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

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
