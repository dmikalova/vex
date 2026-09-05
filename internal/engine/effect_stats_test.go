package engine

import "testing"

func TestGainStats(t *testing.T) {
	g := NewGame("A", "B", 1)
	abond := g.AddToBattleline(testCreature("abond", 3), 0)
	friend := g.AddToBattleline(testCreature("friend", 5), 0)
	g.State.Cards[friend].ArmorRemaining = 0
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := GainStats{Target: Target{Kind: TargetEachOtherFriendlyCreature}, Armor: 1}
	if got := e.Text(); got != "for the remainder of the turn, each other friendly creature gains +1 armor" {
		t.Errorf("text = %q", got)
	}
	both := GainStats{Target: Target{Kind: TargetEachFriendlyCreature}, Power: 2, Armor: 2}
	if got := both.Text(); got != "for the remainder of the turn, each friendly creature gains +2 power and +2 armor" {
		t.Errorf("text = %q", got)
	}
	if (GainStats{Armor: 1}).validate() == nil {
		t.Error("unset target should be invalid")
	}
	if (GainStats{Target: Target{Kind: TargetEachFriendlyCreature}}).validate() == nil {
		t.Error("granting no stats should be invalid")
	}
	if e.validate() != nil {
		t.Error("a set target and a stat should be valid")
	}

	// The Action grants each other friendly creature +1 armor, topping up the
	// armor it has left to absorb damage this turn. Abond itself is excluded.
	e.Resolve(ctx)
	if got := g.Armor(friend); got != 1 {
		t.Errorf("friend armor = %d, want 1", got)
	}
	if got := g.State.Cards[friend].ArmorRemaining; got != 1 {
		t.Errorf("friend ArmorRemaining = %d, want 1", got)
	}
	if g.State.Cards[abond].TempArmorBonus != 0 {
		t.Error("the source creature should not gain armor from EachOtherFriendlyCreature")
	}

	// A power grant raises current power for the turn.
	both.Resolve(ctx)
	if got := g.Power(friend); got != 7 {
		t.Errorf("friend power = %d, want 7", got)
	}

	// The ready phase clears the bonus for every creature.
	g.StartTurn(0)
	g.EndPlayPhase(0)
	if g.State.Cards[friend].TempArmorBonus != 0 || g.State.Cards[friend].TempPowerBonus != 0 {
		t.Error("end of turn should clear the temporary stat bonus")
	}

	// Granting stats to a creature no longer in play is a no-op, not a panic.
	gone := g.AddToBattleline(testCreature("gone", 3), 0)
	g.DestroyEach(0, []LocalID{gone})
	g.GainStats(gone, 1, 1)
}
