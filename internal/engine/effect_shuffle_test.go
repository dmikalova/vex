package engine

import "testing"

// TestSwapDeckAndDiscard covers Reverse Time: the deck and discard pile trade
// places, and the new deck is shuffled.
func TestSwapDeckAndDiscard(t *testing.T) {
	if got := (SwapDeckAndDiscard{}).Text(); got != "swap your deck and your discard pile, then shuffle your deck" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	inDeck := g.AddToDeck(testCreature("deck card", 1), 0)
	inDiscard := g.AddToDiscard(testCreature("discard card", 1), 0)

	SwapDeckAndDiscard{}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != inDiscard {
		t.Error("the discard pile should have become the deck")
	}
	if g.State.Discard[0].Count != 1 || g.State.Discard[0].IDs[0] != inDeck {
		t.Error("the deck should have become the discard pile")
	}
}

// TestSwapDeckAndDiscardFlipsTheStack checks the swap turns the deck over rather
// than pouring it out card by card: the top of the deck ends up at the bottom of
// the new discard pile, in reversed order. Only this direction is observable —
// the pile that becomes the deck is shuffled straight afterwards.
func TestSwapDeckAndDiscardFlipsTheStack(t *testing.T) {
	g := NewGame("A", "B", 1)
	top := g.AddToDeck(testCreature("top", 1), 0)
	middle := g.AddToDeck(testCreature("middle", 1), 0)
	bottom := g.AddToDeck(testCreature("bottom", 1), 0)
	g.AddToDiscard(testCreature("discarded", 1), 0)

	SwapDeckAndDiscard{}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	// A discard pile is stored bottom-first, so the old top of the deck is now
	// buried deepest.
	want := []LocalID{top, middle, bottom}
	got := g.Discard(0)
	if len(got) != len(want) {
		t.Fatalf("discard = %v, want %v", got, want)
	}
	for i, id := range want {
		if got[i] != id {
			t.Fatalf("discard = %v, want %v (bottom first)", got, want)
		}
	}
}
