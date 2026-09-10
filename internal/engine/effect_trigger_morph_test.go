package engine

import "testing"

func TestTriggerMorphValid(t *testing.T) {
	if !(TriggerMorph{From: TriggerAfterPlay, Onto: TriggerAfterReap}).valid() {
		t.Error("play->reap morph should be valid")
	}
	if (TriggerMorph{From: TriggerAfterPlay, Onto: TriggerDestroyed}).valid() {
		t.Error("a morph onto a non-action trigger should be invalid")
	}
	if (TriggerMorph{From: TriggerAction, Onto: TriggerAfterReap}).valid() {
		t.Error("a morph from a non-action trigger should be invalid")
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
	if g.State.MorphCount != 2 {
		t.Fatalf("morph count = %d, want 2", g.State.MorphCount)
	}
	// Reaping fires fight abilities; fighting fires reap abilities.
	fromReap := g.State.Morphs[0].From == TriggerAfterFight &&
		g.State.Morphs[0].Onto == TriggerAfterReap
	fromFight := g.State.Morphs[1].From == TriggerAfterReap &&
		g.State.Morphs[1].Onto == TriggerAfterFight
	if !fromReap || !fromFight {
		t.Errorf("morphs = %+v, want fight->reap and reap->fight", g.State.Morphs[:2])
	}
}

func TestAddLastingMorphCaps(t *testing.T) {
	g := started(t)
	for i := 0; i < maxMorph+2; i++ {
		g.AddLastingMorph(
			LastingMorph{Controller: 0, From: TriggerAfterFight, Onto: TriggerAfterReap},
		)
	}
	if int(g.State.MorphCount) != maxMorph {
		t.Errorf("morph count = %d, want %d (capped)", g.State.MorphCount, maxMorph)
	}
}

func TestClearLastingMorphsKeepsOtherPlayer(t *testing.T) {
	g := started(t)
	g.AddLastingMorph(LastingMorph{Controller: 0, From: TriggerAfterFight, Onto: TriggerAfterReap})
	g.AddLastingMorph(LastingMorph{Controller: 1, From: TriggerAfterReap, Onto: TriggerAfterFight})
	g.clearLastingMorphs(0)
	if g.State.MorphCount != 1 {
		t.Fatalf("morph count = %d, want 1 (the opponent's kept)", g.State.MorphCount)
	}
	if g.State.Morphs[0].Controller != 1 {
		t.Errorf("remaining morph controller = %d, want 1", g.State.Morphs[0].Controller)
	}
	// clearLasting also drops a player's morphs.
	g.clearLasting(1)
	if g.State.MorphCount != 0 {
		t.Errorf("morph count = %d, want 0 after clearing the opponent's turn", g.State.MorphCount)
	}
}

// morphConstant is a creature that, while in play, makes each friendly creature's
// play effect also fire on reap — the Kompsos Haruspex mechanic in miniature.
func morphConstant() CardDefinition {
	return NewCard("Morpher", Brobnar, Creature, Rare, WithPower(4),
		WithConstantAbility(ConstantAbility{
			Target: Target{Kind: TargetEachFriendlyCreature},
			Morphs: []TriggerMorph{{From: TriggerAfterPlay, Onto: TriggerAfterReap}},
		}))
}

func TestMorphedTriggersFromConstant(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(morphConstant(), 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)

	got := g.morphedTriggers(friend, TriggerAfterReap)
	if len(got) != 1 || got[0] != TriggerAfterPlay {
		t.Errorf("morphed triggers on reap = %v, want [play]", got)
	}
	// The constant reaches its own controller's creature, including the source.
	if got := g.morphedTriggers(src, TriggerAfterReap); len(got) != 1 {
		t.Errorf("source should also be reached, got %v", got)
	}
	// No morph maps onto fight, so fighting gathers nothing extra.
	if got := g.morphedTriggers(friend, TriggerAfterFight); got != nil {
		t.Errorf("fight morphs = %v, want none", got)
	}
	// The constant does not reach the opponent's creatures.
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	if got := g.morphedTriggers(foe, TriggerAfterReap); got != nil {
		t.Errorf("enemy morphs = %v, want none", got)
	}
}

func TestMorphedTriggersFromLastingAndDedup(t *testing.T) {
	g := started(t)
	// A constant morph and a lasting morph both map play onto reap; the result is
	// deduplicated to a single trigger.
	g.AddToBattleline(morphConstant(), 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)
	g.AddLastingMorph(LastingMorph{Controller: 0, From: TriggerAfterPlay, Onto: TriggerAfterReap})
	if got := g.morphedTriggers(friend, TriggerAfterReap); len(got) != 1 {
		t.Errorf("duplicate morphs should collapse, got %v", got)
	}
	// A lasting morph is scoped to its owner.
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	g.AddLastingMorph(LastingMorph{Controller: 0, From: TriggerAfterFight, Onto: TriggerAfterReap})
	if got := g.morphedTriggers(foe, TriggerAfterReap); got != nil {
		t.Errorf("enemy lasting morphs = %v, want none", got)
	}
}

func TestConstantMorphFiresPlayAbilityOnReap(t *testing.T) {
	g := started(t)
	g.AddToBattleline(morphConstant(), 0)
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
	// +1 from the reap itself, +1 from the play ability fired by the morph.
	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2 (reap + morphed play)", g.Aember(0))
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
		t.Errorf("aember after reap = %d, want 2 (reap + morphed fight)", g.Aember(0))
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
		t.Errorf("aember after fight = %d, want 3 (morphed reap)", g.Aember(0))
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

func TestConstantMorphText(t *testing.T) {
	def := morphConstant()
	got := constantText(&def)
	want := "Each friendly creature's play effect is a play/reap effect."
	if got != want {
		t.Errorf("constant text = %q, want %q", got, want)
	}
}

func TestConstantMorphInvalidTriggerPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("a constant morph naming a non-action trigger should panic at init")
		}
	}()
	NewCard("Bad Morpher", Brobnar, Creature, Common, WithPower(1),
		WithConstantAbility(ConstantAbility{
			Morphs: []TriggerMorph{{From: TriggerAfterPlay, Onto: TriggerDestroyed}},
		}))
}
