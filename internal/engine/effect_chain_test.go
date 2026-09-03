package engine

import "testing"

func TestGainChainsText(t *testing.T) {
	if got := (GainChains{Amount: 1}).Text(); got != "gain 1 chain" {
		t.Errorf("text = %q, want %q", got, "gain 1 chain")
	}
	if got := (GainChains{Amount: 2}).Text(); got != "gain 2 chains" {
		t.Errorf("text = %q, want %q", got, "gain 2 chains")
	}
}

func TestGainChainsResolve(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	GainChains{Amount: 2}.Resolve(ctx)
	if g.State.Chains[0] != 2 {
		t.Errorf("chains = %d, want 2", g.State.Chains[0])
	}
}

func TestChainsReduceDrawAndShed(t *testing.T) {
	g := NewGame("A", "B", 1)
	for i := 0; i < HandSize; i++ {
		g.AddToDeck(testCreature("d", 1), 0)
	}
	g.State.Chains[0] = 7 // 7 chains: draw two fewer

	g.EndPlayPhase(0)

	if got := len(g.Hand(0)); got != HandSize-2 {
		t.Errorf("hand = %d, want %d (draw reduced by 2 at 7 chains)", got, HandSize-2)
	}
	if g.State.Chains[0] != 6 {
		t.Errorf("chains = %d, want 6 (one shed after being blocked)", g.State.Chains[0])
	}
}

func TestChainsShedOnlyWhenBlocked(t *testing.T) {
	// A full hand with cards to draw is not blocked, so no chain is shed.
	g := NewGame("A", "B", 1)
	for i := 0; i < HandSize; i++ {
		g.AddToHand(testCreature("h", 1), 0)
	}
	g.AddToDeck(testCreature("d", 1), 0)
	g.State.Chains[0] = 3

	g.EndPlayPhase(0)

	if got := len(g.Hand(0)); got != HandSize {
		t.Errorf("hand = %d, want %d (already full)", got, HandSize)
	}
	if g.State.Chains[0] != 3 {
		t.Errorf("chains = %d, want 3 (not blocked, none shed)", g.State.Chains[0])
	}
}

func TestChainsNotShedWhenNothingToDraw(t *testing.T) {
	// Below a full hand but with an empty deck and discard: the empty deck, not the
	// chains, is what stops the draw, so no chain is shed.
	g := NewGame("A", "B", 1)
	g.State.Chains[0] = 2

	g.EndPlayPhase(0)

	if got := len(g.Hand(0)); got != 0 {
		t.Errorf("hand = %d, want 0", got)
	}
	if g.State.Chains[0] != 2 {
		t.Errorf("chains = %d, want 2 (no cards to draw, none shed)", g.State.Chains[0])
	}
}

func TestChainsClampDrawAtZero(t *testing.T) {
	g := NewGame("A", "B", 1)
	for i := 0; i < HandSize; i++ {
		g.AddToDeck(testCreature("d", 1), 0)
	}
	g.State.Chains[0] = 40 // 40 chains: draw seven fewer, clamped to zero

	g.EndPlayPhase(0)

	if got := len(g.Hand(0)); got != 0 {
		t.Errorf("hand = %d, want 0 (reduction exceeds hand size)", got)
	}
	if g.State.Chains[0] != 39 {
		t.Errorf("chains = %d, want 39 (one shed)", g.State.Chains[0])
	}
}
