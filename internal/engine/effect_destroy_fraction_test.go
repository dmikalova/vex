package engine

import "testing"

// TestFractionOf covers the round-up arithmetic on both fractions.
func TestFractionOf(t *testing.T) {
	cases := []struct {
		f    Fraction
		n    int
		want int
	}{
		{OneThird, 0, 0},
		{OneThird, 1, 1},
		{OneThird, 3, 1},
		{OneThird, 4, 2},
		{OneHalf, 3, 2},
		{OneHalf, 4, 2},
	}
	for _, c := range cases {
		if got := c.f.of(c.n); got != c.want {
			t.Errorf("%s.of(%d) = %d, want %d", c.f.word(), c.n, got, c.want)
		}
	}
}

// TestDestroyFractionValidate covers the unset-portion guard and success.
func TestDestroyFractionValidate(t *testing.T) {
	if err := (DestroyFractionOfEachBattleline{}).validate(); err == nil {
		t.Error("an unset Portion should be invalid")
	}
	if err := (DestroyFractionOfEachBattleline{Portion: OneThird}).validate(); err != nil {
		t.Errorf("a set Portion should pass: %v", err)
	}
}

// TestDestroyFractionText covers the printed wording.
func TestDestroyFractionText(t *testing.T) {
	want := "destroy one third of all enemy creatures and one third of all " +
		"friendly creatures (rounding up each time)"
	if got := (DestroyFractionOfEachBattleline{Portion: OneThird}).Text(); got != want {
		t.Errorf("Text() = %q, want %q", got, want)
	}
}

// TestDestroyFractionResolve covers the choose-and-destroy path: ceil(4/3)=2 enemy
// and ceil(4/3)=2 friendly creatures are destroyed, chosen by the controller.
func TestDestroyFractionResolve(t *testing.T) {
	g := NewGame("A", "B", 1)
	f0 := g.AddToBattleline(testCreature("f0", 3), 0)
	f1 := g.AddToBattleline(testCreature("f1", 3), 0)
	f2 := g.AddToBattleline(testCreature("f2", 3), 0)
	f3 := g.AddToBattleline(testCreature("f3", 3), 0)
	e0 := g.AddToBattleline(testCreature("e0", 3), 1)
	e1 := g.AddToBattleline(testCreature("e1", 3), 1)
	e2 := g.AddToBattleline(testCreature("e2", 3), 1)
	e3 := g.AddToBattleline(testCreature("e3", 3), 1)

	// Destroy e0,e1 on the enemy side (chosen first) and f0,f1 on the friendly side.
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{e0, e1, f0, f1}})

	DestroyFractionOfEachBattleline{Portion: OneThird}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.inPlay(e0) || g.inPlay(e1) {
		t.Error("the chosen enemy creatures should be destroyed")
	}
	if !g.inPlay(e2) || !g.inPlay(e3) {
		t.Error("the unchosen enemy creatures should survive")
	}
	if g.inPlay(f0) || g.inPlay(f1) {
		t.Error("the chosen friendly creatures should be destroyed")
	}
	if !g.inPlay(f2) || !g.inPlay(f3) {
		t.Error("the unchosen friendly creatures should survive")
	}
}

// TestDestroyFractionRoundingAndEmpty covers ceil(1/3)=1 on a one-creature side
// (that lone creature is taken) and an empty side (nothing to choose or destroy).
func TestDestroyFractionRoundingAndEmpty(t *testing.T) {
	g := NewGame("A", "B", 1)
	f0 := g.AddToBattleline(testCreature("f0", 3), 0)

	// One friendly creature: ceil(1/3)=1, so it is destroyed. The enemy side is empty.
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{f0}})

	DestroyFractionOfEachBattleline{Portion: OneThird}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.inPlay(f0) {
		t.Error("ceil(1/3)=1, so the lone creature should be destroyed")
	}
}

// TestDestroyFractionNoCreatures covers both sides empty: the loop never runs.
func TestDestroyFractionNoCreatures(_ *testing.T) {
	g := NewGame("A", "B", 1)
	DestroyFractionOfEachBattleline{Portion: OneThird}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})
}

// TestDestroyFractionDeclinedChoice covers the !ok break: a rejecting chooser
// destroys nothing.
func TestDestroyFractionDeclinedChoice(t *testing.T) {
	g := NewGame("A", "B", 1)
	f0 := g.AddToBattleline(testCreature("f0", 3), 0)
	f1 := g.AddToBattleline(testCreature("f1", 3), 0)
	g.SetChooser(0, orderRejectChooser{})

	DestroyFractionOfEachBattleline{Portion: OneThird}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})

	if !g.inPlay(f0) || !g.inPlay(f1) {
		t.Error("a declined choice destroys nothing")
	}
}
