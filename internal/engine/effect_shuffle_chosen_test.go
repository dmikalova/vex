package engine

import "testing"

// TestShuffleChosenCreaturesFromDiscard covers Not Finished with You: the
// controller shuffles any number of chosen creatures out of their discard pile.
func TestShuffleChosenCreaturesFromDiscard(t *testing.T) {
	// Text renders the plain and house-filtered forms.
	if got := (ShuffleChosenCreaturesFromDiscard{}).Text(); got !=
		"shuffle any number of creatures from your discard pile into your deck" {
		t.Errorf("text = %q", got)
	}
	if got := (ShuffleChosenCreaturesFromDiscard{House: Untamed}).Text(); got !=
		"shuffle any number of Untamed creatures from your discard pile into your deck" {
		t.Errorf("house text = %q", got)
	}

	// The default chooser takes every eligible creature, skipping a non-creature.
	g := NewGame("A", "B", 1)
	a := g.AddToDiscard(testCreature("a", 3), 0)
	b := g.AddToDiscard(testCreature("b", 3), 0)
	tactic := g.AddToDiscard(NewCard("t", Dis, Tactic, Common), 0)
	ShuffleChosenCreaturesFromDiscard{}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, Source: a},
	)
	if g.State.Deck[0].Count != 2 {
		t.Errorf("deck count = %d, want 2", g.State.Deck[0].Count)
	}
	if !containsID(g.Discard(0), tactic) {
		t.Error("the non-creature should stay in the discard pile")
	}
	if containsID(g.Discard(0), a) || containsID(g.Discard(0), b) {
		t.Error("chosen creatures should have left the discard pile")
	}

	// The house filter narrows the eligible creatures.
	g2 := NewGame("A", "B", 1)
	untamed := g2.AddToDiscard(NewCard("u", Untamed, Creature, Common, WithPower(1)), 0)
	g2.AddToDiscard(NewCard("d", Dis, Creature, Common, WithPower(1)), 0)
	ShuffleChosenCreaturesFromDiscard{House: Untamed}.
		Resolve(&EffectContext{Resolver: g2, Controller: 0, Source: untamed})
	if g2.State.Deck[0].Count != 1 || containsID(g2.Discard(0), untamed) {
		t.Error("only the Untamed creature should have been shuffled away")
	}

	// Declining shuffles nothing.
	g3 := NewGame("A", "B", 1)
	kept := g3.AddToDiscard(testCreature("kept", 3), 0)
	g3.SetChooser(0, &cardDecliner{decline: true})
	ShuffleChosenCreaturesFromDiscard{}.
		Resolve(&EffectContext{Resolver: g3, Controller: 0, Source: kept})
	if g3.State.Deck[0].Count != 0 || !containsID(g3.Discard(0), kept) {
		t.Error("declining should shuffle nothing")
	}

	// Outside a shuffle batch the move narrates itself directly.
	g4 := NewGame("A", "B", 1)
	lone := g4.AddToDiscard(testCreature("lone", 3), 0)
	g4.ShuffleFromDiscardIntoDeck(lone)
	if g4.State.Deck[0].Count != 1 || containsID(g4.Discard(0), lone) {
		t.Error("a lone shuffle should move the card into the deck")
	}
}
