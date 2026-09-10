package engine

import "testing"

// TestReturnItToHand covers each branch of the effect: its rendered text, the
// no-context no-op, recovering a destroyed creature from the discard pile, and
// returning a creature that is still in play.
func TestReturnItToHand(t *testing.T) {
	if got := (ReturnItToHand{}).Text(); got != "put it into its owner's hand" {
		t.Errorf("Text = %q, want %q", got, "put it into its owner's hand")
	}

	// No card in context: the effect does nothing.
	g := started(t)
	(ReturnItToHand{}).Resolve(&EffectContext{Resolver: g})

	// A destroyed creature (in the discard) is recovered to its owner's hand.
	dead := g.AddToBattleline(testCreature("dead", 3), 1)
	g.DestroyEach(0, []LocalID{dead})
	if g.inPlay(dead) {
		t.Fatal("dead should have left play")
	}
	(ReturnItToHand{}).Resolve(&EffectContext{Resolver: g, It: dead, HasIt: true})
	if !handContains(g, 1, dead) {
		t.Error("destroyed creature should be recovered to its owner's hand")
	}

	// A creature still in play is returned straight from the battleline.
	live := g.AddToBattleline(testCreature("live", 3), 1)
	(ReturnItToHand{}).Resolve(&EffectContext{Resolver: g, It: live, HasIt: true})
	if g.inPlay(live) {
		t.Fatal("live should have left play")
	}
	if !handContains(g, 1, live) {
		t.Error("in-play creature should be returned to its owner's hand")
	}
}

func handContains(g *Game, player int, id LocalID) bool {
	for _, h := range g.Hand(player) {
		if h == id {
			return true
		}
	}
	return false
}
