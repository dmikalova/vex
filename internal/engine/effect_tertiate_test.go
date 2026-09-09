package engine

import "testing"

// TestTertiateText covers the printed wording.
func TestTertiateText(t *testing.T) {
	want := "destroy one third of all enemy creatures and one third of all " +
		"friendly creatures (rounding up each time)"
	if got := (Tertiate{}).Text(); got != want {
		t.Errorf("Text() = %q, want %q", got, want)
	}
}

// TestTertiateResolve covers the choose-and-destroy path: ceil(4/3)=2 enemy and
// ceil(3/3)=1 friendly creatures are destroyed, and the controller's specific
// picks are the ones that die.
func TestTertiateResolve(t *testing.T) {
	g := NewGame("A", "B", 1)
	// Friendly creatures (player 0).
	f0 := g.AddToBattleline(testCreature("f0", 3), 0)
	f1 := g.AddToBattleline(testCreature("f1", 3), 0)
	f2 := g.AddToBattleline(testCreature("f2", 3), 0)
	// Enemy creatures (player 1).
	e0 := g.AddToBattleline(testCreature("e0", 3), 1)
	e1 := g.AddToBattleline(testCreature("e1", 3), 1)
	e2 := g.AddToBattleline(testCreature("e2", 3), 1)
	e3 := g.AddToBattleline(testCreature("e3", 3), 1)

	// Controller (player 0) chooses e0 then e2 on the enemy side; the friendly
	// side's single pick falls back to the first candidate, f0.
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{e0, e2}})

	Tertiate{}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	// Enemy: 2 of 4 destroyed (the chosen ones); the other two survive.
	if g.inPlay(e0) || g.inPlay(e2) {
		t.Error("the chosen enemy creatures should be destroyed")
	}
	if !g.inPlay(e1) || !g.inPlay(e3) {
		t.Error("the unchosen enemy creatures should survive")
	}
	// Friendly: 1 of 3 destroyed.
	if g.inPlay(f0) {
		t.Error("the chosen friendly creature should be destroyed")
	}
	if !g.inPlay(f1) || !g.inPlay(f2) {
		t.Error("only one friendly creature should be destroyed")
	}
}

// TestTertiateRoundingAndEmpty covers ceil(1/3)=1 on a one-creature side (taken
// automatically as the sole candidate) and an empty side (nothing destroyed).
func TestTertiateRoundingAndEmpty(t *testing.T) {
	g := NewGame("A", "B", 1)
	e0 := g.AddToBattleline(testCreature("e0", 3), 1)
	f0 := g.AddToBattleline(testCreature("f0", 3), 0)

	Tertiate{}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	// One enemy -> ceil(1/3)=1 destroyed; one friendly -> ceil(1/3)=1 destroyed.
	if g.inPlay(e0) {
		t.Error("the sole enemy creature should be destroyed")
	}
	if g.inPlay(f0) {
		t.Error("the sole friendly creature should be destroyed")
	}
}

// TestTertiateNoCreatures covers both sides empty: the loop never runs.
func TestTertiateNoCreatures(_ *testing.T) {
	g := NewGame("A", "B", 1)
	Tertiate{}.Resolve(&EffectContext{Resolver: g, Controller: 0})
}

// TestTertiateDeclinedChoice covers the !ok break: a rejecting chooser declines
// the enemy pick, so no enemy creature is destroyed.
func TestTertiateDeclinedChoice(t *testing.T) {
	g := NewGame("A", "B", 1)
	e0 := g.AddToBattleline(testCreature("e0", 3), 1)
	e1 := g.AddToBattleline(testCreature("e1", 3), 1)
	e2 := g.AddToBattleline(testCreature("e2", 3), 1)
	e3 := g.AddToBattleline(testCreature("e3", 3), 1)
	g.SetChooser(0, orderRejectChooser{})

	Tertiate{}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if !g.inPlay(e0) || !g.inPlay(e1) || !g.inPlay(e2) || !g.inPlay(e3) {
		t.Error("a declined choice should destroy no creatures")
	}
}
