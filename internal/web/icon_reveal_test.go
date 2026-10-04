package web

import (
	"fmt"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// TestEveryTopActDraws holds the TopAct family's gate from the other side: every
// routing step a deck read can send its cards through draws a real, distinct strip
// of existing assets, so a step that archives the cards cannot draw what a step
// that discards them draws.
func TestEveryTopActDraws(t *testing.T) {
	assertFamilyDraws(t, "TopAct")
	assertFamilyRowsDrawApart(t, "TopAct")
}

// TestPutBackOrderedIsNotShuffled pins the pair this family could most easily
// collapse: ReorderRest puts the read cards back on top in an order the controller
// picks, and Shuffle puts them back at no known position. Drawing both as "the
// cards go back into the deck" would lose the only thing that separates Eyegor
// from Borr Nit.
func TestPutBackOrderedIsNotShuffled(t *testing.T) {
	ordered := topActGlyphs(engine.ReorderRest{})
	shuffled := topActGlyphs(engine.Shuffle{})
	if len(ordered) == 0 || len(shuffled) == 0 {
		t.Fatal("the two put-back steps drew no glyph")
	}
	if fmt.Sprintf("%+v", ordered) == fmt.Sprintf("%+v", shuffled) {
		t.Errorf("both put-back steps draw %+v, want a known position apart from none", ordered)
	}
	if shuffled[len(shuffled)-1].sigil != sigilShuffle {
		t.Errorf("the shuffling step = %+v, want the shuffle verb", shuffled)
	}
}

// TestRevealedCardsDrawWhereTheyGo checks the steps reach the strip of the read
// carrying them: Borr Nit reveals five cards, purges one and shuffles the rest,
// and before the steps were drawn every multi-step reveal drew the same bare
// "deck, eye" as a plain peek.
func TestRevealedCardsDrawWhereTheyGo(t *testing.T) {
	gs := effectGlyphs(engine.RevealTopOfDeck{
		Amount: 5,
		Then: []engine.TopAct{
			engine.ChooseAndMove{Cards: 1, Dest: engine.IntoPurge},
			engine.Shuffle{},
		},
	})
	want := []string{"zone-deck", "glyph-look", "card-back", "zone-purge", "zone-deck"}
	if len(gs) != len(want) {
		t.Fatalf("routed reveal = %+v, want %d glyphs", gs, len(want))
	}
	for i, asset := range want {
		if gs[i].asset != asset {
			t.Errorf("routed reveal glyph %d = %q, want %q", i, gs[i].asset, asset)
		}
	}
	if gs[len(gs)-1].sigil != sigilShuffle {
		t.Errorf("routed reveal = %+v, want the shuffle closing it", gs)
	}
}

// TestPartitionDrawsBothSidesOfTheSplit checks New Frontiers' split draws where
// each half of the read goes: the cards of the chosen house to one zone, and the
// others behind the branch mark that says they are the other way the same split
// falls. Drawing only the matching half would make a partition read as a plain
// move.
func TestPartitionDrawsBothSidesOfTheSplit(t *testing.T) {
	gs := topActGlyphs(engine.PartitionByChosenHouse{
		Matching: engine.IntoArchives,
		Rest:     engine.IntoDiscard,
	})
	if len(gs) != 3 {
		t.Fatalf("partition = %+v, want three glyphs", gs)
	}
	if gs[0].asset != "house-any" || gs[0].decor&decorChosen == 0 {
		t.Errorf("partition opens with %+v, want the chosen house", gs[0])
	}
	if gs[1].asset != "zone-archives" || gs[2].asset != "zone-discard" {
		t.Errorf("partition = %+v, want the matching zone then the rest's", gs)
	}
	if gs[2].punct != punctBranch {
		t.Errorf("the rest of the read opens with %v, want the branch mark", gs[2].punct)
	}
}
