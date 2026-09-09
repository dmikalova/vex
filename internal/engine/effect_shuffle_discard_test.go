package engine

import "testing"

// TestShuffleMatchingFromDiscardIntoDeck covers Low Dawn: every Untamed creature
// leaves the discard for the deck while other cards stay put.
func TestShuffleMatchingFromDiscardIntoDeck(t *testing.T) {
	if got := (ShuffleMatchingFromDiscardIntoDeck{House: Untamed, Type: Creature}).
		Text(); got != "shuffle each Untamed creature from your discard pile into your deck" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	src := g.AddToDiscard(NewCard("src", Untamed, Tactic, Common), 0)
	u1 := g.AddToDiscard(NewCard("u1", Untamed, Creature, Common, WithPower(3)), 0)
	u2 := g.AddToDiscard(NewCard("u2", Untamed, Creature, Common, WithPower(3)), 0)
	tactic := g.AddToDiscard(NewCard("ut", Untamed, Tactic, Common), 0)
	mars := g.AddToDiscard(NewCard("mc", Mars, Creature, Common, WithPower(3)), 0)

	ShuffleMatchingFromDiscardIntoDeck{House: Untamed, Type: Creature}.
		Resolve(&EffectContext{Resolver: g, Controller: 0, Source: src})

	if g.State.Deck[0].Count != 2 {
		t.Errorf("deck count = %d, want 2", g.State.Deck[0].Count)
	}
	if containsID(g.Discard(0), u1) || containsID(g.Discard(0), u2) {
		t.Error("Untamed creatures should have left the discard pile")
	}
	if !containsID(g.Discard(0), tactic) {
		t.Error("the Untamed tactic should stay in the discard pile")
	}
	if !containsID(g.Discard(0), mars) {
		t.Error("the Mars creature should stay in the discard pile")
	}

	// The batch narrates one grouped, source-attributed line.
	var grouped int
	for _, rec := range g.Log {
		if e, ok := rec.Entry.(CardsShuffledIntoDeckBy); ok {
			grouped++
			if e.Source != src {
				t.Errorf("grouped shuffle source = %v, want %v", e.Source, src)
			}
		}
		if _, ok := rec.Entry.(CardShuffledIntoDeck); ok {
			t.Error("a batched shuffle should not narrate a passive per-card line")
		}
	}
	if grouped != 1 {
		t.Errorf("grouped shuffle lines = %d, want 1", grouped)
	}

	// With nothing matching, the effect shuffles nothing and narrates nothing.
	g2 := NewGame("A", "B", 1)
	only := g2.AddToDiscard(NewCard("mc", Mars, Creature, Common, WithPower(3)), 0)
	ShuffleMatchingFromDiscardIntoDeck{House: Untamed, Type: Creature}.
		Resolve(&EffectContext{Resolver: g2, Controller: 0, Source: only})
	if g2.State.Deck[0].Count != 0 || !containsID(g2.Discard(0), only) {
		t.Error("a non-matching discard pile should be left untouched")
	}
}
