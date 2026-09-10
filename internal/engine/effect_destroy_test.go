package engine

import "testing"

func TestDestroyEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("shaker", 7), 0)
	weakFriendly := g.AddToBattleline(testCreature("weak", 3), 0)
	strongEnemy := g.AddToBattleline(testCreature("strong", 5), 1)
	weakEnemy := g.AddToBattleline(testCreature("weakfoe", 2), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	byPower := Destroy{Target: Target{Kind: TargetEachCreature}.PowerAtMost(3)}
	if byPower.Text() != "destroy each creature with power 3 or lower" {
		t.Errorf("power text = %q", byPower.Text())
	}
	byPower.Resolve(ctx)
	if g.inPlay(weakFriendly) || g.inPlay(weakEnemy) {
		t.Error("power<=3 creatures should be destroyed")
	}
	if !g.inPlay(src) || !g.inPlay(strongEnemy) {
		t.Error("power>3 creatures should survive")
	}

	sci := g.AddToBattleline(
		NewCard("sci", Logos, Creature, Common, WithPower(6), WithTraits(Scientist)),
		1,
	)
	byTrait := Destroy{Target: Target{Kind: TargetEachCreature}.WithTrait(Scientist)}
	if byTrait.Text() != "destroy each Scientist creature" {
		t.Errorf("trait text = %q", byTrait.Text())
	}
	byTrait.Resolve(ctx)
	if g.inPlay(sci) {
		t.Error("Scientist creature should be destroyed")
	}
	if !g.inPlay(strongEnemy) {
		t.Error("non-Scientist creature should survive")
	}
}

// TestDestroyMostPowerfulUnlessReadyHouse covers Quicksand: each player who does
// not control a ready creature of the named house loses their most powerful
// creature; a player fielding a ready one is spared entirely.
func TestDestroyMostPowerfulUnlessReadyHouse(t *testing.T) {
	e := DestroyMostPowerfulUnlessReadyHouse{House: Untamed}
	if got := e.Text(); got != "destroy the most powerful creature controlled by "+
		"each player who does not control a ready Untamed creature" {
		t.Errorf("text = %q", got)
	}
	if (DestroyMostPowerfulUnlessReadyHouse{}).validate() == nil {
		t.Error("validate should reject an unset house")
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate with house set = %v", err)
	}

	g := NewGame("A", "B", 1)
	// P0 controls a ready Untamed creature, so it is spared entirely.
	readyUntamed := g.AddToBattleline(NewCard("ready", Untamed, Creature, Common, WithPower(2)), 0)
	bigP0 := g.AddToBattleline(testCreature("bigP0", 8), 0)
	// P1's only Untamed creature is exhausted, so P1 is not spared.
	exhaustedUntamed := g.AddToBattleline(
		NewCard("weary", Untamed, Creature, Common, WithPower(2)),
		1,
	)
	g.SetExhausted(exhaustedUntamed, true)
	bigP1 := g.AddToBattleline(testCreature("bigP1", 6), 1)
	smallP1 := g.AddToBattleline(testCreature("smallP1", 3), 1)

	e.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if !g.inPlay(readyUntamed) || !g.inPlay(bigP0) {
		t.Error("a player with a ready Untamed creature should be spared entirely")
	}
	if g.inPlay(bigP1) {
		t.Error("the most powerful creature of an unspared player should be destroyed")
	}
	if !g.inPlay(smallP1) || !g.inPlay(exhaustedUntamed) {
		t.Error("only the most powerful creature of an unspared player is destroyed")
	}
}

func TestDestroyChosenArtifact(t *testing.T) {
	g := NewGame("A", "B", 1)
	mine := g.AddArtifact(exAutocannon(), 0)
	theirs := g.AddArtifact(exAutocannon(), 1)
	ctx := &EffectContext{Resolver: g, Source: mine, Controller: 0}

	e := Destroy{Target: Target{Kind: TargetChosenArtifact}}
	if e.Text() != "destroy an artifact" {
		t.Errorf("text = %q", e.Text())
	}
	// The default chooser picks the first candidate (the controller's artifact).
	e.Resolve(ctx)
	if g.inPlay(mine) {
		t.Error("the chosen artifact should be destroyed and removed from play")
	}
	if !g.inPlay(theirs) {
		t.Error("the other artifact should be untouched")
	}
}

func TestDestroySamePower(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 3), 0) // chosen (candidates[0])
	strong := g.AddToBattleline(testCreature("strong", 5), 0)
	c := g.AddToBattleline(testCreature("c", 3), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := Destroy{Target: Target{Kind: TargetEachCreature}.Refine(SamePowerAsChosen)}
	if e.Text() != "choose a creature - destroy each creature with the same power as the chosen creature" {
		t.Errorf("text = %q", e.Text())
	}
	// The default chooser picks a (power 3); every power-3 creature is destroyed.
	e.Resolve(ctx)
	if g.inPlay(a) || g.inPlay(c) {
		t.Error("power-3 creatures should be destroyed, including the chosen one")
	}
	if !g.inPlay(strong) {
		t.Error("the power-5 creature should survive")
	}

	// A rejected choice destroys nothing (a second creature makes the choice real,
	// since a sole candidate would be auto-selected).
	g.AddToBattleline(testCreature("strong2", 5), 1)
	g.SetChooser(0, orderRejectChooser{})
	e.Resolve(ctx)
	if !g.inPlay(strong) {
		t.Error("rejecting the choice should destroy nothing")
	}
}

func TestDestroySamePowerEitherChosen(t *testing.T) {
	g := NewGame("A", "B", 1)
	fChosen := g.AddToBattleline(testCreature("fChosen", 3), 0)
	fShare := g.AddToBattleline(testCreature("fShare", 3), 0)       // shares friendly power
	fEnemyPow := g.AddToBattleline(testCreature("fEnemyPow", 4), 0) // shares enemy power
	fSurvive := g.AddToBattleline(testCreature("fSurvive", 6), 0)
	eChosen := g.AddToBattleline(testCreature("eChosen", 4), 1)
	eShare := g.AddToBattleline(testCreature("eShare", 4), 1) // shares enemy power
	eSurvive := g.AddToBattleline(testCreature("eSurvive", 7), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := Destroy{Target: Target{Kind: TargetEachCreature}.Refine(SamePowerAsEitherChosen)}
	want := "choose a friendly creature and an enemy creature - destroy each " +
		"creature with the same power as either of the chosen creatures"
	if got := e.Text(); got != want {
		t.Errorf("text = %q", got)
	}

	// A declined choice records no power, so nothing is destroyed.
	g.SetChooser(0, orderRejectChooser{})
	e.Resolve(ctx)
	for _, id := range []LocalID{fChosen, eChosen, fShare, eShare} {
		if !g.inPlay(id) {
			t.Fatal("declining the choices should destroy nothing")
		}
	}

	// Choose fChosen (power 3) then eChosen (power 4); the union of both power
	// brackets is destroyed, computed from the pre-destruction board.
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{fChosen, eChosen}})
	e.Resolve(ctx)
	for _, id := range []LocalID{fChosen, fShare, fEnemyPow, eChosen, eShare} {
		if g.inPlay(id) {
			t.Errorf("creature %d should have been destroyed", id)
		}
	}
	if !g.inPlay(fSurvive) {
		t.Error("the power-6 friendly matching neither chosen power should survive")
	}
	if !g.inPlay(eSurvive) {
		t.Error("the power-7 enemy matching neither chosen power should survive")
	}
}
