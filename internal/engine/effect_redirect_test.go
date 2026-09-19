package engine

import "testing"

func TestRedirectFightDamage(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("gabos", 5), 0)
	g.AddToBattleline(testCreature("foe", 3), 1)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}

	e := RedirectFightDamage{Target: Target{Kind: TargetChosenCreature}}
	want := "choose a creature. {self} deals its fight damage to the chosen creature instead of to the creature it is fighting"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
	// The default chooser picks the first candidate (the controller's creature).
	e.Resolve(ctx)
	if g.State.FightDamageRedirect != src {
		t.Errorf("redirect = %d, want %d (first candidate)", g.State.FightDamageRedirect, src)
	}
}

func TestRedirectFightDamageDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("gabos", 5), 0)
	g.AddToBattleline(testCreature("ally", 2), 0) // two candidates, so the chooser is consulted
	g.AddToBattleline(testCreature("foe", 3), 1)
	g.SetChooser(0, orderRejectChooser{})
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}

	RedirectFightDamage{Target: Target{Kind: TargetChosenCreature}}.Resolve(ctx)
	if g.State.FightDamageRedirect != 0 {
		t.Errorf("redirect = %d, want 0 (choice declined)", g.State.FightDamageRedirect)
	}
}
