package engine

import "testing"

// TestRepeatedFightText covers the rendered phrase and the validation of its
// bounds.
func TestRepeatedFightText(t *testing.T) {
	e := RepeatedFight{Times: 3, Target: Target{Kind: TargetChosenFriendlyCreature}}
	want := "ready and fight with a friendly creature 3 times, each time against " +
		"a different enemy creature. Resolve these fights one at a time"
	if got := e.Text(); got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}

	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if err := (RepeatedFight{Times: 3}).validate(); err == nil {
		t.Error("a targetless RepeatedFight should not validate")
	}
	if err := (RepeatedFight{Target: e.Target}).validate(); err == nil {
		t.Error("a RepeatedFight with no fights should not validate")
	}
}

// TestRepeatedFightNeverFightsTheSameEnemyTwice checks each fight is against an
// enemy no earlier fight used, and that the attacker is readied for every one.
func TestRepeatedFightNeverFightsTheSameEnemyTwice(t *testing.T) {
	g := NewGame("A", "B", 1)
	hero := g.AddToBattleline(testCreature("hero", 10), 0)
	foes := []LocalID{
		g.AddToBattleline(testCreature("a", 1), 1),
		g.AddToBattleline(testCreature("b", 1), 1),
		g.AddToBattleline(testCreature("c", 1), 1),
	}
	ctx := &EffectContext{Resolver: g, Controller: 0}

	RepeatedFight{Times: 3, Target: Target{Kind: TargetChosenFriendlyCreature}}.
		Resolve(ctx)

	if len(g.Battleline(1)) != 0 {
		t.Errorf("enemies left = %d, want 0", len(g.Battleline(1)))
	}
	if len(g.Discard(1)) != len(foes) {
		t.Errorf("discarded = %d, want %d", len(g.Discard(1)), len(foes))
	}
	if got := g.Damage(hero); got != 3 {
		t.Errorf("damage taken = %d, want 3 (one per fight)", got)
	}
}

// TestRepeatedFightStopsWithoutAnEnemy checks the effect stops rather than
// readying its creature when there is no untouched enemy to fight.
func TestRepeatedFightStopsWithoutAnEnemy(t *testing.T) {
	g := NewGame("A", "B", 1)
	hero := g.AddToBattleline(testCreature("hero", 10), 0)
	g.State.Cards[hero].Exhausted = true
	ctx := &EffectContext{Resolver: g, Controller: 0}

	RepeatedFight{Times: 3, Target: Target{Kind: TargetChosenFriendlyCreature}}.
		Resolve(ctx)

	if !g.State.Cards[hero].Exhausted {
		t.Error("with no enemy to fight the creature should not be readied")
	}
}

// TestRepeatedFightStopsWithoutACreature checks the effect stops when its
// controller has nobody to fight with.
func TestRepeatedFightStopsWithoutACreature(t *testing.T) {
	g := NewGame("A", "B", 1)
	foe := g.AddToBattleline(testCreature("foe", 1), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	RepeatedFight{Times: 3, Target: Target{Kind: TargetChosenFriendlyCreature}}.
		Resolve(ctx)

	if g.Damage(foe) != 0 {
		t.Error("with no friendly creature nothing should be fought")
	}
}

// TestRepeatedFightStopsWhenDeclined checks a refused enemy choice ends the
// whole effect instead of moving on to the next fight.
func TestRepeatedFightStopsWhenDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	hero := g.AddToBattleline(testCreature("hero", 10), 0)
	g.AddToBattleline(testCreature("a", 1), 1)
	g.AddToBattleline(testCreature("b", 1), 1)
	g.SetChooser(0, orderRejectChooser{})
	ctx := &EffectContext{Resolver: g, Controller: 0}

	RepeatedFight{Times: 2, Target: Target{Kind: TargetChosenFriendlyCreature}}.
		Resolve(ctx)

	if len(g.Battleline(1)) != 2 {
		t.Error("a declined choice should leave every enemy alone")
	}
	if g.Damage(hero) != 0 {
		t.Error("a declined choice should not fight")
	}
}
