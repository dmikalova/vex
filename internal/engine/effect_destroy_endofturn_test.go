package engine

import "testing"

// TestDestroyEachCreatureAtEndOfTurn covers Ragnarok's scheduled board wipe: the
// effect arms a flag during the play phase, the flag survives the ready phase, and
// the end-of-turn phase destroys every creature and clears the flag.
func TestDestroyEachCreatureAtEndOfTurn(t *testing.T) {
	if got := (DestroyEachCreatureAtEndOfTurn{}).Text(); got != "at the end of the turn, destroy each creature" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	src := g.AddToDiscard(NewCard("Ragnarok", Brobnar, Tactic, Rare), 0)
	mine := g.AddToBattleline(NewCard("mine", Brobnar, Creature, Common, WithPower(3)), 0)
	theirs := g.AddToBattleline(NewCard("theirs", Brobnar, Creature, Common, WithPower(3)), 1)

	DestroyEachCreatureAtEndOfTurn{}.Resolve(
		&EffectContext{Resolver: g, Source: src, Controller: 0},
	)
	if !g.State.EndOfTurnDestroyAll.Value {
		t.Fatal("Resolve should arm the end-of-turn wipe")
	}
	if g.State.EndOfTurnDestroyAll.Source != src {
		t.Errorf("armed source = %d, want %d", g.State.EndOfTurnDestroyAll.Source, src)
	}
	// The wipe is scheduled, not immediate: both creatures are still in play.
	if len(g.Battleline(0)) != 1 || len(g.Battleline(1)) != 1 {
		t.Fatalf("battlelines = %v / %v, want both still present", g.Battleline(0), g.Battleline(1))
	}

	g.EndPlayPhase(0)

	if len(g.Battleline(0)) != 0 || len(g.Battleline(1)) != 0 {
		t.Errorf(
			"battlelines after end of turn = %v / %v, want both empty",
			g.Battleline(0),
			g.Battleline(1),
		)
	}
	if g.State.EndOfTurnDestroyAll.Value {
		t.Error("the wipe flag should be cleared once it fires")
	}
	_, _ = mine, theirs
}
