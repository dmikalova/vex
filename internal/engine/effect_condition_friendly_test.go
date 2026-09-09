package engine

import "testing"

func TestItIsFriendly(t *testing.T) {
	g := NewGame("A", "B", 1)
	mine := g.AddToBattleline(testCreature("mine", 3), 0)
	theirs := g.AddToBattleline(testCreature("theirs", 3), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	cond := ItIsFriendly{}
	if got := cond.CondText(); got != "if it is a friendly creature" {
		t.Errorf("CondText = %q", got)
	}

	// No card in context: not met.
	if cond.Met(ctx) {
		t.Error("ItIsFriendly should not be met with no context card")
	}
	ctx.It, ctx.HasIt = mine, true
	if !cond.Met(ctx) {
		t.Error("a card the controller controls should meet ItIsFriendly")
	}
	ctx.It = theirs
	if cond.Met(ctx) {
		t.Error("a card the opponent controls should not meet ItIsFriendly")
	}
}
