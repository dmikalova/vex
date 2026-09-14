package engine

import "testing"

func TestAlsoTriggersOnValid(t *testing.T) {
	if !(AlsoTriggersOn{From: TriggerAfterPlay, Onto: TriggerAfterReap}).valid() {
		t.Error("play->reap rule should be valid")
	}
	if (AlsoTriggersOn{From: TriggerAfterPlay, Onto: TriggerDestroyed}).valid() {
		t.Error("a rule onto a non-action trigger should be invalid")
	}
	if (AlsoTriggersOn{From: TriggerAction, Onto: TriggerAfterReap}).valid() {
		t.Error("a rule from a non-action trigger should be invalid")
	}
}

func TestFuseTriggersForTurnValidate(t *testing.T) {
	if err := (FuseTriggersForTurn{A: TriggerAfterFight, B: TriggerAfterReap}).validate(); err != nil {
		t.Errorf("fight/reap fuse should validate, got %v", err)
	}
	if err := (FuseTriggersForTurn{A: TriggerAfterFight, B: TriggerAfterFight}).validate(); err == nil {
		t.Error("a fuse of a trigger with itself should be rejected")
	}
	if err := (FuseTriggersForTurn{A: TriggerAction, B: TriggerAfterReap}).validate(); err == nil {
		t.Error("a fuse naming a non-action trigger should be rejected")
	}
}

func TestFuseTriggersForTurnText(t *testing.T) {
	got := (FuseTriggersForTurn{A: TriggerAfterFight, B: TriggerAfterReap}).Text()
	want := "each friendly creature's fight effects and reap effects are " +
		"fight/reap effects for the remainder of the turn"
	if got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

func TestFuseTriggersForTurnResolveInstallsBothDirections(t *testing.T) {
	g := started(t)
	(FuseTriggersForTurn{A: TriggerAfterFight, B: TriggerAfterReap}).
		Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.AlsoTriggersCount != 2 {
		t.Fatalf("rule count = %d, want 2", g.State.AlsoTriggersCount)
	}
	// Reaping fires fight abilities; fighting fires reap abilities.
	fromReap := g.State.AlsoTriggers[0].From == TriggerAfterFight &&
		g.State.AlsoTriggers[0].Onto == TriggerAfterReap
	fromFight := g.State.AlsoTriggers[1].From == TriggerAfterReap &&
		g.State.AlsoTriggers[1].Onto == TriggerAfterFight
	if !fromReap || !fromFight {
		t.Errorf("rules = %+v, want fight->reap and reap->fight", g.State.AlsoTriggers[:2])
	}
}

func TestAddLastingAlsoTriggersCaps(t *testing.T) {
	g := started(t)
	for i := 0; i < maxAlsoTriggers+2; i++ {
		g.AddLastingAlsoTriggers(
			LastingAlsoTriggersOn{Controller: 0, From: TriggerAfterFight, Onto: TriggerAfterReap},
		)
	}
	if int(g.State.AlsoTriggersCount) != maxAlsoTriggers {
		t.Errorf("rule count = %d, want %d (capped)", g.State.AlsoTriggersCount, maxAlsoTriggers)
	}
}

func TestClearLastingAlsoTriggersKeepsOtherPlayer(t *testing.T) {
	g := started(t)
	g.AddLastingAlsoTriggers(
		LastingAlsoTriggersOn{Controller: 0, From: TriggerAfterFight, Onto: TriggerAfterReap},
	)
	g.AddLastingAlsoTriggers(
		LastingAlsoTriggersOn{Controller: 1, From: TriggerAfterReap, Onto: TriggerAfterFight},
	)
	g.clearLastingAlsoTriggers(0)
	if g.State.AlsoTriggersCount != 1 {
		t.Fatalf("rule count = %d, want 1 (the opponent's kept)", g.State.AlsoTriggersCount)
	}
	if g.State.AlsoTriggers[0].Controller != 1 {
		t.Errorf("remaining rule controller = %d, want 1", g.State.AlsoTriggers[0].Controller)
	}
	// clearLasting also drops a player's rules.
	g.clearLasting(1)
	if g.State.AlsoTriggersCount != 0 {
		t.Errorf(
			"rule count = %d, want 0 after clearing the opponent's turn",
			g.State.AlsoTriggersCount,
		)
	}
}

// alsoTriggersConstant is a creature that, while in play, makes each friendly
// creature's play effect also fire on reap — the Kompsos Haruspex mechanic in
// miniature.
func alsoTriggersConstant() CardDefinition {
	return NewCard("Echoer", Brobnar, Creature, Rare, WithPower(4),
		WithConstantAbility(ConstantAbility{
			Target:       Target{Kind: TargetEachFriendlyCreature},
			AlsoTriggers: []AlsoTriggersOn{{From: TriggerAfterPlay, Onto: TriggerAfterReap}},
		}))
}

func TestAdditionalTriggersFromConstant(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(alsoTriggersConstant(), 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)

	got := g.additionalTriggers(friend, TriggerAfterReap)
	if len(got) != 1 || got[0] != TriggerAfterPlay {
		t.Errorf("additional triggers on reap = %v, want [play]", got)
	}
	// The constant reaches its own controller's creature, including the source.
	if got := g.additionalTriggers(src, TriggerAfterReap); len(got) != 1 {
		t.Errorf("source should also be reached, got %v", got)
	}
	// No rule maps onto fight, so fighting gathers nothing extra.
	if got := g.additionalTriggers(friend, TriggerAfterFight); got != nil {
		t.Errorf("fight triggers = %v, want none", got)
	}
	// The constant does not reach the opponent's creatures.
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	if got := g.additionalTriggers(foe, TriggerAfterReap); got != nil {
		t.Errorf("enemy triggers = %v, want none", got)
	}
}

func TestAdditionalTriggersFromLastingAndDedup(t *testing.T) {
	g := started(t)
	// A constant rule and a lasting rule both map play onto reap; the result is
	// deduplicated to a single trigger.
	g.AddToBattleline(alsoTriggersConstant(), 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)
	g.AddLastingAlsoTriggers(
		LastingAlsoTriggersOn{Controller: 0, From: TriggerAfterPlay, Onto: TriggerAfterReap},
	)
	if got := g.additionalTriggers(friend, TriggerAfterReap); len(got) != 1 {
		t.Errorf("duplicate triggers should collapse, got %v", got)
	}
	// A lasting rule is scoped to its owner.
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	g.AddLastingAlsoTriggers(
		LastingAlsoTriggersOn{Controller: 0, From: TriggerAfterFight, Onto: TriggerAfterReap},
	)
	if got := g.additionalTriggers(foe, TriggerAfterReap); got != nil {
		t.Errorf("enemy lasting triggers = %v, want none", got)
	}
}

func TestConstantAlsoTriggersFiresPlayAbilityOnReap(t *testing.T) {
	g := started(t)
	g.AddToBattleline(alsoTriggersConstant(), 0)
	friend := g.AddToBattleline(
		testCreature(
			"friend",
			3,
			WithAbility(TriggerAfterPlay, GainAember{Player: Controller, Amount: 1}),
		),
		0,
	)
	if err := g.Reap(0, friend); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	// +1 from the reap itself, +1 from the play ability fired by the rule.
	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2 (reap + echoed play)", g.Aember(0))
	}
}

