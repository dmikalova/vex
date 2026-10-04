package web

import (
	"fmt"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// This file is part of the Iconography pass (ADR 0022): tests for the Filter
// family in icon_filter.go.
//
// The census holds one representative value per axis, which is what the family
// gate measures. The tests beside it walk the axis enums the engine exports
// (engine.Positions, engine.PowerBoundKinds, …), because a census row shows one
// value of its axis and says nothing about the other four places a creature can
// stand.

// TestEveryFilterAxisDraws is the Filter family's gate from the family's own
// side: every axis a card can narrow by draws a strip of existing assets, and no
// two axes draw the same one. Until this landed targetGlyph summarised the Kind
// and left every fine filter to the rules text, so "each creature" and "each
// damaged Mars creature with power 3 or lower" drew one picture.
func TestEveryFilterAxisDraws(t *testing.T) {
	assertFamilyDraws(t, "Filter")
	assertFamilyRowsDrawApart(t, "Filter")
}

// TestFilterPowerBoundsDraw walks engine.PowerBoundKinds() and fails on a bound
// that draws nothing or draws what another bound draws. The census catalogues
// one bound (power 3 or lower), so without this walk a parity bound or a bound
// read against the source card's own power could ship drawing the same "≤ 3" as
// the threshold beside it.
func TestFilterPowerBoundsDraw(t *testing.T) {
	seen := map[string]engine.PowerBoundKind{}
	for _, kind := range engine.PowerBoundKinds() {
		gs := powerBoundGlyphs(engine.PowerBound{Kind: kind, Amount: 3})
		assertGlyphsDrawn(t, gs, "power bound kind %d", kind)
		key := filterStrip(gs)
		if prev, dup := seen[key]; dup {
			t.Errorf("power bound kinds %d and %d both draw %s; mark one apart",
				prev, kind, key)
		}
		seen[key] = kind
	}
}

// TestFilterPositionsDraw walks engine.Positions() and fails on a battleline
// place that draws nothing or draws another place's picture. A place that fell
// through would draw nothing at all, and "each creature in the center of its
// controller's battleline" would read as "each creature".
func TestFilterPositionsDraw(t *testing.T) {
	seen := map[string]engine.Position{}
	for _, position := range engine.Positions() {
		gs := positionGlyphs(position)
		assertGlyphsDrawn(t, gs, "position %d", position)
		key := filterStrip(gs)
		if prev, dup := seen[key]; dup {
			t.Errorf("positions %d and %d both draw %s; mark one apart", prev, position, key)
		}
		seen[key] = position
	}
}

// TestFilterExclusionsDraw walks engine.Exclusions() and fails on an excluded
// card that draws nothing. All three print the same "other", so they are allowed
// to look alike in the text; the strip still says which card is left out, and a
// fourth exclusion cannot ship drawing nothing.
func TestFilterExclusionsDraw(t *testing.T) {
	for _, exclusion := range engine.Exclusions() {
		assertGlyphsDrawn(t, exclusionGlyphs(exclusion), "exclusion %d", exclusion)
	}
}

// TestFilterPresencesDraw walks the damage and Æmber presences and fails on one
// that draws nothing. Each pair is "has some" against "has none", which the
// strip tells apart with the refusal mark rather than with a second drawing.
func TestFilterPresencesDraw(t *testing.T) {
	for _, damage := range engine.DamagePresences() {
		assertGlyphsDrawn(t, damagePresenceGlyphs(damage), "damage presence %d", damage)
	}
	for _, aember := range engine.AemberPresences() {
		assertGlyphsDrawn(t, aemberPresenceGlyphs(aember), "Æmber presence %d", aember)
	}
	some := filterStrip(damagePresenceGlyphs(engine.DamageSome))
	if none := filterStrip(damagePresenceGlyphs(engine.DamageNone)); some == none {
		t.Errorf("a damaged and an undamaged creature both draw %s", some)
	}
	some = filterStrip(aemberPresenceGlyphs(engine.AemberSome))
	if none := filterStrip(aemberPresenceGlyphs(engine.AemberNone)); some == none {
		t.Errorf("Æmber on a card and none on it both draw %s", some)
	}
}

// TestAFilterReshapesTheNounItNarrows pins the three axes that replace the noun
// rather than trailing it: a type, a gigantic creature, and a card named
// outright are what the target *is*, the way the printed text turns "each
// creature" into "each artifact". Drawing them as trailing marks would say the
// target reaches creatures carrying an artifact.
func TestAFilterReshapesTheNounItNarrows(t *testing.T) {
	cases := []struct {
		name   string
		filter engine.Filter
		want   string
	}{
		{"type", engine.Filter{Type: engine.Artifact}, "type-artifact"},
		{"gigantic", engine.Filter{Gigantic: true}, "glyph-gigantic"},
		{"name", engine.Filter{Name: "Ancient Bear"}, "glyph-name"},
	}
	for _, c := range cases {
		gs := filterGlyphs(c.filter, glyph{asset: "type-creature", decor: decorEach})
		if len(gs) != 1 || gs[0].asset != c.want {
			t.Errorf("the %s axis draws %+v, want the one %s noun", c.name, gs, c.want)
			continue
		}
		if gs[0].decor != decorEach {
			t.Errorf("the %s axis dropped the noun's decor: %+v", c.name, gs[0])
		}
	}
}

// TestAMatchAnyFilterDividesItsAxes pins what the filter's "or" draws: the axes
// it disjoins are punctuated with the divider a ChooseOne's options use, so "each
// Mars or Robot creature" (EMP Blast) reads as one noun with two alternatives.
// A conjoining filter sets the same two axes and draws no divider at all.
func TestAMatchAnyFilterDividesItsAxes(t *testing.T) {
	both := engine.Filter{
		House: engine.HouseMatcher{Kind: engine.MatchNamedHouse, House: engine.Mars},
		Trait: engine.Robot,
	}
	noun := glyph{asset: "type-creature", decor: decorEach}
	for _, g := range filterGlyphs(both, noun) {
		if g.punct == punctDivider {
			t.Errorf("a conjoining filter drew the \"or\" divider on %+v", g)
		}
	}
	both.MatchAny = true
	gs := filterGlyphs(both, noun)
	if len(gs) != 3 {
		t.Fatalf("a Mars-or-Robot filter draws %+v, want the noun and two axes", gs)
	}
	if gs[1].punct != punctNone || gs[2].punct != punctDivider {
		t.Errorf("the divider sits wrong in %+v: want it on the second axis alone", gs)
	}
}

// TestAFilterInAUnionDrawsNoNoun pins the one place a filter draws without a
// noun: as a Refinement inside a union, where the target holding the union has
// already drawn one (Regrettable Meteor unions a trait against a power floor).
// Drawing a second noun there would read as two sets rather than one.
func TestAFilterInAUnionDrawsNoNoun(t *testing.T) {
	gs := refinementGlyphs(engine.AnyOf{Refinements: []engine.Refinement{
		engine.Filter{Trait: engine.Dinosaur},
		engine.Filter{Power: engine.PowerBound{Kind: engine.BoundAtLeast, Amount: 6}},
	}})
	want := []string{"glyph-trait", "power"}
	if len(gs) != len(want) {
		t.Fatalf("the union draws %+v, want one glyph per member", gs)
	}
	for i, asset := range want {
		if gs[i].asset != asset {
			t.Errorf("union member %d draws %q, want %q", i, gs[i].asset, asset)
		}
	}
}

// assertGlyphsDrawn fails unless a run draws at least one glyph and every glyph
// in it names an asset on disk. An empty run and an empty stem are both dropped
// silently at render, so neither fails anything without a check like this one
// (ADR 0022).
func assertGlyphsDrawn(t *testing.T, gs []glyph, format string, args ...any) {
	t.Helper()
	subject := fmt.Sprintf(format, args...)
	if len(gs) == 0 {
		t.Errorf("%s draws nothing; give it a glyph (ADR 0022)", subject)
		return
	}
	for _, g := range gs {
		if g.asset == "" {
			t.Errorf("%s draws a glyph naming no asset: %+v", subject, g)
			continue
		}
		if !assetExists(t, g.asset) {
			t.Errorf("%s: no web/assets/%s.svg", subject, g.asset)
		}
	}
}

// filterStrip renders a run of glyphs in full, every axis of every glyph, so two
// drawings can be compared for sameness where asset and decor alone would miss
// it (glyphKey): the comparison mark and the numeral are the whole difference
// between "with power 3 or lower" and "with odd power".
func filterStrip(gs []glyph) string {
	return fmt.Sprintf("%+v", gs)
}
