package engine

import "testing"

// TestShuffleIntoDeckFromZonesText covers the rendered phrase with and
// without a house filter. An in-play source renders as the bare "play" and puts
// the side scope on the object instead ("friendly creatures"), because play
// belongs to neither player and so takes no possessive.
func TestShuffleIntoDeckFromZonesText(t *testing.T) {
	if got := (ShuffleIntoDeck{Player: Controller, From: songOfSpringZones, Selection: Chosen{Type: Creature, Optional: true}, Quantity: AnyNumber{}}).Text(); got !=
		"shuffle any number of friendly creatures from your hand, "+
			"your discard pile, or play into your deck" {
		t.Errorf("Text = %q", got)
	}
	if got := (ShuffleIntoDeck{Player: Controller, From: songOfSpringZones, Selection: Chosen{House: namedHouse(Untamed), Type: Creature, Optional: true}, Quantity: AnyNumber{}}).Text(); got !=
		"shuffle any number of friendly Untamed creatures from your hand, "+
			"your discard pile, or play into your deck" {
		t.Errorf("Text = %q", got)
	}
	if err := (ShuffleIntoDeck{Player: Controller, From: songOfSpringZones, Selection: Chosen{Type: Creature, Optional: true}, Quantity: AnyNumber{}}).validate(); err != nil {
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
		Quantity:  AnyNumber{},
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
		Quantity:  AnyNumber{},
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

// TestShuffleIntoDeckFromPlayFoldsTheWholeBoard pins the Timequake shape: an
// untyped Each over an in-play source reaches every friendly card — creature,
// artifact, and the upgrade attached to a creature — and shuffles each home
// rather than shedding the upgrade to the discard pile. Enemy cards stay put.
func TestShuffleIntoDeckFromPlayFoldsTheWholeBoard(t *testing.T) {
	g := NewGame("A", "B", 1)
	creature := g.AddToBattleline(NewCard("creature", Brobnar, Creature, Common, WithPower(4)), 0)
	upgrade := g.Register(NewCard("boon", Brobnar, Upgrade, Common), 0)
	g.AttachUpgrade(creature, upgrade)
	artifact := g.AddArtifact(NewCard("relic", Brobnar, Artifact, Common), 0)
	enemy := g.AddToBattleline(NewCard("enemy", Brobnar, Creature, Common, WithPower(4)), 1)

	ctx := &EffectContext{Resolver: g, Controller: 0}
	e := ShuffleIntoDeck{Player: Controller, From: []Zone{InPlay}, Selection: Each{}}
	if got := e.Text(); got != "shuffle each friendly card from play into your deck" {
		t.Errorf("Text = %q", got)
	}
	e.Resolve(ctx)

	for _, id := range []LocalID{creature, upgrade, artifact} {
		if !g.State.Deck[0].contains(id) {
			t.Errorf("%s should have been shuffled into the deck", g.Name(id))
		}
		if g.State.Discard[0].contains(id) {
			t.Errorf("%s was discarded instead of shuffled home", g.Name(id))
		}
	}
	if !g.inPlay(enemy) {
		t.Error("the enemy creature should have been left in play")
	}
	if ctx.Produced.Moved != [2]int{3, 0} {
		t.Errorf("Moved = %v, want all 3 friendly cards tallied", ctx.Produced.Moved)
	}
	if got := (CardsShuffledIntoDeck{}).Value(ctx); got != 3 {
		t.Errorf("CardsShuffledIntoDeck.Value = %d, want 3", got)
	}
	if err := g.InvariantError(); err != nil {
		t.Errorf("state should stay sound: %v", err)
	}
}

// TestShuffleIntoDeckFromPlayEnemySide pins that the side scope an in-play source
// puts on the object follows Player: reaching the opponent's board reads "enemy".
func TestShuffleIntoDeckFromPlayEnemySide(t *testing.T) {
	e := ShuffleIntoDeck{Player: Opponent, From: []Zone{InPlay}, Selection: Each{Type: Creature}}
	if got := e.Text(); got != "shuffle each enemy creature from play into your opponent's deck" {
		t.Errorf("Text = %q", got)
	}
}
