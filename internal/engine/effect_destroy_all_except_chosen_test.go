package engine

import "testing"

// TestDestroyAllExceptChosenValidate covers both keep-count guards and success.
func TestDestroyAllExceptChosenValidate(t *testing.T) {
	if err := (DestroyAllExceptChosen{FriendlyKept: 0, EnemyKept: 3}).validate(); err == nil {
		t.Error("a non-positive FriendlyKept should be invalid")
	}
	if err := (DestroyAllExceptChosen{FriendlyKept: 3, EnemyKept: 0}).validate(); err == nil {
		t.Error("a non-positive EnemyKept should be invalid")
	}
	if err := (DestroyAllExceptChosen{FriendlyKept: 3, EnemyKept: 3}).validate(); err != nil {
		t.Errorf("valid keep counts should pass: %v", err)
	}
}

// TestDestroyAllExceptChosenText covers the printed wording.
func TestDestroyAllExceptChosenText(t *testing.T) {
	want := "choose 3 friendly creatures and 3 enemy creatures. Destroy each " +
		"other creature"
	got := (DestroyAllExceptChosen{FriendlyKept: 3, EnemyKept: 3}).Text()
	if got != want {
		t.Errorf("Text() = %q, want %q", got, want)
	}
}

// TestDestroyAllExceptChosenResolve covers more than the keep count on both sides:
// the controller keeps 3 of 4 per side, the unchosen creature on each side is
// destroyed, and the chosen ones survive.
func TestDestroyAllExceptChosenResolve(t *testing.T) {
	g := NewGame("A", "B", 1)
	f0 := g.AddToBattleline(testCreature("f0", 3), 0)
	f1 := g.AddToBattleline(testCreature("f1", 3), 0)
	f2 := g.AddToBattleline(testCreature("f2", 3), 0)
	f3 := g.AddToBattleline(testCreature("f3", 3), 0)
	e0 := g.AddToBattleline(testCreature("e0", 3), 1)
	e1 := g.AddToBattleline(testCreature("e1", 3), 1)
	e2 := g.AddToBattleline(testCreature("e2", 3), 1)
	e3 := g.AddToBattleline(testCreature("e3", 3), 1)

	// Keep f0,f1,f2 on the friendly side and e0,e1,e2 on the enemy side.
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{f0, f1, f2, e0, e1, e2}})

	DestroyAllExceptChosen{FriendlyKept: 3, EnemyKept: 3}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})

	if !g.inPlay(f0) || !g.inPlay(f1) || !g.inPlay(f2) {
		t.Error("the chosen friendly creatures should survive")
	}
	if g.inPlay(f3) {
		t.Error("the unchosen friendly creature should be destroyed")
	}
	if !g.inPlay(e0) || !g.inPlay(e1) || !g.inPlay(e2) {
		t.Error("the chosen enemy creatures should survive")
	}
	if g.inPlay(e3) {
		t.Error("the unchosen enemy creature should be destroyed")
	}
}

// TestDestroyAllExceptChosenFewerAndEmpty covers a side with fewer creatures than
// the keep count (all kept, none destroyed) and an empty side (nothing to choose
// or destroy).
func TestDestroyAllExceptChosenFewerAndEmpty(t *testing.T) {
	g := NewGame("A", "B", 1)
	f0 := g.AddToBattleline(testCreature("f0", 3), 0)
	f1 := g.AddToBattleline(testCreature("f1", 3), 0)

	// Only two friendly creatures exist, so both are kept; the enemy side is empty.
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{f0, f1}})

	DestroyAllExceptChosen{FriendlyKept: 3, EnemyKept: 3}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})

	if !g.inPlay(f0) || !g.inPlay(f1) {
		t.Error("both friendly creatures should survive when fewer than the keep count")
	}
}

// TestDestroyAllExceptChosenDeclined covers the !ok break: a rejecting chooser
// keeps nothing, so every creature is destroyed.
func TestDestroyAllExceptChosenDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	f0 := g.AddToBattleline(testCreature("f0", 3), 0)
	f1 := g.AddToBattleline(testCreature("f1", 3), 0)
	g.SetChooser(0, orderRejectChooser{})

	DestroyAllExceptChosen{FriendlyKept: 3, EnemyKept: 3}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.inPlay(f0) || g.inPlay(f1) {
		t.Error("a declined choice keeps nothing, so all creatures are destroyed")
	}
}
