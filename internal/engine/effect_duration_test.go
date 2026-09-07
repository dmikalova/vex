package engine

import (
	"strings"
	"testing"
)

func TestForDurationValidate(t *testing.T) {
	house := BelongToHouse{
		Target:   Target{Kind: TargetThisCreature},
		House:    Mars,
		Duration: EndOfTurn,
	}
	prevent := PreventDamage{Target: Target{Kind: TargetThisCreature}, Duration: EndOfTurn}

	if err := (ForDuration{Duration: NextTurn, Effects: []Effect{house, prevent}}).validate(); err == nil {
		t.Error("want error for unsupported duration")
	}
	if err := (ForDuration{Duration: EndOfTurn, Effects: []Effect{house}}).validate(); err == nil {
		t.Error("want error for fewer than two effects")
	}
	if err := (ForDuration{Duration: EndOfTurn, Effects: []Effect{house, Heal{Fully: true, Target: Target{Kind: TargetThisCreature}}}}).validate(); err == nil {
		t.Error("want error for a child that renders no duration body")
	}
	badChild := PreventDamage{Duration: EndOfTurn} // no target
	if err := (ForDuration{Duration: EndOfTurn, Effects: []Effect{house, badChild}}).validate(); err == nil {
		t.Error("want error surfaced from a misconfigured child")
	}
	if err := (ForDuration{Duration: EndOfTurn, Effects: []Effect{house, prevent}}).validate(); err != nil {
		t.Errorf("valid ForDuration errored: %v", err)
	}
}

func TestForDurationTextSharedSubject(t *testing.T) {
	e := ForDuration{Duration: EndOfTurn, Effects: []Effect{
		BelongToHouse{
			Target:   Target{Kind: TargetTriggeringCreature},
			House:    Sanctum,
			Duration: EndOfTurn,
		},
		PreventDamage{Target: Target{Kind: TargetTriggeringCreature}, Duration: EndOfTurn},
	}}
	want := "for the remainder of the turn, it belongs to house Sanctum and cannot be dealt damage"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

func TestForDurationTextDistinctSubjects(t *testing.T) {
	e := ForDuration{Duration: EndOfTurn, Effects: []Effect{
		BelongToHouse{
			Target:   Target{Kind: TargetThisCreature},
			House:    Sanctum,
			Duration: EndOfTurn,
		},
		PreventDamage{Target: Target{Kind: TargetEachFriendlyCreature}, Duration: EndOfTurn},
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
	e := ForDuration{Duration: EndOfTurn, Effects: []Effect{
		BelongToHouse{Target: Target{Kind: TargetThisCreature}, House: Mars, Duration: EndOfTurn},
		PreventDamage{Target: Target{Kind: TargetThisCreature}, Duration: EndOfTurn},
	}}
	e.Resolve(&EffectContext{Resolver: g, Source: host, Controller: 0})
	if g.House(host) != Mars {
		t.Errorf("house = %s, want Mars", g.House(host))
	}
	g.applyRawDamage(DamageTarget{ID: host, Amount: 2})
	if g.Damage(host) != 0 {
		t.Errorf("damage = %d, want 0 (prevented)", g.Damage(host))
	}
}
