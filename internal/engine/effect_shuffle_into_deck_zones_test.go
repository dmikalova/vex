package engine

import "testing"

// TestShuffleIntoDeckFromZonesText covers the rendered phrase with and
// without a house filter.
func TestShuffleIntoDeckFromZonesText(t *testing.T) {
	if got := (ShuffleIntoDeck{Player: Controller, From: songOfSpringZones, Selection: Chosen{Type: Creature, Optional: true}, AnyNumber: true}).Text(); got !=
		"shuffle any number of creatures from your hand, "+
			"your discard pile, or in play into your deck" {
		t.Errorf("Text = %q", got)
	}
	if got := (ShuffleIntoDeck{Player: Controller, From: songOfSpringZones, Selection: Chosen{House: namedHouse(Untamed), Type: Creature, Optional: true}, AnyNumber: true}).Text(); got !=
		"shuffle any number of Untamed creatures from your hand, "+
			"your discard pile, or in play into your deck" {
		t.Errorf("Text = %q", got)
	}
	if err := (ShuffleIntoDeck{Player: Controller, From: songOfSpringZones, Selection: Chosen{Type: Creature, Optional: true}, AnyNumber: true}).validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
}

// TestShuffleIntoDeckFromZonesResolve checks a creature is drawn from each
// of the three zones, the house filter spares a non-matching creature, and the
// controller can decline the rest.
func TestShuffleIntoDeckFromZonesResolve(t *testing.T) {
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
	ShuffleIntoDeck{
		Player:    Controller,
		From:      songOfSpringZones,
		Selection: Chosen{House: namedHouse(Untamed), Type: Creature, Optional: true},
		AnyNumber: true,
	}.Resolve(
		ctx,
	)

	for _, id := range []LocalID{inHand, inDiscard, onBoard} {
		if !g.State.Deck[0].contains(id) {
			t.Errorf("%s should have been shuffled into the deck", g.Name(id))
		}
	}
	if !g.State.Hand[0].contains(offHouse) {
		t.Error("the Brobnar creature should have been spared")
	}
}

// TestShuffleIntoDeckFromZonesDeclineImmediately checks the controller can
// decline the very first pick, leaving every creature where it was.
func TestShuffleIntoDeckFromZonesDeclineImmediately(t *testing.T) {
	g := NewGame("A", "B", 1)
	inHand := g.Register(NewCard("inHand", Untamed, Creature, Common, WithPower(3)), 0)
	g.State.Hand[0].add(inHand)
	g.SetChooser(0, &declineAfterChooser{})

	ctx := &EffectContext{Resolver: g, Controller: 0}
	ShuffleIntoDeck{
		Player:    Controller,
		From:      songOfSpringZones,
		Selection: Chosen{Type: Creature, Optional: true},
		AnyNumber: true,
	}.Resolve(
		ctx,
	)

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

// songOfSpringZones is the three-zone source set Song of Spring reaches across.
var songOfSpringZones = []Zone{Hand, Discard, InPlay}

// TestShuffleIntoDeckTalliesByOwner pins that a shuffled card is counted under
// its owner, not under the player whose zone it left: a controlled enemy creature
// returns to its owner's deck, so it must refill the opponent's tally and not
// feed the controller's "draw a card for each card shuffled this way".
func TestShuffleIntoDeckTalliesByOwner(t *testing.T) {
	g := NewGame("A", "B", 1)
	mine := g.AddToBattleline(NewCard("mine", Untamed, Creature, Common, WithPower(3)), 0)
	theirs := g.Register(NewCard("theirs", Untamed, Creature, Common, WithPower(3)), 1)
	g.State.Battleline[0].add(theirs)

	ctx := &EffectContext{Resolver: g, Controller: 0}
	ShuffleIntoDeck{
		Player:    Controller,
		From:      []Zone{InPlay},
		Selection: Each{Type: Creature},
	}.Resolve(ctx)

	if !g.State.Deck[0].contains(mine) || !g.State.Deck[1].contains(theirs) {
		t.Fatal("each creature should return to its own owner's deck")
	}
	if ctx.Produced.Moved != [2]int{1, 1} {
		t.Errorf("Moved = %v, want one card tallied per owner", ctx.Produced.Moved)
	}
	if err := g.InvariantError(); err != nil {
		t.Errorf("state should stay sound: %v", err)
	}
}

// TestShuffleIntoDeckValidatePlayer checks an unset Player is rejected (ADR 0010).
func TestShuffleIntoDeckValidatePlayer(t *testing.T) {
	if (ShuffleIntoDeck{From: []Zone{Discard}, Selection: Chosen{}}).validate() == nil {
		t.Error("an unset Player should be rejected")
	}
}
