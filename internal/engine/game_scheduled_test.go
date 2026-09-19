package engine

import "testing"

// TestScheduleOnLeaveText covers the arming node's rendering and the opponent
// ForgeKey it carries, plus the validation of the effect it can schedule.
func TestScheduleOnLeaveText(t *testing.T) {
	arm := ScheduleOnLeave{Do: ForgeKey{Player: Opponent, FreeOfCost: true}}
	if got := arm.Text(); got != "when {self} leaves play, your opponent forges a key at no cost" {
		t.Errorf("arm text = %q", got)
	}
	if got := (ForgeKey{Player: Opponent, FreeOfCost: true}).Text(); got != "your opponent forges a key at no cost" {
		t.Errorf("opponent forge text = %q", got)
	}
	if err := arm.validate(); err != nil {
		t.Errorf("valid arm should validate, got %v", err)
	}
	if action, ok := scheduledActionOf(ForgeKey{Player: Opponent, FreeOfCost: true}); !ok ||
		action != schedOpponentForgesKeyFree {
		t.Errorf("scheduledActionOf(opponent forge) = %d, %v", action, ok)
	}
}

// TestScheduleOnLeaveValidateRejects covers the two arming errors: no Do, and a Do
// the schedule cannot carry.
func TestScheduleOnLeaveValidateRejects(t *testing.T) {
	if err := (ScheduleOnLeave{}).validate(); err == nil {
		t.Error("a ScheduleOnLeave with no Do should fail validation")
	}
	if err := (ScheduleOnLeave{Do: GainAember{Player: Controller, Amount: 1}}).validate(); err == nil {
		t.Error("a ScheduleOnLeave with an unschedulable Do should fail validation")
	}
	if _, ok := scheduledActionOf(GainAember{Player: Controller, Amount: 1}); ok {
		t.Error("GainAember is not a schedulable action")
	}
}

// TestForgeKeyOpponentValidate covers the opponent forge only being supported at
// no cost.
func TestForgeKeyOpponentValidate(t *testing.T) {
	if err := (ForgeKey{Player: Opponent, FreeOfCost: true}).validate(); err != nil {
		t.Errorf("a free opponent forge should validate, got %v", err)
	}
	if err := (ForgeKey{Player: Opponent}).validate(); err == nil {
		t.Error("an opponent forge that is not free should fail validation")
	}
}

// TestScheduleOnLeaveForgesForOpponent covers the whole arc: arming the leave-play
// schedule, surviving the end-of-turn window, and forcing the opponent to forge a
// key when the source leaves play.
func TestScheduleOnLeaveForgesForOpponent(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Dis); err != nil {
		t.Fatal(err)
	}
	src := g.AddToBattleline(NewCard("Turnkey", Dis, Creature, Rare, WithPower(2)), 0)

	ScheduleOnLeave{Do: ForgeKey{Player: Opponent, FreeOfCost: true}}.Resolve(
		&EffectContext{Resolver: g, Source: src, Controller: 0},
	)
	if g.State.ScheduledCount != 1 {
		t.Fatalf("arming should schedule one effect, got %d", g.State.ScheduledCount)
	}
	if g.State.Scheduled[0].Duration != UntilThisLeavesPlay {
		t.Errorf("scheduled duration = %v, want UntilThisLeavesPlay", g.State.Scheduled[0].Duration)
	}
	if g.State.Scheduled[0].Do != schedOpponentForgesKeyFree {
		t.Errorf("scheduled action = %d, want schedOpponentForgesKeyFree", g.State.Scheduled[0].Do)
	}

	// The end-of-turn window neither fires nor clears a leave-play schedule.
	if pending := g.scheduledEndOfTurn(0); len(pending) != 0 {
		t.Errorf("a leave-play schedule should not be an end-of-turn entry, got %d", len(pending))
	}
	g.clearScheduled()
	if g.State.ScheduledCount != 1 {
		t.Fatal("clearScheduled should keep the leave-play schedule")
	}
	if g.Keys(1) != 0 {
		t.Fatalf("the opponent should not have forged yet, keys = %d", g.Keys(1))
	}

	// The source leaving play fires and removes the schedule.
	g.fireScheduledOnLeave(src)
	if g.Keys(1) != 1 {
		t.Errorf("the opponent should have forged one key on leave, keys = %d", g.Keys(1))
	}
	if g.State.ScheduledCount != 0 {
		t.Errorf("the fired schedule should be removed, count = %d", g.State.ScheduledCount)
	}
}

