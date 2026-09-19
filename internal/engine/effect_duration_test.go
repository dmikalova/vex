package engine

import (
	"strings"
	"testing"
)

func TestForDurationValidate(t *testing.T) {
	house := BelongToHouse{
		Target:   Target{Kind: TargetThisCreature},
		House:    Mars,
		Duration: RemainderOfPlayerTurn,
	}
	cannotDamage := CannotBeDealtDamage{
		Target:   Target{Kind: TargetThisCreature},
		Duration: RemainderOfPlayerTurn,
	}

	if err := (ForDuration{
		Duration: OpponentNextTurn,
		Effects:  []Effect{house, cannotDamage},
	}).validate(); err == nil {
		t.Error("want error for unsupported duration")
	}
	if err := (ForDuration{
		Duration: RemainderOfPlayerTurn,
		Effects:  []Effect{house},
	}).validate(); err == nil {
		t.Error("want error for fewer than two effects")
	}
	if err := (ForDuration{
		Duration: RemainderOfPlayerTurn,
		Effects: []Effect{house, Heal{
			Fully:  true,
			Target: Target{Kind: TargetThisCreature},
		}},
	}).validate(); err == nil {
		t.Error("want error for a child that renders no duration body")
	}
	badChild := CannotBeDealtDamage{Duration: RemainderOfPlayerTurn} // no target
	if err := (ForDuration{
		Duration: RemainderOfPlayerTurn,
		Effects:  []Effect{house, badChild},
	}).validate(); err == nil {
		t.Error("want error surfaced from a misconfigured child")
	}
	if err := (ForDuration{
		Duration: RemainderOfPlayerTurn,
		Effects:  []Effect{house, cannotDamage},
	}).validate(); err != nil {
		t.Errorf("valid ForDuration errored: %v", err)
	}
}

func TestForDurationTextSharedSubject(t *testing.T) {
	e := ForDuration{Duration: RemainderOfPlayerTurn, Effects: []Effect{
		BelongToHouse{
			Target:   Target{Kind: TargetTriggeringCreature},
			House:    Sanctum,
			Duration: RemainderOfPlayerTurn,
		},
		CannotBeDealtDamage{
			Target:   Target{Kind: TargetTriggeringCreature},
			Duration: RemainderOfPlayerTurn,
		},
	}}
	want := "for the remainder of the turn, it belongs to house Sanctum and cannot be dealt damage"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

func TestForDurationTextDistinctSubjects(t *testing.T) {
	e := ForDuration{Duration: RemainderOfPlayerTurn, Effects: []Effect{
		BelongToHouse{
			Target:   Target{Kind: TargetThisCreature},
			House:    Sanctum,
			Duration: RemainderOfPlayerTurn,
		},
		CannotBeDealtDamage{
			Target:   Target{Kind: TargetEachFriendlyCreature},
			Duration: RemainderOfPlayerTurn,
		},
	}}
	got := e.Text()
	if !strings.Contains(
		got,
		SelfName+" belongs to house Sanctum and each friendly creature cannot be dealt damage",
	) {
		t.Errorf("text = %q", got)
	}
}

func TestForDurationResolvesChildren(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(NewCard("Host", Brobnar, Creature, Common, WithPower(3)), 0)
	e := ForDuration{Duration: RemainderOfPlayerTurn, Effects: []Effect{
		BelongToHouse{
			Target:   Target{Kind: TargetThisCreature},
			House:    Mars,
			Duration: RemainderOfPlayerTurn,
		},
		CannotBeDealtDamage{
			Target:   Target{Kind: TargetThisCreature},
			Duration: RemainderOfPlayerTurn,
		},
	}}
	e.Resolve(&EffectContext{
		Resolver:   g,
		Source:     host,
		Controller: 0,
	})
	if g.House(host) != Mars {
		t.Errorf("house = %s, want Mars", g.House(host))
	}
	g.applyRawDamage(DamageTarget{
		ID:     host,
		Amount: 2,
	})
	if g.Damage(host) != 0 {
		t.Errorf("damage = %d, want 0 (prevented)", g.Damage(host))
	}
}

// TestGainUntilNextTurnText covers the folded clause: children sharing a subject
// name it once and join their predicates under one "until the start of your next
// turn" suffix.
func TestGainUntilNextTurnText(t *testing.T) {
	shared := GainUntilNextTurn{
		Effects: []Effect{
			GainKeywords{
				Target:   Target{Kind: TargetTriggeringCreature},
				Keywords: []Keyword{Elusive},
				Duration: StartOfPlayerNextTurn,
			},
			GainTrait{
				Target: Target{Kind: TargetTriggeringCreature},
				Trait:  Mutant,
			},
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
			GainKeywords{
				Target:   Target{Kind: TargetTriggeringCreature},
				Keywords: []Keyword{Elusive},
				Duration: StartOfPlayerNextTurn,
			},
			GainTrait{
				Target: Target{Kind: TargetThisCreature},
				Trait:  Mutant,
			},
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
			GainKeywords{
				Target:   Target{Kind: TargetTriggeringCreature},
				Keywords: []Keyword{Elusive},
				Duration: StartOfPlayerNextTurn,
			},
		},
	}
	if one.validate() == nil {
		t.Error("fewer than two effects should be invalid")
	}
	notScoped := GainUntilNextTurn{
		Effects: []Effect{
			GainKeywords{
				Target:   Target{Kind: TargetTriggeringCreature},
				Keywords: []Keyword{Elusive},
				Duration: StartOfPlayerNextTurn,
			},
			Draw{Amount: 1},
		},
	}
	if notScoped.validate() == nil {
		t.Error("a child that renders no duration body should be invalid")
	}
	badChild := GainUntilNextTurn{
		Effects: []Effect{
			GainKeywords{
				Target:   Target{Kind: TargetTriggeringCreature},
				Keywords: []Keyword{Elusive},
				Duration: StartOfPlayerNextTurn,
			},
			GainTrait{Target: Target{Kind: TargetTriggeringCreature}},
		},
	}
	if badChild.validate() == nil {
		t.Error("an invalid child should be rejected")
	}
	if (GainUntilNextTurn{
		Effects: []Effect{
			GainKeywords{
				Target:   Target{Kind: TargetTriggeringCreature},
				Keywords: []Keyword{Elusive},
				Duration: StartOfPlayerNextTurn,
			},
			GainTrait{
				Target: Target{Kind: TargetTriggeringCreature},
				Trait:  Mutant,
			},
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
			GainKeywords{
				Target:   Target{Kind: TargetTriggeringCreature},
				Keywords: []Keyword{Skirmish},
				Duration: StartOfPlayerNextTurn,
			},
			GainAssaultUntilNextTurn{
				Target: Target{Kind: TargetTriggeringCreature},
				Amount: Fixed(3),
			},
			GainTrait{
				Target: Target{Kind: TargetTriggeringCreature},
				Trait:  Mutant,
			},
		},
	}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		It:         beast,
		HasIt:      true,
	})

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
