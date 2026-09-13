package engine

import "testing"

// TestDestroyEachCreatureAtEndOfTurn covers Ragnarok's scheduled board wipe: the
// effect schedules the wipe during the play phase, the schedule survives the ready
// phase, and the end-of-turn window destroys every creature and clears the schedule.
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
	if g.State.ScheduledCount != 1 {
		t.Fatalf("Resolve should schedule one effect, got count %d", g.State.ScheduledCount)
	}
	if g.State.Scheduled[0].Source != src {
		t.Errorf("scheduled source = %d, want %d", g.State.Scheduled[0].Source, src)
	}
	if g.State.Scheduled[0].Do != schedDestroyEachCreature {
		t.Errorf("scheduled action = %d, want schedDestroyEachCreature", g.State.Scheduled[0].Do)
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
	if g.State.ScheduledCount != 0 {
		t.Error("the schedule should be cleared once its window fires")
	}
	_, _ = mine, theirs
}

// TestScheduleAtEndOfTurnFull covers the schedule silently dropping entries once
// the fixed array is full, and the scheduledEffectOf fallback for an unset action.
func TestScheduleAtEndOfTurnFull(t *testing.T) {
	g := NewGame("A", "B", 1)
	for i := 0; i < maxScheduled+2; i++ {
		g.ScheduleAtEndOfTurn(LocalID(i), schedDestroyEachCreature)
	}
	if int(g.State.ScheduledCount) != maxScheduled {
		t.Errorf("scheduled count = %d, want capped at %d", g.State.ScheduledCount, maxScheduled)
	}
	if scheduledEffectOf(schedUnset) != nil {
		t.Error("scheduledEffectOf(schedUnset) should be nil")
	}
}
