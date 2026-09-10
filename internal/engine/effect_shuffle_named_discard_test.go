package engine

import "testing"

// TestShuffleNamedFromDiscardIntoDeck covers Chain Gang: one card of the named
// card leaves the discard for the deck while a differently-named card stays put.
func TestShuffleNamedFromDiscardIntoDeck(t *testing.T) {
	if got := (ShuffleNamedFromDiscardIntoDeck{Name: "Subtle Chain"}).Text(); got !=
		"shuffle a Subtle Chain from your discard pile into your deck" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	chain := g.AddToDiscard(NewCard("Subtle Chain", Dis, Tactic, Common), 0)
	other := g.AddToDiscard(NewCard("Mind Barb", Dis, Tactic, Common), 0)

	ShuffleNamedFromDiscardIntoDeck{Name: "Subtle Chain"}.
		Resolve(&EffectContext{Resolver: g, Controller: 0, Source: other})

	if containsID(g.Discard(0), chain) {
		t.Error("the named card should have left the discard pile")
	}
	if !containsID(g.Discard(0), other) {
		t.Error("a differently-named card should stay in the discard pile")
	}
	if g.State.Deck[0].Count != 1 {
		t.Errorf("deck count = %d, want 1", g.State.Deck[0].Count)
	}

	// With no card of that name in the discard pile, the effect shuffles nothing.
	g2 := NewGame("A", "B", 1)
	only := g2.AddToDiscard(NewCard("Mind Barb", Dis, Tactic, Common), 0)
	ShuffleNamedFromDiscardIntoDeck{Name: "Subtle Chain"}.
		Resolve(&EffectContext{Resolver: g2, Controller: 0, Source: only})
	if g2.State.Deck[0].Count != 0 || !containsID(g2.Discard(0), only) {
		t.Error("a discard pile without the named card should be left untouched")
	}
}
