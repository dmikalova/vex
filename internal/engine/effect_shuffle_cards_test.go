package engine

import "testing"

// TestShuffleCardsFromDiscard covers Shard of Life: the controller shuffles one
// card from their discard pile back into their deck for each friendly Shard.
func TestShuffleCardsFromDiscard(t *testing.T) {
	count := InPlay{Player: Controller, Trait: Shard}

	if got := (ShuffleCardsFromDiscard{Count: count}).Text(); got !=
		"for each friendly Shard, shuffle a card from your discard pile into your deck" {
		t.Errorf("text = %q", got)
	}
	if (ShuffleCardsFromDiscard{}).validate() == nil {
		t.Error("a missing Count should fail validation")
	}
	if err := (ShuffleCardsFromDiscard{Count: count}).validate(); err != nil {
		t.Errorf("a set Count should validate, got %v", err)
	}

	// Two friendly Shards in play shuffle two discard cards back.
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("s1", 3, WithTraits(Shard)), 0)
	g.AddToBattleline(testCreature("s2", 3, WithTraits(Shard)), 0)
	g.AddToDiscard(NewCard("d1", Dis, Creature, Common, WithPower(1)), 0)
	g.AddToDiscard(NewCard("d2", Dis, Creature, Common, WithPower(1)), 0)
	g.AddToDiscard(NewCard("d3", Dis, Creature, Common, WithPower(1)), 0)
	src := g.AddArtifact(NewCard("Shard of Life", Untamed, Artifact, Rare), 0)
	ShuffleCardsFromDiscard{Count: count}.
		Resolve(&EffectContext{Resolver: g, Controller: 0, Source: src})
	if g.State.Deck[0].Count != 2 {
		t.Errorf("deck count = %d, want 2", g.State.Deck[0].Count)
	}
	if len(g.Discard(0)) != 1 {
		t.Errorf("discard size = %d, want 1", len(g.Discard(0)))
	}

	// No friendly Shard shuffles nothing.
	g2 := NewGame("A", "B", 1)
	g2.AddToDiscard(NewCard("d", Dis, Creature, Common, WithPower(1)), 0)
	src2 := g2.AddArtifact(NewCard("Shard of Life", Untamed, Artifact, Rare), 0)
	ShuffleCardsFromDiscard{Count: count}.
		Resolve(&EffectContext{Resolver: g2, Controller: 0, Source: src2})
	if g2.State.Deck[0].Count != 0 || len(g2.Discard(0)) != 1 {
		t.Error("with no friendly Shard nothing should be shuffled")
	}

	// Fewer discard cards than the count shuffles only what is there.
	g3 := NewGame("A", "B", 1)
	g3.AddToBattleline(testCreature("s", 3, WithTraits(Shard)), 0)
	g3.AddToBattleline(testCreature("s2", 3, WithTraits(Shard)), 0)
	only := g3.AddToDiscard(NewCard("only", Dis, Creature, Common, WithPower(1)), 0)
	src3 := g3.AddArtifact(NewCard("Shard of Life", Untamed, Artifact, Rare), 0)
	ShuffleCardsFromDiscard{Count: count}.
		Resolve(&EffectContext{Resolver: g3, Controller: 0, Source: src3})
	if g3.State.Deck[0].Count != 1 || containsID(g3.Discard(0), only) {
		t.Error("the single discard card should have been shuffled away")
	}
}
