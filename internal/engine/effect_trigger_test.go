package engine

import "testing"

func TestTriggerAbility(t *testing.T) {
	e := TriggerAbility{
		Trigger: TriggerAfterReap,
		Target:  Target{Kind: TargetChosenCreature}.Other(),
	}
	want := "trigger the reap effect of another creature"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
	if (TriggerAbility{Trigger: TriggerAfterReap}).validate() == nil {
		t.Error("unset target should be invalid")
	}
	if (TriggerAbility{Trigger: TriggerAfterForgeKey, Target: Target{Kind: TargetChosenCreature}}).
		validate() == nil {
		t.Error("a trigger with no effect noun should be invalid")
	}
	if e.validate() != nil {
		t.Error("play/fight/reap triggers should be valid")
	}
	for _, tr := range []Trigger{TriggerAfterPlay, TriggerAfterFight} {
		if triggerEffectNoun(tr) == "" {
			t.Errorf("trigger %v should name an effect", tr)
		}
	}

	// The source reaps and fires the only other creature carrying a Reap ability.
	g := started(t)
	gainer := testCreature("Gainer", 2, WithAbility(
		TriggerAfterReap, GainAember{Amount: 1, Player: Controller}))
	g.AddToBattleline(testCreature("Source", 2), 0)
	g.AddToBattleline(gainer, 0)
	g.AddToBattleline(testCreature("Bystander", 2), 0)
	src := g.Battleline(0)[0]
	other := g.Battleline(0)[1]

	e.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: src})
	if g.Aember(0) != 1 {
		t.Errorf("Æmber = %d, want 1", g.Aember(0))
	}
	if g.Exhausted(other) {
		t.Error("the triggered creature should not exhaust")
	}
}

func TestTriggerAbilityBoundedByRuleOfSix(t *testing.T) {
	// Two Replicator-like creatures reach for each other's reap effect; the chain
	// bounces until the Rule of Six stops it. Each pass gains 1 Æmber, so the pool
	// counts the resolutions.
	replicate := Sequence{Effects: []Effect{
		GainAember{Amount: 1, Player: Controller},
		TriggerAbility{
			Trigger: TriggerAfterReap,
			Target:  Target{Kind: TargetChosenCreature}.Other(),
		},
	}}
	g := started(t)
	def := testCreature("Replicant", 2, WithAbility(TriggerAfterReap, replicate))
	g.AddToBattleline(def, 0)
	g.AddToBattleline(def, 0)
	src := g.Battleline(0)[0]

	replicate.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: src})

	if g.Aember(0) != RuleOfSix+1 {
		t.Errorf("Æmber = %d, want %d (the first pass plus %d bounces)",
			g.Aember(0), RuleOfSix+1, RuleOfSix)
	}
	if g.TriggerDepth() != 0 {
		t.Errorf("trigger depth = %d, want 0 once resolution ends", g.TriggerDepth())
	}
}

func TestTriggerAbilityNoCandidate(t *testing.T) {
	g := started(t)
	g.AddToBattleline(testCreature("Lonely", 2), 0)
	src := g.Battleline(0)[0]
	TriggerAbility{
		Trigger: TriggerAfterReap,
		Target:  Target{Kind: TargetChosenCreature}.Other(),
	}.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: src})
}
