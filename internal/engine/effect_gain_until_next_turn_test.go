package engine

import "testing"

// TestGainUntilNextTurnText covers the folded clause: children sharing a subject
// name it once and join their predicates under one "until the start of your next
// turn" suffix.
func TestGainUntilNextTurnText(t *testing.T) {
	shared := GainUntilNextTurn{
		Effects: []Effect{
			GainKeyword{Target: Target{Kind: TargetTriggeringCreature}, Keyword: Elusive},
			GainTrait{Target: Target{Kind: TargetTriggeringCreature}, Trait: Mutant},
		},
	}
	want := "it gains elusive and the Mutant trait until the start of your next turn"
	if got := shared.Text(); got != want {
		t.Errorf("shared text = %q, want %q", got, want)
	}

	// Children acting on different subjects render each full body, still under one
	// shared suffix.
	unshared := GainUntilNextTurn{
		Effects: []Effect{
			GainKeyword{Target: Target{Kind: TargetTriggeringCreature}, Keyword: Elusive},
			GainTrait{Target: Target{Kind: TargetThisCreature}, Trait: Mutant},
		},
	}
	want = "it gains elusive and " + SelfName +
		" gains the Mutant trait until the start of your next turn"
	if got := unshared.Text(); got != want {
		t.Errorf("unshared text = %q, want %q", got, want)
	}
}

// TestGainUntilNextTurnValidate rejects fewer than two children and a child that
// renders no duration body.
func TestGainUntilNextTurnValidate(t *testing.T) {
	one := GainUntilNextTurn{
		Effects: []Effect{
			GainKeyword{Target: Target{Kind: TargetTriggeringCreature}, Keyword: Elusive},
		},
	}
	if one.validate() == nil {
		t.Error("fewer than two effects should be invalid")
	}
	notScoped := GainUntilNextTurn{
		Effects: []Effect{
			GainKeyword{Target: Target{Kind: TargetTriggeringCreature}, Keyword: Elusive},
			Draw{Amount: 1},
		},
	}
	if notScoped.validate() == nil {
		t.Error("a child that renders no duration body should be invalid")
	}
	badChild := GainUntilNextTurn{
		Effects: []Effect{
			GainKeyword{Target: Target{Kind: TargetTriggeringCreature}, Keyword: Elusive},
			GainTrait{Target: Target{Kind: TargetTriggeringCreature}},
		},
	}
	if badChild.validate() == nil {
		t.Error("an invalid child should be rejected")
	}
	if (GainUntilNextTurn{
		Effects: []Effect{
			GainKeyword{Target: Target{Kind: TargetTriggeringCreature}, Keyword: Elusive},
			GainTrait{Target: Target{Kind: TargetTriggeringCreature}, Trait: Mutant},
		},
	}).validate() != nil {
		t.Error("a valid fold should pass")
	}
}

// TestGainUntilNextTurnResolve applies every folded grant to the creature in
// context in order.
func TestGainUntilNextTurnResolve(t *testing.T) {
	g := started(t)
	beast := g.AddToBattleline(testCreature("beast", 4), 0)

	GainUntilNextTurn{
		Effects: []Effect{
			GainKeyword{Target: Target{Kind: TargetTriggeringCreature}, Keyword: Skirmish},
			GainAssaultUntilNextTurn{
				Target: Target{Kind: TargetTriggeringCreature},
				Amount: Fixed(3),
			},
			GainTrait{Target: Target{Kind: TargetTriggeringCreature}, Trait: Mutant},
		},
	}.Resolve(&EffectContext{Resolver: g, Controller: 0, It: beast, HasIt: true})

	if !g.hasKeyword(beast, Skirmish) {
		t.Error("the creature should have gained skirmish")
	}
	if g.assault(beast) != 3 {
		t.Errorf("assault = %d, want 3", g.assault(beast))
	}
	if !g.HasTrait(beast, Mutant) {
		t.Error("the creature should have gained the Mutant trait")
	}
}
