package engine

import "testing"

func TestChooseCreatureThen(t *testing.T) {
	g := NewGame("A", "B", 1)
	ally := g.AddToBattleline(testCreature("ally", 3), 0)
	g.State.Cards[ally].Damage = 2
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := ChooseCreatureThen{
		Target: Target{Kind: TargetChosenCreature},
		Then: Sequence{Effects: []Effect{
			Heal{Fully: true, Target: Target{Kind: TargetTriggeringCreature}},
			PreventDamage{Target: Target{Kind: TargetTriggeringCreature}, Duration: EndOfTurn},
		}},
	}
	want := "choose a creature - fully heal it, and for the remainder of the turn, it cannot be dealt damage"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}

	g.SetChooser(0, &idQueueChooser{ids: []LocalID{ally}})
	e.Resolve(ctx)
	if g.State.Cards[ally].Damage != 0 {
		t.Error("the chosen creature should have been healed")
	}
	if !g.State.Cards[ally].DamageImmune {
		t.Error("the chosen creature should be protected from damage")
	}
}

func TestChooseCreatureThenNoCandidates(_ *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// An empty battleline offers nothing to choose, so Then never resolves (no
	// panic, no effect).
	ChooseCreatureThen{
		Target: Target{Kind: TargetChosenCreature},
		Then:   PreventDamage{Target: Target{Kind: TargetTriggeringCreature}, Duration: EndOfTurn},
	}.Resolve(ctx)
}

func TestChooseCreatureThenValidate(t *testing.T) {
	unsetTarget := ChooseCreatureThen{
		Then: PreventDamage{Target: Target{Kind: TargetTriggeringCreature}, Duration: EndOfTurn},
	}
	if validateEffect(unsetTarget) == nil {
		t.Error("ChooseCreatureThen with an unset Target should fail validation")
	}

	badThen := ChooseCreatureThen{
		Target: Target{Kind: TargetChosenCreature},
		Then:   Heal{Fully: true, Amount: 1, Target: Target{Kind: TargetTriggeringCreature}},
	}
	if validateEffect(badThen) == nil {
		t.Error("ChooseCreatureThen should surface an invalid Then via validate")
	}

	good := ChooseCreatureThen{
		Target: Target{Kind: TargetChosenCreature},
		Then:   PreventDamage{Target: Target{Kind: TargetTriggeringCreature}, Duration: EndOfTurn},
	}
	if validateEffect(good) != nil {
		t.Error("ChooseCreatureThen with a valid Target and Then should pass validation")
	}
}