// TestFireScheduledOnLeaveSkipsAndCompacts covers fireScheduledOnLeave skipping the
// entries that do not match the leaving card, an end-of-turn schedule ignored by it,
// removeScheduledAt sliding later entries down, and a nil-effect entry firing
// harmlessly.
func TestFireScheduledOnLeaveSkipsAndCompacts(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Dis); err != nil {
		t.Fatal(err)
	}
	other := g.AddToBattleline(NewCard("Other", Dis, Creature, Common, WithPower(1)), 0)

	// An unset leave-play action for `other` (fires as a no-op), an end-of-turn
	// entry that must be left alone, and a real opponent forge for `other`.
	g.State.Scheduled[0] = ScheduledEffect{
		Source:   other,
		Do:       schedUnset,
		Duration: UntilThisLeavesPlay,
	}
	g.State.Scheduled[1] = ScheduledEffect{
		Source:   LocalID(99),
		Do:       schedDestroyEachCreature,
		Duration: RemainderOfPlayerTurn,
	}
	g.State.Scheduled[2] = ScheduledEffect{
		Source:   other,
		Do:       schedOpponentForgesKeyFree,
		Duration: UntilThisLeavesPlay,
	}
	g.State.ScheduledCount = 3

	g.fireScheduledOnLeave(other)

	if g.Keys(1) != 1 {
		t.Errorf(
			"the opponent should have forged from the matching entry, keys = %d",
			g.Keys(1),
		)
	}
	if g.State.ScheduledCount != 1 {
		t.Fatalf("only the end-of-turn entry should remain, count = %d", g.State.ScheduledCount)
	}
	if g.State.Scheduled[0].Do != schedDestroyEachCreature {
		t.Errorf(
			"surviving entry = %d, want the untouched end-of-turn schedule",
			g.State.Scheduled[0].Do,
		)
	}
}

// TestForgeKeyFreeForcedResolver covers the resolver wrapper and the opponent
// ForgeKey resolving through it.
func TestForgeKeyFreeForcedResolver(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Dis); err != nil {
		t.Fatal(err)
	}
	if !g.ForgeKeyFreeForced(1) {
		t.Fatal("forcing a free forge should report a key forged")
	}
	if g.Keys(1) != 1 {
		t.Errorf("forced forge keys = %d, want 1", g.Keys(1))
	}

	ForgeKey{Player: Opponent, FreeOfCost: true}.Resolve(
		&EffectContext{Resolver: g, Source: 0, Controller: 0},
	)
	if g.Keys(1) != 2 {
		t.Errorf("the opponent ForgeKey should forge again, keys = %d", g.Keys(1))
	}
}

// TestForgeKeyFreeForcedGuards covers the forced forge being stopped before it
// happens: a barred key ordinal, and a before-forge guard (Keyforgery) cancelling
// it on a wrong guess.
func TestForgeKeyFreeForcedGuards(t *testing.T) {
	barred := NewGame("A", "B", 1)
	barred.AddToBattleline(NewCard("Bronze Key Imp", Dis, Creature, Common,
		WithPower(2), WithRestrictions(Restrictions{NoForgeKeyNumber: 1})), 0)
	if barred.forgeKeyFreeForced(1) {
		t.Error("a barred first key should stop the forced forge")
	}
	if barred.Keys(1) != 0 {
		t.Errorf("a barred forced forge should forge nothing, keys = %d", barred.Keys(1))
	}

	guarded := NewGame("A", "B", 1)
	guarded.AddArtifact(keyforgeryCard(), 0) // guard is the forger's opponent
	guarded.AddToHand(NewCard("Logos Card", Logos, Creature, Common, WithPower(3)), 0)
	guarded.SetChooser(1, optionPicker{idx: 0}) // forger names Brobnar; the reveal is Logos
	if guarded.forgeKeyFreeForced(1) {
		t.Error("a before-forge guard should cancel the forced forge on a wrong guess")
	}
	if guarded.Keys(1) != 0 {
		t.Errorf("a cancelled forced forge should forge nothing, keys = %d", guarded.Keys(1))
	}
}

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
	for i := range maxScheduled + 2 {
		g.ScheduleAtEndOfTurn(LocalID(i), schedDestroyEachCreature)
	}
	if int(g.State.ScheduledCount) != maxScheduled {
		t.Errorf("scheduled count = %d, want capped at %d", g.State.ScheduledCount, maxScheduled)
	}
	if scheduledEffectOf(schedUnset) != nil {
		t.Error("scheduledEffectOf(schedUnset) should be nil")
	}
}
