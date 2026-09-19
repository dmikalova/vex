package engine

import "testing"

// TestDiscardHandText covers the per-player rendering and the unset-player guard.
func TestDiscardHandText(t *testing.T) {
	if got := (DiscardHand{Player: Controller}).Text(); got != "discard your hand" {
		t.Errorf("controller = %q", got)
	}
	if got := (DiscardHand{Player: Opponent}).Text(); got != "your opponent discards their hand" {
		t.Errorf("opponent = %q", got)
	}
	if got := (DiscardHand{Player: EachPlayer}).Text(); got != "each player discards their hand" {
		t.Errorf("each = %q", got)
	}
	if err := (DiscardHand{}).validate(); err == nil {
		t.Error("an unset player should be rejected")
	}
	if err := (DiscardHand{Player: Controller}).validate(); err != nil {
		t.Errorf("a set player should validate, got %v", err)
	}
}

// TestRefillHandText covers the per-player rendering and the unset-player guard.
func TestRefillHandText(t *testing.T) {
	if got := (RefillHand{Player: Controller}).Text(); got != "refill your hand as if it were the end of the turn" {
		t.Errorf("controller = %q", got)
	}
	if got := (RefillHand{Player: Opponent}).Text(); got != "your opponent refills their hand as if it were the end of their turn" {
		t.Errorf("opponent = %q", got)
	}
	if got := (RefillHand{Player: EachPlayer}).Text(); got != "each player refills their hand as if it were the end of their turn" {
		t.Errorf("each = %q", got)
	}
	if err := (RefillHand{}).validate(); err == nil {
		t.Error("an unset player should be rejected")
	}
	if err := (RefillHand{Player: Controller}).validate(); err != nil {
		t.Errorf("a set player should validate, got %v", err)
	}
}

// TestDiscardAndRefillHandEachPlayer checks both players discard their hand and
// draw a fresh full one (Punctuated Equilibrium's two passes).
func TestDiscardAndRefillHandEachPlayer(t *testing.T) {
	g := NewGame("A", "B", 1)
	old := g.AddToHand(testCreature("old", 3), 0)
	for range HandSize + 2 {
		g.AddToDeck(testCreature("fresh0", 1), 0)
	}
	for range HandSize + 2 {
		g.AddToDeck(testCreature("fresh1", 1), 1)
	}
	g.AddToHand(testCreature("theirOld", 3), 1)

	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	DiscardHand{Player: EachPlayer}.Resolve(ctx)
	RefillHand{Player: EachPlayer}.Resolve(ctx)

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

// TestDiscardAndRefillHandRespectsChains checks a chained player draws fewer cards
// and sheds one chain, exactly like an end-of-turn draw.
func TestDiscardAndRefillHandRespectsChains(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Chains[0] = 6 // reduces the draw by one card
	for range HandSize {
		g.AddToHand(testCreature("hand", 1), 0)
		g.AddToDeck(testCreature("deck", 1), 0)
	}

	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	DiscardHand{Player: EachPlayer}.Resolve(ctx)
	RefillHand{Player: EachPlayer}.Resolve(ctx)

	if got := int(g.State.Hand[0].Count); got != HandSize-1 {
		t.Errorf("chained hand = %d, want %d", got, HandSize-1)
	}
	if g.State.Chains[0] != 5 {
		t.Errorf("chains = %d, want 5 (one shed)", g.State.Chains[0])
	}
}

// TestDiscardHandControllerOnly covers the single-player (non-EachPlayer) path.
func TestDiscardHandControllerOnly(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToHand(testCreature("c", 2), 0)
	DiscardHand{Player: Controller}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
	if !g.State.Discard[0].contains(c) {
		t.Error("the controller's hand should be discarded")
	}
}
