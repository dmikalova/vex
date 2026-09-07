package engine

import "testing"

// TestEachPlayerDiscardsAndRefillsHandText covers the rendered sentence.
func TestEachPlayerDiscardsAndRefillsHandText(t *testing.T) {
	want := "each player discards their hand, then refills their hand as if it " +
		"were the end of their turn"
	if got := (EachPlayerDiscardsAndRefillsHand{}).Text(); got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}
	if err := (EachPlayerDiscardsAndRefillsHand{}).validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
}

// TestEachPlayerDiscardsAndRefillsHand checks both players discard their hand and
// draw a fresh full one.
func TestEachPlayerDiscardsAndRefillsHand(t *testing.T) {
	g := NewGame("A", "B", 1)
	old := g.AddToHand(testCreature("old", 3), 0)
	for i := 0; i < HandSize+2; i++ {
		g.AddToDeck(testCreature("fresh0", 1), 0)
	}
	for i := 0; i < HandSize+2; i++ {
		g.AddToDeck(testCreature("fresh1", 1), 1)
	}
	g.AddToHand(testCreature("theirOld", 3), 1)

	ctx := &EffectContext{Resolver: g, Controller: 0}
	EachPlayerDiscardsAndRefillsHand{}.Resolve(ctx)

	if g.State.Discard[0].Count == 0 || !g.State.Discard[0].contains(old) {
		t.Error("the controller's old hand should have been discarded")
	}
	if got := int(g.State.Hand[0].Count); got != HandSize {
		t.Errorf("controller hand = %d, want %d", got, HandSize)
	}
	if got := int(g.State.Hand[1].Count); got != HandSize {
		t.Errorf("opponent hand = %d, want %d", got, HandSize)
	}
}

// TestEachPlayerDiscardsAndRefillsHandRespectsChains checks a chained player draws
// fewer cards and sheds one chain, exactly like an end-of-turn draw.
func TestEachPlayerDiscardsAndRefillsHandRespectsChains(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Chains[0] = 6 // reduces the draw by one card
	for i := 0; i < HandSize; i++ {
		g.AddToHand(testCreature("hand", 1), 0)
		g.AddToDeck(testCreature("deck", 1), 0)
	}

	ctx := &EffectContext{Resolver: g, Controller: 0}
	EachPlayerDiscardsAndRefillsHand{}.Resolve(ctx)

	if got := int(g.State.Hand[0].Count); got != HandSize-1 {
		t.Errorf("chained hand = %d, want %d", got, HandSize-1)
	}
	if g.State.Chains[0] != 5 {
		t.Errorf("chains = %d, want 5 (one shed)", g.State.Chains[0])
	}
}
