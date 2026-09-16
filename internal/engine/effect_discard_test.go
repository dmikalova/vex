package engine

import "testing"

func TestMoveFromDiscardToHand(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.Register(testCreature("c", 3), 0)
	g.State.Discard[0].add(c)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := PutFromDiscard{Selection: Chosen{}, Destination: ToHand}
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

func TestPutFromDiscardBindsReturnedCard(t *testing.T) {
	// Bind leaves the returned card in context so a following effect can act on it.
	g := NewGame("A", "B", 1)
	c := g.Register(testCreature("c", 3), 0)
	g.State.Discard[0].add(c)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	PutFromDiscard{Selection: Chosen{}, Destination: ToHand, Bind: true}.Resolve(ctx)
	if !ctx.HasIt || ctx.It != c {
		t.Errorf("ctx.It = %v (HasIt %v), want %d bound", ctx.It, ctx.HasIt, c)
	}

	// Nothing to return binds nothing.
	ctx2 := &EffectContext{Resolver: g, Controller: 0}
	PutFromDiscard{Selection: Chosen{}, Destination: ToHand, Bind: true}.Resolve(ctx2)
	if ctx2.HasIt {
		t.Error("an empty discard should leave ctx.It unset")
	}
}

func TestPutFromDiscardByTrait(t *testing.T) {
	g := NewGame("A", "B", 1)
	horseman := g.Register(
		NewCard(
			"Rider",
			Sanctum,
			Creature,
			Common,
			WithPower(5),
			WithTraits(Horseman),
		),
		0,
	)
	other := g.Register(
		NewCard(
			"Squire",
			Sanctum,
			Creature,
			Common,
			WithPower(3),
			WithTraits(Human),
		),
		0,
	)
	g.State.Discard[0].add(horseman)
	g.State.Discard[0].add(other)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := PutFromDiscard{
		Selection:   Each{Type: Creature, Trait: Horseman},
		Destination: ToHand,
	}
	if e.Text() != "put each Horseman creature from your discard pile into your hand" {
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

func TestPutFromDiscardByTraitChoose(t *testing.T) {
	g := NewGame("A", "B", 1)
	horseman := g.Register(
		NewCard(
			"Rider",
			Sanctum,
			Creature,
			Common,
			WithPower(5),
			WithTraits(Horseman),
		),
		0,
	)
	other := g.Register(
		NewCard(
			"Squire",
			Sanctum,
			Creature,
			Common,
			WithPower(3),
			WithTraits(Human),
		),
		0,
	)
	g.State.Discard[0].add(horseman)
	g.State.Discard[0].add(other)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Not All: the non-Horseman card is filtered out of the candidates, leaving
	// only the Horseman for the controller to choose.
	e := PutFromDiscard{Selection: Chosen{Type: Creature, Trait: Horseman}, Destination: ToHand}
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
	act := g.Register(
		NewCard("a", Dis, Tactic, Common),
		0,
	) // Dis but not a creature
	for _, id := range []LocalID{dis1, dis2, sanc, act} {
		g.State.Discard[0].add(id)
	}
	ctx := &EffectContext{Resolver: g, Controller: 0, ChosenHouse: Dis}

	e := PutFromDiscard{
		Selection:   Each{Type: Creature, House: chosenHouse},
		Destination: ToHand,
	}
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

func TestPutFromDiscardByName(t *testing.T) {
	g := NewGame("A", "B", 1)
	bind1 := g.Register(NewCard("Ortannu's Binding", Dis, Tactic, Common), 0)
	bind2 := g.Register(NewCard("Ortannu's Binding", Dis, Tactic, Common), 0)
	other := g.Register(NewCard("Gateway to Dis", Dis, Tactic, Common), 0)
	for _, id := range []LocalID{bind1, bind2, other} {
		g.State.Discard[0].add(id)
	}
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := PutFromDiscard{
		Selection:   Each{Name: "Ortannu's Binding"},
		Destination: ToHand,
	}
	if e.Text() != "put each Ortannu's Binding from your discard pile into your hand" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	hand := g.Hand(0)
	if len(hand) != 2 || !containsID(hand, bind1) || !containsID(hand, bind2) {
		t.Errorf("hand = %v, want both Bindings", hand)
	}
	if d := g.Discard(0); len(d) != 1 || d[0] != other {
		t.Errorf("discard = %v, want just the differently named card", d)
	}
	if ctx.Produced.Returned != 2 {
		t.Errorf("Returned = %d, want 2", ctx.Produced.Returned)
	}
}

func TestPutFromDiscardByNameChoose(t *testing.T) {
	g := NewGame("A", "B", 1)
	bind := g.Register(NewCard("Ortannu's Binding", Dis, Tactic, Common), 0)
	other := g.Register(NewCard("Gateway to Dis", Dis, Tactic, Common), 0)
	g.State.Discard[0].add(bind)
	g.State.Discard[0].add(other)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Not All: the differently named card is filtered out of the candidates,
	// leaving only the Binding for the controller to choose.
	e := PutFromDiscard{Selection: Chosen{Name: "Ortannu's Binding"}, Destination: ToHand}
	e.Resolve(ctx)
	if !g.State.Hand[0].contains(bind) {
		t.Error("the named card should return to hand")
	}
	if g.State.Hand[0].contains(other) {
		t.Error("the differently named card should stay in the discard pile")
	}
}

func TestReturnCreatureFromDiscardToDeck(t *testing.T) {
	g := NewGame("A", "B", 1)
	act := g.Register(
		NewCard("act", Brobnar, Tactic, Common),
		0,
	) // non-creature
	g.State.Discard[0].add(act)
	crea := g.Register(testCreature("crea", 3), 0)
	g.State.Discard[0].add(crea)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := PutFromDiscard{Selection: Chosen{Type: Creature}, Destination: ToTopOfDeck}
	if e.Text() != "put a creature from your discard pile on top of your deck" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	// The action is skipped (not a creature); the creature goes to the deck top.
	if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != crea {
		t.Errorf(
			"deck top = %v, want the creature %d",
			g.State.Deck[0].IDs[:g.State.Deck[0].Count],
			crea,
		)
	}
	if d := g.Discard(0); len(d) != 1 || d[0] != act {
		t.Errorf("discard = %v, want just the action %d", d, act)
	}
}

func TestPutFromDiscardTypeOrTrait(t *testing.T) {
	g := NewGame("A", "B", 1)
	upgrade := g.Register(NewCard("chip", StarAlliance, Upgrade, Common), 0)
	robot := g.Register(
		NewCard(
			"droid",
			StarAlliance,
			Creature,
			Common,
			WithPower(3),
			WithTraits(Robot),
		),
		0,
	)
	human := g.Register(
		NewCard(
			"pilot",
			StarAlliance,
			Creature,
			Common,
			WithPower(3),
			WithTraits(Human),
		),
		0,
	)
	for _, id := range []LocalID{upgrade, robot, human} {
		g.State.Discard[0].add(id)
	}
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := PutFromDiscard{
		Selection:   Each{Type: Upgrade, Or: []CardFilter{{Trait: Robot}}},
		Destination: ToHand,
	}
	if e.Text() != "put each upgrade or Robot card from your discard pile into your hand" {
		t.Errorf("text = %q", e.Text())
	}
	choose := PutFromDiscard{
		Selection:   Chosen{Type: Upgrade, Or: []CardFilter{{Trait: Robot}}},
		Destination: ToHand,
	}
	if choose.Text() != "put an upgrade or Robot card from your discard pile into your hand" {
		t.Errorf("choose text = %q", choose.Text())
	}
	e.Resolve(ctx)

	// The upgrade (matches Type) and the Robot creature (matches OrTrait) return;
	// the Human creature matches neither and stays in the discard pile.
	hand := g.Hand(0)
	if len(hand) != 2 || !containsID(hand, upgrade) ||
		!containsID(hand, robot) {
		t.Errorf(
			"hand = %v, want the upgrade %d and the Robot %d",
			hand,
			upgrade,
			robot,
		)
	}
	if d := g.Discard(0); len(d) != 1 || d[0] != human {
		t.Errorf("discard = %v, want just the Human creature %d", d, human)
	}
}

// May{PutFromDiscard} with no matching card in the discard is vacuous, so the
// controller is never asked (Chief Engineer Walls with no upgrade or Robot in the
// discard prompts for nothing).
func TestPutFromDiscardVacuousUnderMay(t *testing.T) {
	g := NewGame("A", "B", 1)
	human := g.Register(
		NewCard(
			"pilot",
			StarAlliance,
			Creature,
			Common,
			WithPower(3),
			WithTraits(Human),
		),
		0,
	)
	g.State.Discard[0].add(human)
	ch := &optionRecorder{}
	g.SetChooser(0, ch)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	may := May{
		Do: PutFromDiscard{
			Selection:   Chosen{Type: Upgrade, Or: []CardFilter{{Trait: Robot}}},
			Destination: ToHand,
		},
	}
	may.Resolve(ctx)
	if ch.asked != 0 {
		t.Errorf("prompts with nothing to recover = %d, want 0", ch.asked)
	}
	if len(g.Hand(0)) != 0 {
		t.Error("nothing should have moved to hand")
	}

	// Add a Robot: the choice becomes real, so the card is offered and recovered.
	robot := g.Register(
		NewCard(
			"droid",
			StarAlliance,
			Creature,
			Common,
			WithPower(3),
			WithTraits(Robot),
		),
		0,
	)
	g.State.Discard[0].add(robot)
	may.Resolve(ctx)
	if !g.State.Hand[0].contains(robot) {
		t.Error("the Robot should be recovered once the choice is real")
	}
}

func TestDiscardFromHandEach(t *testing.T) {
	g := NewGame("A", "B", 1)
	// Opponent (player 1) hand: a Mars creature, a Mars action, a Sanctum creature.
	marsCreature := g.AddToHand(
		NewCard("mc", Mars, Creature, Common, WithPower(2)),
		1,
	)
	marsAction := g.AddToHand(NewCard("ma", Mars, Tactic, Common), 1)
	sanctumCreature := g.AddToHand(
		NewCard("sc", Sanctum, Creature, Common, WithPower(2)),
		1,
	)
	ctx := &EffectContext{Resolver: g, Controller: 0, ChosenHouse: Mars}

	// Each from an opponent's hand reads as a controller-directed discard, not
	// "your opponent discards".
	e := DiscardCard{
		Player:    Opponent,
		Zone:      Hand,
		Selection: Each{Type: Creature, House: chosenHouse},
	}
	if e.Text() != "discard each creature of the chosen house from your opponent's hand" {
		t.Errorf("text = %q", e.Text())
	}
	if plain := (DiscardCard{Player: Controller, Zone: Hand, Selection: Each{}}).Text(); plain != "discard each card from your hand" {
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

func TestDiscardRandomFromHand(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToHand(NewCard("a", Mars, Tactic, Common), 1)
	b := g.AddToHand(NewCard("b", Mars, Tactic, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DiscardCard{Player: Opponent, Zone: Hand, Selection: Random{}}
	if e.Text() != "your opponent discards a random card from their hand" {
		t.Errorf("text = %q", e.Text())
	}
	if self := (DiscardCard{Player: Controller, Zone: Hand, Selection: Random{}}).Text(); self != "discard a random card from your hand" {
		t.Errorf("self text = %q", self)
	}
	if (DiscardCard{Zone: Hand, Selection: Random{}}).validate() == nil {
		t.Error("unset player should be invalid")
	}
	if (DiscardCard{Player: Opponent, Zone: Hand, Selection: Random{}}).validate() != nil {
		t.Error("set player should be valid")
	}

	e.Resolve(ctx)
	if g.State.Hand[1].Count != 1 {
		t.Errorf(
			"hand count = %d, want 1 after one discard",
			g.State.Hand[1].Count,
		)
	}
	if g.State.Discard[1].Count != 1 {
		t.Errorf("discard count = %d, want 1", g.State.Discard[1].Count)
	}
	// The discarded card is one of the two; the other remains.
	if got := g.State.Hand[1].contains(a) == g.State.Hand[1].contains(b); got {
		t.Error("exactly one of the two cards should remain in hand")
	}

	// An empty hand is a no-op.
	e.Resolve(ctx)
	e.Resolve(ctx) // hand now empty
	e.Resolve(ctx)
	if g.State.Discard[1].Count != 2 {
		t.Errorf(
			"discard count = %d, want 2 (empty-hand discards are no-ops)",
			g.State.Discard[1].Count,
		)
	}
}

func TestDiscardRandomFromHandAmount(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToHand(NewCard("a", Mars, Tactic, Common), 0)
	g.AddToHand(NewCard("b", Mars, Tactic, Common), 0)
	g.AddToHand(NewCard("c", Mars, Tactic, Common), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DiscardCard{
		Player:    Controller,
		Zone:      Hand,
		Selection: Random{},
		Amount:    2,
	}
	if got := e.Text(); got != "discard 2 random cards from your hand" {
		t.Errorf("text = %q", got)
	}
	e.Resolve(ctx)
	if g.State.Hand[0].Count != 1 {
		t.Errorf(
			"hand count = %d, want 1 after discarding 2",
			g.State.Hand[0].Count,
		)
	}
	if g.State.Discard[0].Count != 2 {
		t.Errorf("discard count = %d, want 2", g.State.Discard[0].Count)
	}
}

func TestDiscardFromArchives(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToArchives(NewCard("a", Mars, Tactic, Common), 1)
	b := g.AddToArchives(NewCard("b", Mars, Tactic, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DiscardCard{Player: Opponent, Zone: Archives, Selection: Random{}}
	if e.Text() != "your opponent discards a random card from their archives" {
		t.Errorf("text = %q", e.Text())
	}
	if self := (DiscardCard{Player: Controller, Zone: Archives, Selection: Random{}}).Text(); self != "discard a random card from your archives" {
		t.Errorf("self text = %q", self)
	}
	if owner := (DiscardCard{Player: ItsOwner, Zone: Archives, Selection: Random{}}).Text(); owner != "its owner discards a random card from their archives" {
		t.Errorf("owner text = %q", owner)
	}
	if (DiscardCard{Zone: Archives, Selection: Random{}}).validate() == nil {
		t.Error("unset player should be invalid")
	}
	if (DiscardCard{Player: Opponent, Zone: Archives, Selection: Random{}}).validate() != nil {
		t.Error("set player should be valid")
	}

	e.Resolve(ctx)
	if g.State.Archives[1].Count != 1 {
		t.Errorf(
			"archives count = %d, want 1 after one discard",
			g.State.Archives[1].Count,
		)
	}
	if g.State.Discard[1].Count != 1 {
		t.Errorf("discard count = %d, want 1", g.State.Discard[1].Count)
	}
	// The discarded card is one of the two; the other remains.
	if g.State.Archives[1].contains(a) == g.State.Archives[1].contains(b) {
		t.Error("exactly one of the two cards should remain in archives")
	}

	// Empty archives is a no-op.
	e.Resolve(ctx)
	e.Resolve(ctx) // archives now empty
	e.Resolve(ctx)
	if g.State.Discard[1].Count != 2 {
		t.Errorf(
			"discard count = %d, want 2 (empty-archives discards are no-ops)",
			g.State.Discard[1].Count,
		)
	}

	// A card that is not in the named archives is left where it is.
	loose := g.AddToHand(NewCard("loose", Mars, Tactic, Common), 0)
	before := g.State.Discard[0].Count
	g.DiscardCardFromArchives(0, loose)
	if g.State.Discard[0].Count != before || !g.State.Hand[0].contains(loose) {
		t.Error("discarding a card absent from archives should be a no-op")
	}
}

func TestDiscardFromHandEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToHand(NewCard("a", Logos, Tactic, Common), 0)
	g.AddToHand(NewCard("b", Logos, Tactic, Common), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	one := DiscardCard{
		Player:    Controller,
		Zone:      Hand,
		Selection: Chosen{},
		Amount:    1,
	}
	if one.Text() != "discard a card from your hand" {
		t.Errorf("text = %q", one.Text())
	}
	two := DiscardCard{
		Player:    Controller,
		Zone:      Hand,
		Selection: Chosen{},
		Amount:    2,
	}
	if two.Text() != "discard 2 cards from your hand" {
		t.Errorf("plural text = %q", two.Text())
	}

	// The default chooser discards the first hand card (a).
	one.Resolve(ctx)
	if g.State.Hand[0].contains(a) || !g.State.Discard[0].contains(a) {
		t.Error("chosen card should be discarded")
	}

	// Discarding more than the hand holds stops when the hand empties.
	(DiscardCard{Player: Controller, Zone: Hand, Selection: Chosen{}, Amount: 5}).Resolve(
		ctx,
	)
	if g.State.Hand[0].Count != 0 {
		t.Errorf("hand should be empty, got %d", g.State.Hand[0].Count)
	}
}

func TestDiscardFromHandEffectDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToHand(NewCard("c", Logos, Tactic, Common), 0)
	g.AddToHand(NewCard("d", Logos, Tactic, Common), 0)
	g.SetChooser(0, orderRejectChooser{})
	ctx := &EffectContext{Resolver: g, Controller: 0}
	(DiscardCard{Player: Controller, Zone: Hand, Selection: Chosen{}, Amount: 1}).Resolve(
		ctx,
	)
	if g.State.Discard[0].Count != 0 {
		t.Error("a declined discard choice should discard nothing")
	}
}

func TestDiscardFromHandAnyNumber(t *testing.T) {
	t.Run("text renders any number", func(t *testing.T) {
		e := DiscardCard{
			Player:    Controller,
			Zone:      Hand,
			Selection: Chosen{Optional: true},
			AnyNumber: true,
		}
		if got := e.Text(); got != "discard any number of cards from your hand" {
			t.Errorf("text = %q", got)
		}
	})

	t.Run(
		"discards every card and records each on the context",
		func(t *testing.T) {
			g := NewGame("A", "B", 1)
			g.AddToHand(NewCard("a", Logos, Tactic, Common), 0)
			g.AddToHand(NewCard("b", Logos, Tactic, Common), 0)
			g.AddToHand(NewCard("c", Logos, Tactic, Common), 0)
			ctx := &EffectContext{Resolver: g, Controller: 0}
			(DiscardCard{Player: Controller, Zone: Hand, Selection: Chosen{Optional: true}, AnyNumber: true}).Resolve(
				ctx,
			)
			if g.State.Hand[0].Count != 0 {
				t.Errorf("hand = %d, want 0", g.State.Hand[0].Count)
			}
			if got := len(ctx.Produced.Discarded); got != 3 {
				t.Errorf("recorded discards = %d, want 3", got)
			}
		},
	)

	t.Run("declining discards nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.AddToHand(NewCard("d", Logos, Tactic, Common), 0)
		g.SetChooser(0, &declineAfterChooser{})
		ctx := &EffectContext{Resolver: g, Controller: 0}
		(DiscardCard{Player: Controller, Zone: Hand, Selection: Chosen{Optional: true}, AnyNumber: true}).Resolve(
			ctx,
		)
		if g.State.Discard[0].Count != 0 {
			t.Error("a declined discard should discard nothing")
		}
		if len(ctx.Produced.Discarded) != 0 {
			t.Error("a declined discard should record nothing")
		}
	})
}

func TestDiscardFromHandCreaturesOnlyGate(t *testing.T) {
	g := NewGame("A", "B", 1)
	creature := g.AddToHand(
		NewCard("beast", Mars, Creature, Common, WithPower(2)),
		0,
	)
	g.AddToHand(NewCard("tactic", Mars, Tactic, Common), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DiscardCard{
		Player:    Controller,
		Zone:      Hand,
		Selection: Chosen{Type: Creature},
		Amount:    1,
	}
	if e.Text() != "discard a creature from your hand" {
		t.Errorf("text = %q", e.Text())
	}
	plural := DiscardCard{
		Player:    Controller,
		Zone:      Hand,
		Selection: Chosen{Type: Creature},
		Amount:    2,
	}
	if plural.Text() != "discard 2 creatures from your hand" {
		t.Errorf("plural text = %q", plural.Text())
	}

	// Only the creature is a candidate, so it is discarded and the gate reports true.
	if !e.resolveGate(ctx) {
		t.Error("gate should report a discard happened")
	}
	if g.State.Hand[0].contains(creature) {
		t.Error("the creature should have been discarded")
	}

	// With no creatures left in hand, the gate reports false.
	if e.resolveGate(ctx) {
		t.Error("gate should report false when no creature can be discarded")
	}
}

func TestDiscardFromHandValidate(t *testing.T) {
	if (DiscardCard{Zone: Hand, Selection: Random{}}).validate() == nil {
		t.Error("an unset player should be invalid")
	}
	if (DiscardCard{Player: Controller, Zone: Hand}).validate() == nil {
		t.Error("a nil selection should be invalid")
	}
	if (DiscardCard{Player: Controller, Selection: Random{}}).validate() == nil {
		t.Error("an unset zone should be invalid")
	}
	if (DiscardCard{Player: Controller, Zone: Hand, Selection: Random{}, AnyNumber: true, Amount: 2}).validate() == nil {
		t.Error("AnyNumber paired with Amount should be invalid")
	}
	if (DiscardCard{Player: Controller, Zone: Hand, Selection: Random{}, Amount: 2}).validate() != nil {
		t.Error("a fixed amount should be valid")
	}
	if (DiscardCard{Player: Controller, Zone: Archives, Selection: Random{}, OrArchives: true}).validate() == nil {
		t.Error("OrArchives should require Zone Hand")
	}
	if (DiscardCard{Player: Controller, Zone: Hand, Selection: Random{}, OrArchives: true}).validate() != nil {
		t.Error("OrArchives with Zone Hand should be valid")
	}
}

// TestDiscardFromHandOrArchives covers the combined hand-or-archives source
// (Munchling, Novu Dynamo): the discard draws from both piles, names them
// together, and removes the picked card from whichever pile it sat in.
func TestDiscardFromHandOrArchives(t *testing.T) {
	e := DiscardCard{Player: Controller, Zone: Hand, OrArchives: true, Selection: Chosen{}}
	if got := e.Text(); got != "discard a card from your hand or archives" {
		t.Errorf("text = %q", got)
	}

	t.Run("picks a card from the hand", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		inHand := g.AddToHand(NewCard("h", Logos, Tactic, Common), 0)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		e.Resolve(ctx)
		if !g.State.Discard[0].contains(inHand) {
			t.Error("the hand card should be discarded")
		}
	})

	t.Run("picks a card from the archives", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		inArchives := g.AddToArchives(NewCard("a", Logos, Tactic, Common), 0)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		e.Resolve(ctx)
		if g.State.Archives[0].Count != 0 {
			t.Errorf("archives count = %d, want 0", g.State.Archives[0].Count)
		}
		if !g.State.Discard[0].contains(inArchives) {
			t.Error("the archived card should be discarded from the archives")
		}
	})
}
