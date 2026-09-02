package engine

import "testing"

// TestTheirCreaturesThisWayCounts covers the two "... this way" counts that
// read the per-player tallies, and the tally itself as Destroy and
// PutFromPlay record it.
func TestCreaturesRemovedThisWayCounts(t *testing.T) {
	texts := []struct {
		count Count
		want  string
	}{
		{
			CreaturesDestroyedThisWay{Player: Controller},
			"creature they controlled that was destroyed this way",
		},
		{
			CreaturesShuffledIntoDeckThisWay{Player: Controller},
			"creature shuffled into their deck this way",
		},
		{
			CreaturesDestroyedThisWay{Player: Opponent},
			"creature your opponent controlled that was destroyed this way",
		},
		{
			CreaturesShuffledIntoDeckThisWay{Player: Opponent},
			"creature shuffled into your opponent's deck this way",
		},
	}
	for _, tc := range texts {
		if got := tc.count.CountText(); got != tc.want {
			t.Errorf("CountText = %q, want %q", got, tc.want)
		}
	}

	ctx := &EffectContext{Controller: 1}
	ctx.Produced.Destroyed = [2]int{4, 7}
	ctx.Produced.Moved = [2]int{5, 9}
	if got := (CreaturesDestroyedThisWay{Player: Controller}).Value(ctx); got != 7 {
		t.Errorf("destroyed Value = %d, want 7", got)
	}
	if got := (CreaturesShuffledIntoDeckThisWay{Player: Controller}).Value(ctx); got != 9 {
		t.Errorf("shuffled Value = %d, want 9", got)
	}
	if got := ctx.Produced.TotalDestroyed(); got != 11 {
		t.Errorf("TotalDestroyed = %d, want 11", got)
	}
}

// TestDestroyTalliesPerController checks Destroy splits its "this way"
// tally by the controller of each creature that actually left play.
func TestDestroyTalliesRemovalsPerController(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("mine1", Dis, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("mine2", Dis, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("theirs", Dis, Creature, Common, WithPower(2)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	Destroy{Target: Target{Kind: TargetEachCreature}.OfHouse(Dis)}.Resolve(ctx)

	if ctx.Produced.Destroyed != [2]int{2, 1} {
		t.Errorf("Destroyed = %v, want [2 1]", ctx.Produced.Destroyed)
	}
}

// TestPutFromPlayTalliesPerController checks PutFromPlay records the
// same per-controller tally when it shuffles creatures away.
func TestPutFromPlayTalliesRemovalsPerController(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("mine", Mars, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("theirs1", Mars, Creature, Common, WithPower(2)), 1)
	g.AddToBattleline(NewCard("theirs2", Mars, Creature, Common, WithPower(2)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	PutFromPlay{
		Target:      Target{Kind: TargetEachCreature}.OfHouse(Mars),
		Destination: ToDeckShuffled,
	}.Resolve(ctx)

	if ctx.Produced.Moved != [2]int{1, 2} {
		t.Errorf("Moved = %v, want [1 2]", ctx.Produced.Moved)
	}
}

// TestPutFromPlaySkipsCardsAlreadyGone checks a card a "Leaves Play:" ability
// destroyed mid-move is neither moved again nor counted in the tally.
func TestPutFromPlaySkipsCardsAlreadyGone(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("bomb", Mars, Creature, Common, WithPower(2),
		WithAbility(TriggerLeavesPlay, Destroy{
			Target: Target{Kind: TargetEachEnemyCreature},
		})), 0)
	g.AddToBattleline(NewCard("theirs1", Mars, Creature, Common, WithPower(2)), 1)
	g.AddToBattleline(NewCard("theirs2", Mars, Creature, Common, WithPower(2)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	PutFromPlay{
		Target:      Target{Kind: TargetEachCreature}.OfHouse(Mars),
		Destination: ToDeckShuffled,
	}.Resolve(ctx)

	if ctx.Produced.Moved != [2]int{1, 0} {
		t.Errorf("Moved = %v, want [1 0]: the destroyed creatures were not shuffled",
			ctx.Produced.Moved)
	}
}

// TestGainAemberEachPlayer checks an EachPlayer gain pays both players, and
// scales each player's share by the count read from their own side.
func TestGainAemberEachPlayer(t *testing.T) {
	gain := GainAember{
		Player: EachPlayer,
		Amount: 1,
		Per:    CreaturesDestroyedThisWay{Player: Controller},
	}
	if got, want := gain.Text(),
		"for each creature they controlled that was destroyed this way, "+
			"each player gains 1 Æmber"; got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}

	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	ctx.Produced.Destroyed = [2]int{3, 2}
	gain.Resolve(ctx)

	if got := g.State.Aember[0]; got != 3 {
		t.Errorf("controller Æmber = %d, want 3", got)
	}
	if got := g.State.Aember[1]; got != 2 {
		t.Errorf("opponent Æmber = %d, want 2", got)
	}
}

// TestController checks the port exposes which player a card in play answers to.
func TestController(t *testing.T) {
	g := NewGame("A", "B", 1)
	id := g.AddToBattleline(NewCard("theirs", Dis, Creature, Common, WithPower(2)), 1)
	if got := g.Controller(id); got != 1 {
		t.Errorf("Controller = %d, want 1", got)
	}
}
