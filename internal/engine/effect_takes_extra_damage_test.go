package engine

import "testing"

func TestTakesExtraDamageText(t *testing.T) {
	e := TakesExtraDamage{Target: Target{Kind: TargetChosenCreature}, Amount: 2}
	want := "for the remainder of the turn, whenever a creature takes damage, " +
		"it takes an additional 2 damage"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
}

func TestTakesExtraDamageValidate(t *testing.T) {
	if err := validateEffect(
		TakesExtraDamage{Target: Target{Kind: TargetChosenCreature}, Amount: 2},
	); err != nil {
		t.Errorf("valid effect rejected: %v", err)
	}
	if validateEffect(TakesExtraDamage{Amount: 2}) == nil {
		t.Error("want error for missing target")
	}
	if validateEffect(
		TakesExtraDamage{Target: Target{Kind: TargetChosenCreature}, Amount: 0},
	) == nil {
		t.Error("want error for non-positive amount")
	}
}

func TestTakesExtraDamageAugmentsDamage(t *testing.T) {
	g := started(t)
	victim := g.AddToBattleline(testCreature("victim", 10), 0)
	TakesExtraDamage{Target: Target{Kind: TargetThisCreature}, Amount: 2}.Resolve(
		&EffectContext{Resolver: g, Source: victim, Controller: 0},
	)
	// A single instance of 3 damage lands as 3 + 2 = 5.
	g.dealDamage(0, DamageTarget{ID: victim, Amount: 3})
	if got := g.Damage(victim); got != 5 {
		t.Errorf("damage = %d, want 5 (3 + 2 bonus)", got)
	}
	// A creature with no augmentation takes only what it is dealt.
	bare := g.AddToBattleline(testCreature("bare", 10), 0)
	g.dealDamage(0, DamageTarget{ID: bare, Amount: 3})
	if got := g.Damage(bare); got != 3 {
		t.Errorf("bare damage = %d, want 3", got)
	}
}

func TestLastingExtraDamageQuery(t *testing.T) {
	g := started(t)
	a := g.AddToBattleline(testCreature("a", 5), 0)
	b := g.AddToBattleline(testCreature("b", 5), 0)
	g.AddLasting(
		LastingEffect{On: EventCreatureTakesDamage, Amount: 2, Subject: a, HasSubject: true},
	)
	if got := g.lastingExtraDamage(a); got != 2 {
		t.Errorf("extra damage on a = %d, want 2", got)
	}
	if got := g.lastingExtraDamage(b); got != 0 {
		t.Errorf("extra damage on b = %d, want 0", got)
	}
}
