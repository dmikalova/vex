package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// TestEveryRefinementDraws holds the family's gate from the other side: every
// catalogued refinement draws a real strip of existing assets, and two
// refinements that teach a player different words draw different pictures, so
// "the 3 most powerful creatures" cannot quietly draw what "each creature with
// the lowest power" draws.
func TestEveryRefinementDraws(t *testing.T) {
	assertFamilyDraws(t, "Refinement")
	assertFamilyRowsDrawApart(t, "Refinement")
}

// TestPowerRefinementDecorSaysHowManyItKeeps pins the family's grammar for the
// power rules: a tier keeps every creature at the extreme and wears the each
// stack, a singular selector keeps one and wears the chosen outline, and which
// end of the scale it takes is the comparison mark rather than a drawing of its
// own.
func TestPowerRefinementDecorSaysHowManyItKeeps(t *testing.T) {
	cases := map[string]struct {
		refinement engine.Refinement
		decor      decor
		comparison engine.Comparison
	}{
		"most powerful":  {engine.MostPowerful, decorChosen, engine.AtLeast},
		"least powerful": {engine.LeastPowerful{}, decorChosen, engine.AtMost},
		"highest power":  {engine.HighestPower{}, decorEach, engine.AtLeast},
		"lowest power":   {engine.LowestPower{}, decorEach, engine.AtMost},
	}
	for name, tc := range cases {
		gs := refinementGlyphs(tc.refinement)
		if len(gs) != 1 {
			t.Fatalf("%s drew %+v, want one glyph", name, gs)
		}
		if gs[0].asset != "power" {
			t.Errorf("%s drew %q, want the power noun", name, gs[0].asset)
		}
		if gs[0].decor != tc.decor || gs[0].comparison != tc.comparison {
			t.Errorf("%s drew decor %v comparison %v, want %v and %v",
				name, gs[0].decor, gs[0].comparison, tc.decor, tc.comparison)
		}
	}
}

// TestTopNRefinementPrintsItsNumeralOnlyAboveOne pins what the numeral says: the
// top-N form prints its count the way its clause does ("the 3 most powerful
// creatures"), and the one-creature form prints none ("the most powerful
// creature"), so the badge is never a "1" the card does not say.
func TestTopNRefinementPrintsItsNumeralOnlyAboveOne(t *testing.T) {
	three := refinementGlyphs(engine.MostPowerfulN{N: 3})
	if len(three) != 1 || three[0].qty != 3 {
		t.Errorf("MostPowerfulN{N: 3} drew %+v, want a 3 on the power noun", three)
	}
	one := refinementGlyphs(engine.MostPowerfulN{N: 1})
	if len(one) != 1 || one[0].qty != 0 {
		t.Errorf("MostPowerfulN{N: 1} drew %+v, want no numeral", one)
	}
}

// TestSamePowerRefinementDrawsEveryPick checks that both chosen-power forms draw
// the creatures they measure against: Dance of Doom's single pick names no side,
// and Quintrino Flux's pair is tinted friendly and enemy, so the two read apart.
func TestSamePowerRefinementDrawsEveryPick(t *testing.T) {
	one := refinementGlyphs(engine.SamePowerAsChosen)
	if len(one) != 2 || one[1].decor != decorChosen {
		t.Errorf("SamePowerAsChosen drew %+v, want power and one pick", one)
	}
	pair := refinementGlyphs(engine.SamePowerAsEitherChosen)
	if len(pair) != 3 {
		t.Fatalf("SamePowerAsEitherChosen drew %+v, want power and two picks", pair)
	}
	if pair[1].decor != decorChosen|decorFriendly || pair[2].decor != decorChosen|decorEnemy {
		t.Errorf("SamePowerAsEitherChosen picks drew %v and %v, want friendly then enemy",
			pair[1].decor, pair[2].decor)
	}
}

// TestSparingRefinementsDrawTheRefusalMark pins where the refusal mark belongs: a
// refinement naming what it spares rather than what it takes closes with the ban,
// the way every other barred thing is drawn. The complement draws its inner
// refinement first, so what is spared is visible rather than implied.
func TestSparingRefinementsDrawTheRefusalMark(t *testing.T) {
	spared := refinementGlyphs(engine.AllExcept{Inner: engine.MostPowerful})
	if len(spared) != 2 || spared[0].asset != "power" ||
		spared[1].asset != "glyph-ban" {
		t.Errorf("AllExcept drew %+v, want the inner refinement then the ban", spared)
	}
	kept := refinementGlyphs(engine.KeepPerSide{N: 3})
	if len(kept) != 3 || kept[2].asset != "glyph-ban" {
		t.Fatalf("KeepPerSide drew %+v, want both sides then the ban", kept)
	}
	if kept[0].qty != 3 || kept[1].qty != 3 {
		t.Errorf("KeepPerSide drew %d and %d per side, want 3 each", kept[0].qty, kept[1].qty)
	}
	// The portion takes what it names, so it draws the same pair without the ban.
	taken := refinementGlyphs(engine.PortionPerSide{Share: engine.ThirdRoundedUp})
	if len(taken) != 2 || taken[0].qty != 3 {
		t.Errorf("PortionPerSide drew %+v, want both sides under the fraction's 3", taken)
	}
}

// TestUnionRefinementDrawsEveryMember checks that a union draws what it
// composes: Standardized Testing's two tiers are two power nouns, and
// Regrettable Meteor's Filter member draws the axis it narrows by — the trait
// mark — with no noun of its own, since the target holding the union already
// drew one (filterAxisGlyphs).
func TestUnionRefinementDrawsEveryMember(t *testing.T) {
	tiers := refinementGlyphs(engine.AnyOf{Refinements: []engine.Refinement{
		engine.LowestPower{},
		engine.HighestPower{},
	}})
	if len(tiers) != 2 {
		t.Errorf("AnyOf over two tiers drew %+v, want both", tiers)
	}
	filters := refinementGlyphs(engine.AnyOf{Refinements: []engine.Refinement{
		engine.Filter{Trait: engine.Dinosaur},
	}})
	if len(filters) != 1 || filters[0].asset != "glyph-trait" {
		t.Errorf("AnyOf over a Filter drew %+v, want the trait mark alone", filters)
	}
}

// TestRefinementDrawsAfterTheNounItNarrows pins the composition rule: a
// refinement is not a clause of its own, so a target draws its noun first and the
// narrowing after it, and a target carrying no refinement draws its noun alone.
func TestRefinementDrawsAfterTheNounItNarrows(t *testing.T) {
	plain := targetGlyphs(engine.Target{Kind: engine.TargetEachCreature})
	if len(plain) != 1 || plain[0].asset != "type-creature" {
		t.Fatalf("an unrefined target drew %+v, want its noun alone", plain)
	}
	refined := targetGlyphs(
		engine.Target{Kind: engine.TargetEachCreature}.Refine(engine.MostPowerfulN{N: 3}),
	)
	if len(refined) != 2 || refined[0] != plain[0] || refined[1].asset != "power" {
		t.Errorf("a refined target drew %+v, want the noun then the refinement", refined)
	}
}
