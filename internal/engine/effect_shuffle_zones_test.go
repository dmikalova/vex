package engine

import "testing"

// TestShuffleChosenCreaturesFromZonesText covers the rendered phrase with and
// without a house filter.
func TestShuffleChosenCreaturesFromZonesText(t *testing.T) {
	if got := (ShuffleChosenCreaturesFromZones{}).Text(); got !=
		"shuffle any number of friendly creatures from your hand, "+
			"discard pile, or battleline into your deck" {
		t.Errorf("Text = %q", got)
	}
	if got := (ShuffleChosenCreaturesFromZones{House: Untamed}).Text(); got !=
		"shuffle any number of friendly Untamed creatures from your hand, "+
			"discard pile, or battleline into your deck" {
		t.Errorf("Text = %q", got)
	}
	if err := (ShuffleChosenCreaturesFromZones{}).validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
}

// TestShuffleChosenCreaturesFromZonesResolve checks a creature is drawn from each
// of the three zones, the house filter spares a non-matching creature, and the
// controller can decline the rest.
func TestShuffleChosenCreaturesFromZonesResolve(t *testing.T) {
	g := NewGame("A", "B", 1)
	inHand := g.Register(NewCard("inHand", Untamed, Creature, Common, WithPower(3)), 0)
	g.State.Hand[0].add(inHand)
	inDiscard := g.Register(NewCard("inDiscard", Untamed, Creature, Common, WithPower(3)), 0)
	g.State.Discard[0].add(inDiscard)
	onBoard := g.AddToBattleline(NewCard("onBoard", Untamed, Creature, Common, WithPower(3)), 0)
	offHouse := g.Register(NewCard("offHouse", Brobnar, Creature, Common, WithPower(3)), 0)
	g.State.Hand[0].add(offHouse)
	notCreature := g.Register(NewCard("tactic", Untamed, Tactic, Common), 0)
	g.State.Discard[0].add(notCreature)

	g.SetChooser(0, &declineAfterChooser{ids: []LocalID{inHand, inDiscard, onBoard}})

	ctx := &EffectContext{Resolver: g, Controller: 0}
	ShuffleChosenCreaturesFromZones{House: Untamed}.Resolve(ctx)

	for _, id := range []LocalID{inHand, inDiscard, onBoard} {
		if !g.State.Deck[0].contains(id) {
			t.Errorf("%s should have been shuffled into the deck", g.Name(id))
		}
	}
	if !g.State.Hand[0].contains(offHouse) {
		t.Error("the Brobnar creature should have been spared")
	}
}

// TestShuffleChosenCreaturesFromZonesDeclineImmediately checks the controller can
// decline the very first pick, leaving every creature where it was.
func TestShuffleChosenCreaturesFromZonesDeclineImmediately(t *testing.T) {
	g := NewGame("A", "B", 1)
	inHand := g.Register(NewCard("inHand", Untamed, Creature, Common, WithPower(3)), 0)
	g.State.Hand[0].add(inHand)
	g.SetChooser(0, &declineAfterChooser{})

	ctx := &EffectContext{Resolver: g, Controller: 0}
	ShuffleChosenCreaturesFromZones{}.Resolve(ctx)

	if !g.State.Hand[0].contains(inHand) {
		t.Error("declining should leave the hand creature in place")
	}
}

// TestShuffleFromHandIntoDeckNarratesOutsideBatch checks the standalone narration
// path when no shuffle batch is open.
func TestShuffleFromHandIntoDeckNarratesOutsideBatch(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetRecording(true)
	id := g.Register(NewCard("card", Untamed, Creature, Common, WithPower(3)), 0)
	g.State.Hand[0].add(id)

	g.ShuffleFromHandIntoDeck(id)

	if !g.State.Deck[0].contains(id) || g.State.Hand[0].contains(id) {
		t.Error("the card should have moved from hand to deck")
	}
}
