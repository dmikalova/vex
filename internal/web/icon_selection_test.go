package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// TestEverySelectionDraws holds the family's gate from the other side: every
// catalogued selection draws a real strip of existing assets, so a verb that takes
// the top card cannot draw what a verb that takes a chosen one draws.
func TestEverySelectionDraws(t *testing.T) {
	assertFamilyDraws(t, "Selection")
}

// TestPickedCardDecorSaysHowItWasPicked pins the family's grammar: Chosen, Each
// and Self are one noun under three decorations, not three drawings, so the way
// a card was picked is carried by the same chosen/stack/self marks every other
// noun in the strip wears.
func TestPickedCardDecorSaysHowItWasPicked(t *testing.T) {
	cases := map[string]struct {
		selection engine.Selection
		want      decor
	}{
		"chosen": {engine.Chosen{Filter: engine.Filter{Type: engine.Creature}}, decorChosen},
		"each":   {engine.Each{Filter: engine.Filter{Type: engine.Creature}}, decorEach},
		"self":   {engine.Self{}, decorThis},
	}
	for name, tc := range cases {
		gs := selectionGlyphs(tc.selection)
		if len(gs) != 1 {
			t.Fatalf("%s drew %+v, want one glyph", name, gs)
		}
		if gs[0].decor != tc.want {
			t.Errorf("%s drew decor %v, want %v", name, gs[0].decor, tc.want)
		}
	}
}

// TestChosenDrawsItsIdentityFilters checks that the filters deciding which cards
// a pick can reach at all are drawn — the house it must belong to, the trait and
// the name it narrows by — in the order a count of the same cards draws them, so
// "choose an Untamed Beast creature named Flaxia" reads as one phrase rather than
// as a bare card.
func TestChosenDrawsItsIdentityFilters(t *testing.T) {
	gs := selectionGlyphs(engine.Chosen{
		Filter: engine.Filter{
			House: engine.HouseMatcher{Kind: engine.MatchNamedHouse, House: engine.Untamed},
			Type:  engine.Creature,
			Trait: engine.Beast,
			Name:  "Flaxia",
		},
	})
	if len(gs) != 4 {
		t.Fatalf("filtered pick = %+v, want four glyphs", gs)
	}
	want := []string{"house-untamed", "type-creature", "glyph-trait", "glyph-name"}
	for i, asset := range want {
		if gs[i].asset != asset {
			t.Errorf("filtered pick glyph %d = %q, want %q", i, gs[i].asset, asset)
		}
	}
}

// TestOptionalPickWearsTheMayVerbOnItsCard pins where the may mark lands: on the
// card the controller may decline, not on the head of the run, which a house
// filter occupies. "You may discard a Dis card" is the card that is optional.
func TestOptionalPickWearsTheMayVerbOnItsCard(t *testing.T) {
	gs := selectionGlyphs(engine.Chosen{
		Filter: engine.Filter{
			House: engine.HouseMatcher{Kind: engine.MatchNamedHouse, House: engine.Dis},
		},
		Optional: true,
	})
	if len(gs) != 2 {
		t.Fatalf("optional pick = %+v, want the house and the card", gs)
	}
	if gs[0].sigil != sigilNone {
		t.Errorf("the house wears %v, want no verb", gs[0].sigil)
	}
	if gs[1].sigil != sigilMay {
		t.Errorf("the card wears %v, want the may verb", gs[1].sigil)
	}
}

// TestPositionalPicksDrawOppositeEndsOfThePile pins that the top and the bottom
// of an ordered pile are drawn apart. They are the one pair in this family that
// could collapse into a single "a card from a pile" mark, and collapsing them
// would make Random Access Archives and a bottom-of-deck pick read alike.
func TestPositionalPicksDrawOppositeEndsOfThePile(t *testing.T) {
	top := selectionGlyphs(engine.Top{})
	bottom := selectionGlyphs(engine.Bottom{})
	if len(top) != 1 || len(bottom) != 1 {
		t.Fatalf("the positional picks drew %+v and %+v, want one glyph each", top, bottom)
	}
	if top[0].asset == bottom[0].asset {
		t.Errorf("the top and the bottom both draw %q, want opposite ends", top[0].asset)
	}
}
