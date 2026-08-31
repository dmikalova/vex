package engine

import "testing"

// midThenFirst picks mid for the first creature choice, then the first candidate
// for any later choice (the neighbor pick).
type midThenFirst struct {
	FirstChooser
	mid LocalID
}

func (c midThenFirst) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	for _, x := range cands {
		if x == c.mid {
			return x, true
		}
	}
	return cands[0], true
}

func TestDamageIfSurvives(t *testing.T) {
	g := NewGame("A", "B", 1)
	survivor := g.AddToBattleline(testCreature("s", 5), 1)
	g.AddToHand(testCreature("h1", 1), 1)
	g.AddToHand(testCreature("h2", 1), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DamageIfSurvives{
		Amount: 3,
		Target: Target{Kind: TargetChosenCreature},
		Then:   DiscardRandomFromHand{Player: ItsOwner},
	}
	if e.Text() != "deal 3 damage to a creature. If it is not destroyed, its owner discards a random card from their hand" {
		t.Errorf("text = %q", e.Text())
	}
	if (DamageIfSurvives{Then: DiscardRandomFromHand{Player: Opponent}}).validate() == nil {
		t.Error("unset target should be invalid")
	}
	if (DamageIfSurvives{Target: Target{Kind: TargetChosenCreature}, Then: DiscardRandomFromHand{Player: Opponent}}).validate() != nil {
		t.Error("a set target with a valid follow-up should pass")
	}

	e.Resolve(ctx)
	if g.Damage(survivor) != 3 {
		t.Errorf("survivor damage = %d, want 3", g.Damage(survivor))
	}
	if g.State.Hand[1].Count != 1 {
		t.Errorf("owner hand = %d, want 1 (a card discarded)", g.State.Hand[1].Count)
	}

	// A creature that dies triggers no follow-up.
	g2 := NewGame("A", "B", 1)
	dead := g2.AddToBattleline(testCreature("d", 2), 1)
	g2.AddToHand(testCreature("keep", 1), 1)
	(DamageIfSurvives{Amount: 3, Target: Target{Kind: TargetChosenCreature}, Then: DiscardRandomFromHand{Player: ItsOwner}}).Resolve(
		&EffectContext{Resolver: g2, Controller: 0},
	)
	if g2.inPlay(dead) {
		t.Error("the creature should have been destroyed")
	}
	if g2.State.Hand[1].Count != 1 {
		t.Error("a destroyed creature's owner should not discard")
	}

	// No creature to target: nothing happens.
	g3 := NewGame("A", "B", 1)
	(DamageIfSurvives{Amount: 3, Target: Target{Kind: TargetChosenCreature}, Then: DiscardRandomFromHand{Player: ItsOwner}}).Resolve(
		&EffectContext{Resolver: g3, Controller: 0},
	)
}

func TestDamageCreatureAndNeighbor(t *testing.T) {
	e := DamageCreatureAndNeighbor{Amount: 3, NeighborAmount: 3}
	if e.Text() != "deal 3 damage to a creature and 3 damage to a neighbor of that creature" {
		t.Errorf("text = %q", e.Text())
	}

	// The default chooser picks the left flank; its only neighbor is the middle.
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 10), 1)
	mid := g.AddToBattleline(testCreature("mid", 10), 1)
	right := g.AddToBattleline(testCreature("right", 10), 1)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.Damage(left) != 3 || g.Damage(mid) != 3 || g.Damage(right) != 0 {
		t.Errorf("damage = %d/%d/%d, want 3/3/0", g.Damage(left), g.Damage(mid), g.Damage(right))
	}

	// A middle creature has two neighbors, so the controller picks one.
	g2 := NewGame("A", "B", 1)
	l2 := g2.AddToBattleline(testCreature("l", 10), 1)
	m2 := g2.AddToBattleline(testCreature("m", 10), 1)
	r2 := g2.AddToBattleline(testCreature("r", 10), 1)
	g2.SetChooser(0, midThenFirst{mid: m2})
	e.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.Damage(m2) != 3 || g2.Damage(l2) != 3 || g2.Damage(r2) != 0 {
		t.Errorf(
			"damage = %d/%d/%d, want mid 3 + one neighbor 3",
			g2.Damage(l2),
			g2.Damage(m2),
			g2.Damage(r2),
		)
	}

	// No creatures: nothing happens.
	g3 := NewGame("A", "B", 1)
	e.Resolve(&EffectContext{Resolver: g3, Controller: 0})
}

