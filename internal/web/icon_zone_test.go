package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// TestEveryQuantityDraws holds the Quantity family's gate from the other side:
// every way a verb can say how many cards it takes draws a real, distinct strip of
// existing assets, so "purge 2 cards" cannot draw what "purge up to 2 cards"
// draws.
func TestEveryQuantityDraws(t *testing.T) {
	assertFamilyDraws(t, "Quantity")
	assertFamilyRowsDrawApart(t, "Quantity")
}

// TestDeclinableQuantitiesWearTheMayMark pins the rule the family teaches: "up to
// N" and "any number" both let the controller stop early, and that stop is the
// only thing separating them from a flat count of the same size. The may mark is
// where it is drawn, on the noun the verb already prints — the same mark a
// declinable pick wears (optionalVerb).
func TestDeclinableQuantitiesWearTheMayMark(t *testing.T) {
	for _, tc := range []struct {
		name     string
		quantity engine.Quantity
		want     sigil
	}{
		{"takes", engine.Takes{N: engine.Fixed(2)}, sigilNone},
		{"up to", engine.UpTo{N: engine.Fixed(2)}, sigilMay},
		{"any number", engine.AnyNumber{}, sigilMay},
	} {
		got := quantityOn(glyph{asset: "zone-archives"}, tc.quantity)
		if got.sigil != tc.want {
			t.Errorf("%s drew sigil %v, want %v", tc.name, got.sigil, tc.want)
		}
	}
}

// TestQuantityBadgesTheVerbsOwnNoun checks the quantity reaches the strip of the
// verb carrying it without drawing a card of its own: Creeping Oblivion purges up
// to two cards, which is the purge's own pile wearing a 2 and the may mark, not a
// purge followed by a picture of two cards.
func TestQuantityBadgesTheVerbsOwnNoun(t *testing.T) {
	gs := effectGlyphs(engine.ArchiveCard{
		Selection: engine.Chosen{},
		Quantity:  engine.UpTo{N: engine.Fixed(2)},
	})
	if len(gs) != 1 {
		t.Fatalf("bounded archive = %+v, want one glyph", gs)
	}
	if gs[0].asset != "zone-archives" || gs[0].qty != 2 || gs[0].sigil != sigilMay {
		t.Errorf("bounded archive = %+v, want the archives wearing 2 and the may mark", gs[0])
	}
	plain := effectGlyphs(engine.ArchiveCard{Selection: engine.Chosen{}})
	if plain[0].qty != 0 || plain[0].sigil != sigilNone {
		t.Errorf("a single archive = %+v, want a bare pile with no numeral", plain[0])
	}
}
