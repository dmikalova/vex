package engine

import "testing"

func TestOpponentExcessCreatures(t *testing.T) {
	if (OpponentExcessCreatures{}).CountText() != "creature your opponent controls in excess of you" {
		t.Errorf("count text = %q", (OpponentExcessCreatures{}).CountText())
	}

	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("o1", 3), 1)
	g.AddToBattleline(testCreature("o2", 3), 1)
	g.AddToBattleline(testCreature("m1", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if got := (OpponentExcessCreatures{}).Value(ctx); got != 1 {
		t.Errorf("excess = %d, want 1 (opponent 2, you 1)", got)
	}

	// When the controller controls at least as many, the excess floors at zero.
	g.AddToBattleline(testCreature("m2", 3), 0)
	g.AddToBattleline(testCreature("m3", 3), 0)
	if got := (OpponentExcessCreatures{}).Value(ctx); got != 0 {
		t.Errorf("excess = %d, want 0 when you control more", got)
	}
}

func TestMayUseFriendlyHouse(t *testing.T) {
	if got := (MayUseFriendlyHouse{House: Sanctum}).Text(); got != "for the remainder of the turn, you may use friendly Sanctum creatures" {
		t.Errorf("text = %q", got)
	}
	if (MayUseFriendlyHouse{}).validate() == nil {
		t.Error("unset house should be invalid")
	}
	if (MayUseFriendlyHouse{House: Sanctum}).validate() != nil {
		t.Error("a set house should be valid")
	}

	g := NewGame("A", "B", 1)
	g.BeginTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	c := g.AddToBattleline(NewCard("cleric", Sanctum, Creature, Common, WithPower(3)), 0)
	if g.usableInActiveHouse(c) {
		t.Fatal("an off-house creature should not be usable before the grant")
	}

	MayUseFriendlyHouse{House: Sanctum}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.MayUseHouse[0] != Sanctum {
		t.Fatal("the grant should record the house")
	}
	if !g.usableInActiveHouse(c) {
		t.Error("the granted-house creature should be usable")
	}

	g.EndTurn(0)
	if g.State.MayUseHouse[0] != HouseNone {
		t.Error("the grant should clear at end of turn")
	}
}
