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
	if (TriggerAbility{
		Trigger: TriggerAfterForgeKey,
		Target:  Target{Kind: TargetChosenCreature},
	}).
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
		TriggerAfterReap, GainAember{
			Amount: 1,
			Player: Controller,
		}))
	g.AddToBattleline(testCreature("Source", 2), 0)
	g.AddToBattleline(gainer, 0)
	g.AddToBattleline(testCreature("Bystander", 2), 0)
	src := g.Battleline(0)[0]
	other := g.Battleline(0)[1]

	e.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     src,
	})
	if g.Aember(0) != 1 {
		t.Errorf("Æmber = %d, want 1", g.Aember(0))
	}
	if g.Exhausted(other) {
		t.Error("the triggered creature should not exhaust")
	}
	if got := g.nameUsagesThisTurn(src); got != 0 {
		t.Errorf("the use's first trigger should ride free, pool = %d, want 0", got)
	}
}

func TestTriggerAbilityBoundedByRuleOfSix(t *testing.T) {
	// Two Replicator-like creatures reach for each other's reap effect. The use
	// fires the first pass free; each further trigger spends one from the pool of
	// six the two share, and the pass whose trigger is finally blocked still gains
	// its Æmber — so the chain yields RuleOfSix+2 gains.
	replicate := Sequence{Effects: []Effect{
		GainAember{
			Amount: 1,
			Player: Controller,
		},
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

	replicate.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     src,
	})

	if g.Aember(0) != RuleOfSix+2 {
		t.Errorf(
			"Æmber = %d, want %d (the free first pass, six charged bounces, plus the blocked pass's gain)",
			g.Aember(0),
			RuleOfSix+2,
		)
	}
}

func TestTriggerAbilityChargesTheCascadeRoot(t *testing.T) {
	// A Replicator and a differently-named creature reach for each other's reap
	// effect. The whole chain charges the card that started it, so its name pool
	// fills to six while the other name is never touched — not six each.
	reachReap := TriggerAbility{
		Trigger: TriggerAfterReap,
		Target:  Target{Kind: TargetChosenCreature}.Other(),
	}
	g := started(t)
	g.AddToBattleline(testCreature("Replicator", 2, WithAbility(TriggerAfterReap, reachReap)), 0)
	g.AddToBattleline(testCreature("Doppelganger", 2, WithAbility(TriggerAfterReap, reachReap)), 0)
	root := g.Battleline(0)[0]
	other := g.Battleline(0)[1]

	reachReap.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     root,
	})

	if got := g.nameUsagesThisTurn(root); got != RuleOfSix {
		t.Errorf("root name pool = %d, want %d", got, RuleOfSix)
	}
	if got := g.nameUsagesThisTurn(other); got != 0 {
		t.Errorf("other name pool = %d, want 0 (the chain charges the starter)", got)
	}
}

func TestTriggerAbilityNoCandidate(t *testing.T) {
	g := started(t)
	g.AddToBattleline(testCreature("Lonely", 2), 0)
	src := g.Battleline(0)[0]
	TriggerAbility{
		Trigger: TriggerAfterReap,
		Target:  Target{Kind: TargetChosenCreature}.Other(),
	}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     src,
	})
}
