package engine

import "testing"

// TestGainAssaultText covers the printed clause, standalone and folded under a
// shared duration with a keyword grant.
func TestGainAssaultText(t *testing.T) {
	e := GainAssault{Target: Target{Kind: TargetTriggeringCreature}, Amount: PowerOfChosen{}}
	if got := e.Text(); got != "for the remainder of the turn, it gains assault equal to its power" {
		t.Errorf("text = %q", got)
	}
	folded := ForDuration{
		Duration: RemainderOfPlayerTurn,
		Effects: []Effect{
			GainKeywords{
				Target:   Target{Kind: TargetTriggeringCreature},
				Keywords: []Keyword{Skirmish},
				Duration: RemainderOfPlayerTurn,
			},
			GainAssault{Target: Target{Kind: TargetTriggeringCreature}, Amount: PowerOfChosen{}},
		},
	}
	want := "for the remainder of the turn, it gains skirmish and assault equal to its power"
	if got := folded.Text(); got != want {
		t.Errorf("folded text = %q, want %q", got, want)
	}

	// A Fixed amount renders a plain number, not "equal to".
	fixed := GainAssault{Target: Target{Kind: TargetThisCreature}, Amount: Fixed(2)}
	if got := fixed.Text(); got != "for the remainder of the turn, "+SelfName+" gains assault 2" {
		t.Errorf("fixed text = %q", got)
	}
}

// TestGainAssaultValidate rejects a missing target and a nil amount.
func TestGainAssaultValidate(t *testing.T) {
	if (GainAssault{Amount: PowerOfChosen{}}).validate() == nil {
		t.Error("unset target should be invalid")
	}
	if (GainAssault{Target: Target{Kind: TargetTriggeringCreature}}).validate() == nil {
		t.Error("a nil amount should be invalid")
	}
	if (GainAssault{
		Target: Target{Kind: TargetTriggeringCreature},
		Amount: Fixed(2),
	}).validate() != nil {
		t.Error("a set target and amount should be valid")
	}
}

// TestGainAssaultResolve grants Assault equal to the chosen creature's power for
// the turn, folds into the creature's Assault value, and lifts at the ready phase.
func TestGainAssaultResolve(t *testing.T) {
	g := started(t)
	g.SetRecording(true)
	beast := g.AddToBattleline(testCreature("beast", 4), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0, It: beast, HasIt: true}

	GainAssault{
		Target: Target{Kind: TargetTriggeringCreature},
		Amount: PowerOfChosen{},
	}.Resolve(
		ctx,
	)
	if got := g.assault(beast); got != 4 {
		t.Errorf("assault = %d, want 4 (equal to its power)", got)
	}
	if !hasLogLine(g, "beast gains assault 4") {
		t.Error("the grant should be logged")
	}

	// The ready phase clears the bonus for every creature.
	g.StartTurn(0)
	g.EndPlayPhase(0)
	if g.State.Cards[beast].TempAssaultBonus != 0 {
		t.Error("end of turn should clear the temporary assault bonus")
	}
}

// TestGainAssaultZeroAndGone covers the no-op branches: a non-positive amount
// grants nothing, and granting to a creature no longer in play does not panic.
func TestGainAssaultZeroAndGone(t *testing.T) {
	g := started(t)
	beast := g.AddToBattleline(testCreature("beast", 4), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0, It: beast, HasIt: true}

	GainAssault{Target: Target{Kind: TargetTriggeringCreature}, Amount: Fixed(0)}.Resolve(ctx)
	if g.State.Cards[beast].TempAssaultBonus != 0 {
		t.Error("a non-positive amount should grant nothing")
	}

	gone := g.AddToBattleline(testCreature("gone", 3), 0)
	g.DestroyEach(0, []LocalID{gone})
	g.GainAssault(gone, 2)
}

// TestGainAssaultInCombat covers Assault gained for the turn dealing its damage in
// a fight.
func TestGainAssaultInCombat(t *testing.T) {
	g := started(t)
	attacker := g.AddToBattleline(testCreature("attacker", 3), 0)
	defender := g.AddToBattleline(testCreature("defender", 6), 1)
	g.GainAssault(attacker, 3)

	g.FightWith(attacker, defender)
	// 3 Assault + 3 fight damage = 6 on the 6-power defender: destroyed.
	if g.inPlay(defender) {
		t.Error("the defender should be destroyed by assault plus fight damage")
	}
}

