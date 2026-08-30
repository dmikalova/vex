package engine

import "testing"

func TestMoveFromDiscardToHand(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.Register(testCreature("c", 3), 0)
	g.State.Discard[0].add(c)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := PutFromDiscard{Destination: ToHand}
	if e.Text() != "put a card from your discard pile into your hand" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if len(g.Hand(0)) != 1 || g.Hand(0)[0] != c {
		t.Errorf("hand = %v, want [%d]", g.Hand(0), c)
	}
	if len(g.Discard(0)) != 0 {
		t.Error("the card should have left the discard pile")
	}

	// Empty discard: no candidate, so nothing happens.
	e.Resolve(ctx)
	if len(g.Hand(0)) != 1 {
		t.Error("an empty discard should return nothing")
	}
}

func TestPutFromDiscardByTrait(t *testing.T) {
	g := NewGame("A", "B", 1)
	horseman := g.Register(NewCard("Rider", Sanctum, Creature, Common, WithPower(5), WithTraits("Horseman")), 0)
	other := g.Register(NewCard("Squire", Sanctum, Creature, Common, WithPower(3), WithTraits("Human")), 0)
	g.State.Discard[0].add(horseman)
	g.State.Discard[0].add(other)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := PutFromDiscard{Type: Creature, Trait: "Horseman", All: true, Destination: ToHand}
	if e.Text() != "put each Horseman trait creature from your discard pile into your hand" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if !g.State.Hand[0].contains(horseman) {
		t.Error("the Horseman creature should return to hand")
	}
	if g.State.Hand[0].contains(other) {
		t.Error("the non-Horseman creature should stay in the discard pile")
	}
}

func TestMoveFromDiscardAll(t *testing.T) {
	g := NewGame("A", "B", 1)
	dis1 := g.Register(NewCard("d1", Dis, Creature, Common, WithPower(2)), 0)
	dis2 := g.Register(NewCard("d2", Dis, Creature, Common, WithPower(2)), 0)
	sanc := g.Register(NewCard("s", Sanctum, Creature, Common, WithPower(2)), 0)
	act := g.Register(NewCard("a", Dis, Tactic, Common), 0) // Dis but not a creature
	for _, id := range []LocalID{dis1, dis2, sanc, act} {
		g.State.Discard[0].add(id)
	}
	ctx := &EffectContext{Resolver: g, Controller: 0, ChosenHouse: Dis}

	e := PutFromDiscard{Type: Creature, Destination: ToHand, All: true, OfChosenHouse: true}
	if e.Text() != "put each creature of the chosen house from your discard pile into your hand" {
		t.Errorf("text = %q", e.Text())
	}
	if e.validate() != nil {
		t.Errorf("validate(All) = %v, want nil", e.validate())
	}
	e.Resolve(ctx)

	// Both Dis creatures return to hand; the Sanctum creature and the Dis action stay.
	hand := g.Hand(0)
	if len(hand) != 2 || !containsID(hand, dis1) || !containsID(hand, dis2) {
		t.Errorf("hand = %v, want [%d %d]", hand, dis1, dis2)
	}
	d := g.Discard(0)
	if len(d) != 2 || !containsID(d, sanc) || !containsID(d, act) {
		t.Errorf("discard = %v, want [%d %d]", d, sanc, act)
	}
}

func TestReturnCreatureFromDiscardToDeck(t *testing.T) {
	g := NewGame("A", "B", 1)
	act := g.Register(NewCard("act", Brobnar, Tactic, Common), 0) // non-creature
	g.State.Discard[0].add(act)
	crea := g.Register(testCreature("crea", 3), 0)
	g.State.Discard[0].add(crea)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := PutFromDiscard{Type: Creature, Destination: ToTopOfDeck}
	if e.Text() != "put a creature from your discard pile on top of your deck" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	// The action is skipped (not a creature); the creature goes to the deck top.
	if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != crea {
		t.Errorf("deck top = %v, want the creature %d", g.State.Deck[0].IDs[:g.State.Deck[0].Count], crea)
	}
	if d := g.Discard(0); len(d) != 1 || d[0] != act {
		t.Errorf("discard = %v, want just the action %d", d, act)
	}
}

func TestDiscardHand(t *testing.T) {
	g := NewGame("A", "B", 1)
	// Opponent (player 1) hand: a Mars creature, a Mars action, a Sanctum creature.
	marsCreature := g.AddToHand(NewCard("mc", Mars, Creature, Common, WithPower(2)), 1)
	marsAction := g.AddToHand(NewCard("ma", Mars, Tactic, Common), 1)
	sanctumCreature := g.AddToHand(NewCard("sc", Sanctum, Creature, Common, WithPower(2)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0, ChosenHouse: Mars}

	e := DiscardHand{Player: Opponent, CreaturesOnly: true, OfChosenHouse: true}
	if e.Text() != "discard each creature of the chosen house from your opponent's hand" {
		t.Errorf("text = %q", e.Text())
	}
	if plain := (DiscardHand{Player: Controller}).Text(); plain != "discard each card from your hand" {
		t.Errorf("plain text = %q", plain)
	}

	e.Resolve(ctx)
	if g.State.Hand[1].contains(marsCreature) {
		t.Error("the Mars creature should be discarded")
	}
	if !g.State.Hand[1].contains(marsAction) {
		t.Error("the Mars action is not a creature and should stay")
	}
	if !g.State.Hand[1].contains(sanctumCreature) {
		t.Error("the Sanctum creature is the wrong house and should stay")
	}
	if !g.State.Discard[1].contains(marsCreature) {
		t.Error("the Mars creature should be in the discard pile")
	}

	// Discarding a card no longer in the hand is a no-op.
	g.DiscardCardFromHand(1, marsCreature)
	if n := g.State.Discard[1].Count; n != 1 {
		t.Errorf("discard count = %d, want 1 (no double-discard)", n)
	}
}

func TestDiscardRandom(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToHand(NewCard("a", Mars, Tactic, Common), 1)
	b := g.AddToHand(NewCard("b", Mars, Tactic, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DiscardRandom{Player: Opponent}
	if e.Text() != "your opponent discards a random card from their hand" {
		t.Errorf("text = %q", e.Text())
	}
	if self := (DiscardRandom{Player: Controller}).Text(); self != "discard a random card from your hand" {
		t.Errorf("self text = %q", self)
	}
	if (DiscardRandom{}).validate() == nil {
		t.Error("unset player should be invalid")
	}
	if (DiscardRandom{Player: Opponent}).validate() != nil {
		t.Error("set player should be valid")
	}

	e.Resolve(ctx)
	if g.State.Hand[1].Count != 1 {
		t.Errorf("hand count = %d, want 1 after one discard", g.State.Hand[1].Count)
	}
	if g.State.Discard[1].Count != 1 {
		t.Errorf("discard count = %d, want 1", g.State.Discard[1].Count)
	}
	// The discarded card is one of the two; the other remains.
	if got := g.State.Hand[1].contains(a) == g.State.Hand[1].contains(b); got {
		t.Error("exactly one of the two cards should remain in hand")
	}

	// An empty hand is a no-op.
	g.DiscardRandomFromHand(1)
	g.DiscardRandomFromHand(1) // hand now empty
	g.DiscardRandomFromHand(1)
	if g.State.Discard[1].Count != 2 {
		t.Errorf("discard count = %d, want 2 (empty-hand discards are no-ops)", g.State.Discard[1].Count)
	}
}
