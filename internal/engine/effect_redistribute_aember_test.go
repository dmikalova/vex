package engine

import "testing"

func TestRedistributeCapturedAemberText(t *testing.T) {
	if got := (RedistributeCapturedAember{Side: Controller}).Text(); got != "redistribute the Æmber on friendly creatures among friendly creatures" {
		t.Errorf("friendly text = %q", got)
	}
	if got := (RedistributeCapturedAember{Side: Opponent}).Text(); got != "redistribute the Æmber on enemy creatures among enemy creatures" {
		t.Errorf("enemy text = %q", got)
	}
}

func TestRedistributeCapturedAemberValidate(t *testing.T) {
	if err := validateEffect(RedistributeCapturedAember{Side: Controller}); err != nil {
		t.Errorf("validate = %v", err)
	}
}

func TestRedistributeCapturedAemberMovesAll(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 3), 0)
	b := g.AddToBattleline(testCreature("b", 3), 0)
	g.AddAmberOn(a, 4)
	g.AddAmberOn(b, 0)
	// Pile all 4 onto b.
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{b, b, b, b}})

	ctx := &EffectContext{Resolver: g, Source: a, Controller: 0}
	RedistributeCapturedAember{Side: Controller}.Resolve(ctx)

	if got := g.AmberOn(a); got != 0 {
		t.Errorf("a = %d, want 0", got)
	}
	if got := g.AmberOn(b); got != 4 {
		t.Errorf("b = %d, want 4", got)
	}
}

func TestRedistributeCapturedAemberEnemySide(t *testing.T) {
	g := NewGame("A", "B", 1)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	g.AddAmberOn(foe, 2)

	// One enemy creature: the choice auto-resolves back onto it.
	ctx := &EffectContext{Resolver: g, Source: foe, Controller: 0}
	RedistributeCapturedAember{Side: Opponent}.Resolve(ctx)

	if got := g.AmberOn(foe); got != 2 {
		t.Errorf("foe = %d, want 2", got)
	}
}

func TestRedistributeCapturedAemberNoAember(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 3), 0)

	ctx := &EffectContext{Resolver: g, Source: a, Controller: 0}
	RedistributeCapturedAember{Side: Controller}.Resolve(ctx)

	if got := g.AmberOn(a); got != 0 {
		t.Errorf("a = %d, want 0", got)
	}
}

func TestRedistributeCapturedAemberDeclineFallsBackToFirst(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 3), 0)
	b := g.AddToBattleline(testCreature("b", 3), 0)
	g.AddAmberOn(a, 1)
	// A chooser that declines: the effect falls back to the first candidate (a).
	g.SetChooser(0, orderRejectChooser{})

	ctx := &EffectContext{Resolver: g, Source: a, Controller: 0}
	RedistributeCapturedAember{Side: Controller}.Resolve(ctx)

	if got := g.AmberOn(a); got != 1 {
		t.Errorf("a = %d, want 1", got)
	}
	if got := g.AmberOn(b); got != 0 {
		t.Errorf("b = %d, want 0", got)
	}
}
