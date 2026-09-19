package engine

import "testing"

func TestItIsEnemy(t *testing.T) {
	g := NewGame("A", "B", 1)
	mine := g.AddToBattleline(testCreature("mine", 3), 0)
	theirs := g.AddToBattleline(testCreature("theirs", 3), 1)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	cond := ItIsEnemy{}
	if got := cond.CondText(); got != "if it is an enemy creature" {
		t.Errorf("CondText = %q", got)
	}

	// No card in context: not met.
	if cond.Met(ctx) {
		t.Error("ItIsEnemy should not be met with no context card")
	}
	ctx.It, ctx.HasIt = theirs, true
	if !cond.Met(ctx) {
		t.Error("a card the opponent controls should meet ItIsEnemy")
	}
	ctx.It = mine
	if cond.Met(ctx) {
		t.Error("a card the controller controls should not meet ItIsEnemy")
	}
}