func TestOverwhelmed(t *testing.T) {
	if (Overwhelmed{}).CondText() != "if you are overwhelmed" {
		t.Errorf("cond text = %q", (Overwhelmed{}).CondText())
	}
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("o1", 2), 1)
	g.AddToBattleline(testCreature("o2", 2), 1)
	g.AddToBattleline(testCreature("m1", 2), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if !(Overwhelmed{}).Met(ctx) {
		t.Error("should be overwhelmed when the opponent controls more creatures")
	}
	g.AddToBattleline(testCreature("m2", 2), 0)
	if (Overwhelmed{}).Met(ctx) {
		t.Error("should not be overwhelmed at parity")
	}
}

func TestRepeatOnCondition(t *testing.T) {
	e := RepeatOnCondition{
		Do:   Destroy{Target: Target{Kind: TargetEachEnemyCreature}.Selector(LeastPowerful)},
		Cond: Overwhelmed{},
	}
	if e.Text() != "destroy the least powerful enemy creature -> if you are overwhelmed, repeat this effect" {
		t.Errorf("text = %q", e.Text())
	}
	if (RepeatOnCondition{Do: Destroy{}, Cond: Overwhelmed{}}).validate() == nil {
		t.Error("unset destroy target should be invalid")
	}

	// Opponent has 3, controller 1: destroys until no longer overwhelmed.
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("o1", 2), 1)
	g.AddToBattleline(testCreature("o2", 2), 1)
	g.AddToBattleline(testCreature("o3", 2), 1)
	g.AddToBattleline(testCreature("m1", 2), 0)
	RepeatOnCondition{
		Do:   Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
		Cond: Overwhelmed{},
	}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if len(g.Battleline(1)) != 1 {
		t.Errorf("opponent creatures = %d, want 1 (destroyed down to parity)", len(g.Battleline(1)))
	}

	// No enemy creatures: the gate stops the loop immediately.
	g2 := NewGame("A", "B", 1)
	g2.AddToBattleline(testCreature("m", 2), 0)
	RepeatOnCondition{
		Do:   Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
		Cond: Overwhelmed{},
	}.Resolve(
		&EffectContext{Resolver: g2, Controller: 0},
	)
}

func TestTakeControlTargeted(t *testing.T) {
	if got := (TakeControl{Duration: UntilThisLeavesPlay}).Text(); got != "take control of this creature until "+UpgradeName+" leaves play" {
		t.Errorf("host text = %q", got)
	}
	tgt := TakeControl{
		Target:   Target{Kind: TargetChosenEnemyCreature}.OnFlank(),
		Duration: UntilThisLeavesPlay,
	}
	if got := tgt.Text(); got != "take control of an enemy flank creature until "+SelfName+" leaves play" {
		t.Errorf("targeted text = %q", got)
	}
	if (TakeControl{Target: Target{Kind: TargetChosenEnemyCreature}, Duration: EndOfTurn}).validate() == nil {
		t.Error("only UntilThisLeavesPlay should be valid")
	}

	g := NewGame("A", "B", 1)
	harland := g.AddToBattleline(testCreature("harland", 1), 0)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	tgt.Resolve(&EffectContext{Resolver: g, Source: harland, Controller: 0})
	if g.controller(foe) != 0 {
		t.Errorf("controller of the seized creature = %d, want 0", g.controller(foe))
	}
}

func TestEnemyCreatureDestroyedReaction(t *testing.T) {
	if !EventEnemyCreatureDestroyed.isReaction() {
		t.Error("EventEnemyCreatureDestroyed should be a reaction point")
	}
	if EventEnemyCreatureDestroyed.clause() != "each time an enemy creature is destroyed" {
		t.Errorf("clause = %q", EventEnemyCreatureDestroyed.clause())
	}

	g := NewGame("A", "B", 1)
	g.AddLasting(EventEnemyCreatureDestroyed, actGainAember, 0, 1)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	g.destroyEach(0, []LocalID{foe})
	if g.State.Aember[0] != 1 {
		t.Errorf(
			"controller Æmber = %d, want 1 after an enemy creature was destroyed",
			g.State.Aember[0],
		)
	}
}