func TestLastingFuseFiresBothWays(t *testing.T) {
	g := started(t)
	(FuseTriggersForTurn{A: TriggerAfterFight, B: TriggerAfterReap}).
		Resolve(&EffectContext{Resolver: g, Controller: 0})

	// A creature with only a Fight ability fires it on reap.
	reaper := g.AddToBattleline(
		testCreature(
			"reaper",
			3,
			WithAbility(TriggerAfterFight, GainAember{Player: Controller, Amount: 1}),
		),
		0,
	)
	if err := g.Reap(0, reaper); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if g.Aember(0) != 2 {
		t.Errorf("aember after reap = %d, want 2 (reap + echoed fight)", g.Aember(0))
	}

	// A creature with only a Reap ability fires it on fight.
	fighter := g.AddToBattleline(
		testCreature(
			"fighter",
			8,
			WithAbility(TriggerAfterReap, GainAember{Player: Controller, Amount: 1}),
		),
		0,
	)
	foe := g.AddToBattleline(testCreature("foe", 1), 1)
	if err := g.Fight(0, fighter, foe); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.Aember(0) != 3 {
		t.Errorf("aember after fight = %d, want 3 (echoed reap)", g.Aember(0))
	}

	// The fuse lasts only the turn: once cleared, reaping fires nothing extra.
	g.clearLasting(0)
	reaper2 := g.AddToBattleline(
		testCreature(
			"reaper2",
			3,
			WithAbility(TriggerAfterFight, GainAember{Player: Controller, Amount: 1}),
		),
		0,
	)
	before := g.Aember(0)
	if err := g.Reap(0, reaper2); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if g.Aember(0) != before+1 {
		t.Errorf("aember after cleared reap = %d, want %d (reap only)", g.Aember(0), before+1)
	}
}

func TestConstantAlsoTriggersText(t *testing.T) {
	def := alsoTriggersConstant()
	got := constantText(&def)
	want := "Each friendly creature's play effect is a play/reap effect."
	if got != want {
		t.Errorf("constant text = %q, want %q", got, want)
	}
}

func TestConstantAlsoTriggersInvalidTriggerPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("a constant rule naming a non-action trigger should panic at init")
		}
	}()
	NewCard("Bad Echoer", Brobnar, Creature, Common, WithPower(1),
		WithConstantAbility(ConstantAbility{
			AlsoTriggers: []AlsoTriggersOn{{From: TriggerAfterPlay, Onto: TriggerDestroyed}},
		}))
}