// TestChooseCreatureGainsSkirmishAndAssault covers the composed grant Creed of
// Nature drives: choose a creature, then for the remainder of the turn it gains
// skirmish and assault equal to its power. The chosen creature then fights, dealing
// its assault damage and taking no return damage from skirmish.
func TestChooseCreatureGainsSkirmishAndAssault(t *testing.T) {
	g := started(t)
	chosen := g.AddToBattleline(testCreature("chosen", 4), 0)
	defender := g.AddToBattleline(testCreature("defender", 5), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	grant := ChooseCreatureThen{
		Target: Target{Kind: TargetChosenCreature},
		Then: ForDuration{
			Duration: RemainderOfPlayerTurn,
			Effects: []Effect{
				GainKeywords{
					Target:   Target{Kind: TargetTriggeringCreature},
					Keywords: []Keyword{Skirmish},
					Duration: RemainderOfPlayerTurn,
				},
				GainAssault{
					Target: Target{Kind: TargetTriggeringCreature},
					Amount: PowerOfChosen{},
				},
			},
		},
	}
	want := "choose a creature - for the remainder of the turn, it gains skirmish and assault equal to its power"
	if got := grant.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}

	g.SetChooser(0, &idQueueChooser{ids: []LocalID{chosen}})
	grant.Resolve(ctx)
	if !g.hasKeyword(chosen, Skirmish) {
		t.Error("the chosen creature should have gained skirmish")
	}
	if got := g.assault(chosen); got != 4 {
		t.Errorf("assault = %d, want 4 (equal to its power)", got)
	}

	// 4 Assault + 4 fight damage = 8 destroys the 5-power defender; skirmish spares
	// the 4-power attacker the 5 return damage, so it survives.
	g.FightWith(chosen, defender)
	if g.inPlay(defender) {
		t.Error("the defender should be destroyed by assault plus fight damage")
	}
	if !g.inPlay(chosen) {
		t.Error("skirmish should spare the attacker its return damage")
	}
}

// TestGainAssaultUntilNextTurnText covers the printed clause: a Fixed amount reads
// as a plain number, a scaling one reads "equal to".
func TestGainAssaultUntilNextTurnText(t *testing.T) {
	fixed := GainAssaultUntilNextTurn{
		Target: Target{Kind: TargetTriggeringCreature},
		Amount: Fixed(3),
	}
	if got := fixed.Text(); got != "it gains assault 3 until the start of your next turn" {
		t.Errorf("fixed text = %q", got)
	}
	scaling := GainAssaultUntilNextTurn{
		Target: Target{Kind: TargetTriggeringCreature},
		Amount: PowerOfChosen{},
	}
	want := "it gains assault equal to its power until the start of your next turn"
	if got := scaling.Text(); got != want {
		t.Errorf("scaling text = %q, want %q", got, want)
	}
}

// TestGainAssaultUntilNextTurnValidate rejects a missing target and a nil amount.
func TestGainAssaultUntilNextTurnValidate(t *testing.T) {
	if (GainAssaultUntilNextTurn{Amount: Fixed(3)}).validate() == nil {
		t.Error("unset target should be invalid")
	}
	if (GainAssaultUntilNextTurn{Target: Target{Kind: TargetTriggeringCreature}}).validate() == nil {
		t.Error("a nil amount should be invalid")
	}
	if (GainAssaultUntilNextTurn{
		Target: Target{Kind: TargetTriggeringCreature},
		Amount: Fixed(3),
	}).validate() != nil {
		t.Error("a set target and amount should be valid")
	}
}

// TestGainAssaultUntilNextTurnResolve grants Assault that folds into the creature's
// Assault value and, unlike GainAssault, survives the opponent's turn and lifts only
// at the start of the controller's own next turn. A non-positive amount grants
// nothing, and granting to a creature no longer in play does not panic.
func TestGainAssaultUntilNextTurnResolve(t *testing.T) {
	g := started(t)
	g.SetRecording(true)
	beast := g.AddToBattleline(testCreature("beast", 4), 0)

	GainAssaultUntilNextTurn{Target: Target{Kind: TargetTriggeringCreature}, Amount: Fixed(3)}.
		Resolve(&EffectContext{Resolver: g, Controller: 0, It: beast, HasIt: true})
	if got := g.assault(beast); got != 3 {
		t.Errorf("assault = %d, want 3", got)
	}
	if !hasLogLine(g, "beast gains assault 3") {
		t.Error("the grant should be logged")
	}

	// It survives both ready phases and the opponent's start-of-turn, and lifts at
	// the start of the controller's own next turn.
	g.readyPhase(0)
	g.startOfTurnPhase(1)
	if g.assault(beast) != 3 {
		t.Error("the grant should survive the opponent's turn")
	}
	g.startOfTurnPhase(0)
	if g.assault(beast) != 0 {
		t.Error("the grant should lift at the start of the controller's next turn")
	}

	// A non-positive amount grants nothing.
	GainAssaultUntilNextTurn{Target: Target{Kind: TargetTriggeringCreature}, Amount: Fixed(0)}.
		Resolve(&EffectContext{Resolver: g, Controller: 0, It: beast, HasIt: true})
	if g.assault(beast) != 0 {
		t.Error("a non-positive amount should grant nothing")
	}

	// Granting to a creature no longer in play is a safe no-op.
	gone := g.AddToBattleline(testCreature("gone", 3), 0)
	g.DestroyEach(0, []LocalID{gone})
	g.GrantAssaultUntilNextTurn(gone, 2)
}
