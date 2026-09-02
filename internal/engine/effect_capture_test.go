package engine

import "testing"

func TestCaptureAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 2

	e := CaptureAember{Amount: 3, Target: Target{Kind: TargetThisCreature}, Source: Opponent}
	if e.Text() != "{self} captures 3 Æmber from your opponent" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx) // opponent has only 2
	if g.State.Cards[src].Amber != 2 {
		t.Errorf("captured = %d, want 2", g.State.Cards[src].Amber)
	}
	if g.State.Aember[1] != 0 {
		t.Errorf("opponent aember = %d, want 0", g.State.Aember[1])
	}
}

func TestCaptureAllAember(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("drumble", 2), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 7

	e := CaptureAember{All: true, Target: Target{Kind: TargetThisCreature}, Source: Opponent}
	if e.Text() != "{self} captures all your opponent's Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.AmberOn(src) != 7 {
		t.Errorf("captured = %d, want 7 (all of the opponent's pool)", g.AmberOn(src))
	}
	if g.Aember(1) != 0 {
		t.Errorf("opponent aember = %d, want 0", g.Aember(1))
	}
}

func TestCaptureAemberBy(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("gate", 5), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 9

	e := CaptureAember{By: AllBut(5), Target: Target{Kind: TargetThisCreature}, Source: Opponent}
	if want := "{self} captures all but 5 Æmber from your opponent"; e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
	e.Resolve(ctx)
	if g.AmberOn(src) != 4 {
		t.Errorf("captured = %d, want 4", g.AmberOn(src))
	}
	if g.Aember(1) != 5 {
		t.Errorf("opponent aember = %d, want 5", g.Aember(1))
	}
}

func TestCaptureAemberText(t *testing.T) {
	this := Target{Kind: TargetThisCreature}
	enemy := Target{Kind: TargetChosenEnemyCreature}
	cases := []struct {
		name string
		e    CaptureAember
		want string
	}{
		{
			"friendly from opponent",
			CaptureAember{Amount: 1, Target: this, Source: Opponent},
			"{self} captures 1 Æmber from your opponent",
		},
		{
			"friendly from own side",
			CaptureAember{Amount: 1, Target: this, Source: Controller},
			"{self} captures 1 Æmber from your own side",
		},
		{
			"enemy from their own side",
			CaptureAember{Amount: 1, Target: enemy, Source: Opponent},
			"an enemy creature captures 1 Æmber from their own side",
		},
		{
			"enemy from your pool",
			CaptureAember{Amount: 1, Target: enemy, Source: Controller},
			"an enemy creature captures 1 Æmber from you",
		},
		{
			"all from own pool",
			CaptureAember{All: true, Target: this, Source: Controller},
			"{self} captures all your Æmber",
		},
	}
	for _, tc := range cases {
		if got := tc.e.Text(); got != tc.want {
			t.Errorf("%s: text = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestCaptureAemberValidate(t *testing.T) {
	this := Target{Kind: TargetThisCreature}
	if err := validateEffect(CaptureAember{Amount: 1, Source: Opponent}); err == nil {
		t.Error("unset Target should fail validation")
	}
	if err := validateEffect(CaptureAember{Amount: 1, Target: this}); err == nil {
		t.Error("unset Source should fail validation")
	}
	if err := validateEffect(CaptureAember{Amount: 1, Target: this, Source: Opponent}); err != nil {
		t.Errorf("valid: %v, want nil", err)
	}
	both := CaptureAember{Amount: 1, By: AllBut(5), Target: this, Source: Opponent}
	if err := validateEffect(both); err == nil {
		t.Error("setting both Amount and By should fail validation")
	}
}

func marsCreature(name string, power int) CardDefinition {
	return NewCard(name, Mars, Creature, Common, WithPower(power), WithTraits("Martian"))
}

func TestCaptureAemberByEnemy(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(marsCreature("command", 1), 0)
	g.AddToBattleline(marsCreature("m1", 3), 0) // two friendly Mars creatures -> Per = 2
	foe := g.AddToBattleline(testCreature("foe", 4), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 3

	e := CaptureAember{
		Amount: 1,
		Target: Target{Kind: TargetChosenEnemyCreature},
		Source: Opponent,
		Per:    InPlay{Player: Controller, Type: Creature, House: Mars},
	}
	want := "for each friendly Mars creature, an enemy creature captures 1 Æmber from their own side"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
	e.Resolve(ctx)
	// Two friendly Mars creatures, so the enemy captures twice; the default chooser
	// picks foe both times.
	if g.AmberOn(foe) != 2 {
		t.Errorf("captured on foe = %d, want 2", g.AmberOn(foe))
	}
	if g.Aember(1) != 1 {
		t.Errorf("opponent pool = %d, want 1 (3 - 2 captured)", g.Aember(1))
	}
}

func TestCaptureAemberByEnemyNoEnemies(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(marsCreature("command", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 2

	// One friendly Mars creature but no enemy creatures: the loop finds nothing to
	// capture with and stops.
	e := CaptureAember{
		Amount: 1,
		Target: Target{Kind: TargetChosenEnemyCreature},
		Source: Opponent,
		Per:    InPlay{Player: Controller, Type: Creature, House: Mars},
	}
	e.Resolve(ctx)
	if g.Aember(1) != 2 {
		t.Errorf("opponent pool = %d, want 2 (nothing captured)", g.Aember(1))
	}
}

func TestCaptureAemberByEnemyDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(marsCreature("command", 1), 0)
	g.AddToBattleline(testCreature("foe1", 4), 1)
	g.AddToBattleline(testCreature("foe2", 5), 1) // two candidates, so the chooser is consulted
	g.SetChooser(0, orderRejectChooser{})
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 2

	e := CaptureAember{
		Amount: 1,
		Target: Target{Kind: TargetChosenEnemyCreature},
		Source: Opponent,
		Per:    InPlay{Player: Controller, Type: Creature, House: Mars},
	}
	e.Resolve(ctx)
	if g.Aember(1) != 2 {
		t.Errorf("opponent pool = %d, want 2 (choice declined)", g.Aember(1))
	}
}
