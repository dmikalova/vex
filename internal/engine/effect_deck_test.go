package engine

import (
	"slices"
	"strings"
	"testing"
)

func TestChaosPortalComposition(t *testing.T) {
	effect := ChooseHouseThen{Then: Sentences{Effects: []Effect{
		RevealTopOfDeck{},
		Conditional{Cond: ItIsOfHouse{House: TheChosenHouse}, Then: PlayRevealedCard{}},
	}}}
	if got := effect.Text(); got != "choose a house - reveal the top card of your deck. If it is of the chosen house, play it." {
		t.Errorf("text = %q", got)
	}

	g := started(t) // Brobnar is active; the chosen house can still play a Logos card.
	top := g.AddToDeck(NewCard("Portal Scout", Logos, Creature, Common, WithPower(2)), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0, ChosenHouse: Logos}

	Sequence{
		Effects: []Effect{
			RevealTopOfDeck{},
			Conditional{Cond: ItIsOfHouse{House: TheChosenHouse}, Then: PlayRevealedCard{}},
		},
	}.Resolve(
		ctx,
	)

	if g.State.Deck[0].Count != 0 {
		t.Errorf("deck count = %d, want 0", g.State.Deck[0].Count)
	}
	if got := g.Battleline(0); len(got) != 1 || got[0] != top {
		t.Errorf("battleline = %v, want [%d]", got, top)
	}
	if !g.State.Cards[top].Exhausted {
		t.Error("played creature should enter exhausted")
	}
	if played := g.PlayedThisTurn(0); len(played) != 1 || g.House(played[0]) != Logos {
		t.Errorf("played this turn = %v, want one Logos card", played)
	}
	if log := g.LogText(); len(log) == 0 || !strings.Contains(log[len(log)-1], "Portal Scout") {
		t.Errorf("log = %v, want the top card revealed and played", log)
	}
}

func TestChaosPortalMissesAndGuards(t *testing.T) {
	g := started(t)
	top := g.AddToDeck(NewCard("Wrong House", Dis, Tactic, Common), 0)
	reveal := Sequence{
		Effects: []Effect{
			RevealTopOfDeck{},
			Conditional{Cond: ItIsOfHouse{House: TheChosenHouse}, Then: PlayRevealedCard{}},
		},
	}
	reveal.Resolve(&EffectContext{Resolver: g, Controller: 0, ChosenHouse: Logos})
	if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != top {
		t.Errorf("non-matching top card moved: deck = %v, want [%d]", g.State.Deck[0].slice(), top)
	}
	if len(g.PlayedThisTurn(0)) != 0 {
		t.Errorf(
			"non-matching card should not count as played, got %d",
			len(g.PlayedThisTurn(0)),
		)
	}

	// An empty deck reveals nothing and plays nothing.
	g2 := started(t)
	before := len(g2.Log)
	reveal.Resolve(&EffectContext{Resolver: g2, Controller: 0, ChosenHouse: Logos})
	if len(g2.Log) != before {
		t.Error("an empty deck should not reveal anything")
	}
}

func TestPlayTopOfDeckLeavesUnplayableCardOnTop(t *testing.T) {
	t.Run("game state gates", func(t *testing.T) {
		cases := []struct {
			name  string
			setup func(*Game)
			err   error
		}{
			{
				name: "game over",
				setup: func(g *Game) {
					g.State.Winner = 0
				},
				err: ErrGameOver,
			},
			{
				name: "not active player",
				setup: func(g *Game) {
					g.State.ActivePlayer = 1
				},
				err: ErrNotActivePlayer,
			},
			{
				name: "card play limit",
				setup: func(g *Game) {
					g.AddToBattleline(NewCard(
						"Imp",
						Brobnar,
						Creature,
						Common,
						WithPower(1),
						WithRestrictions(
							Restrictions{
								PlayCardLimit: PlayCardLimit{Player: Controller, Amount: 1},
							},
						),
					), 0)
					g.State.PlayedThisTurn[0].Count = 1
				},
				err: ErrCardPlayLimit,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				g := started(t)
				top := g.AddToDeck(
					NewCard("Portal Scout", Logos, Creature, Common, WithPower(2)),
					0,
				)
				tc.setup(g)
				if _, err := g.playCardFromZone(
					0,
					top,
					func() { g.State.Deck[0].removeAt(0) },
					playCardOptions{},
				); err != tc.err {
					t.Fatalf("playCardFromZone err = %v, want %v", err, tc.err)
				}
				if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != top {
					t.Errorf(
						"guarded play should leave top card in deck, got %v",
						g.State.Deck[0].slice(),
					)
				}
			})
		}

	})

	t.Run("creature play restriction", func(t *testing.T) {
		g := started(t)
		g.AddToBattleline(NewCard("Blocker", Brobnar, Creature, Common, WithPower(1),
			WithRestrictions(Restrictions{CannotPlay: Creature})), 0)
		top := g.AddToDeck(NewCard("Portal Scout", Logos, Creature, Common, WithPower(2)), 0)
		if _, err := g.playCardFromZone(
			0,
			top,
			func() { g.State.Deck[0].removeAt(0) },
			playCardOptions{},
		); err != ErrCannotPlayCreature {
			t.Fatalf("playCardFromZone err = %v, want %v", err, ErrCannotPlayCreature)
		}
		if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != top {
			t.Errorf("barred creature should stay in deck, got %v", g.State.Deck[0].slice())
		}
	})

	t.Run("unmet play requirement", func(t *testing.T) {
		g := started(t)
		g.State.Aember[0] = 6
		top := g.AddToDeck(NewCard("Kelifi Dragon", Brobnar, Creature, Rare,
			WithPower(12), WithPlayRequirement(AemberThreshold(7))), 0)
		PlayTopOfDeck{}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != top {
			t.Errorf("a card whose Æmber threshold is unmet should stay on top of the deck, got %v",
				g.State.Deck[0].slice())
		}
	})

	t.Run("unknown card type", func(t *testing.T) {
		g := started(t)
		top := g.AddToDeck(NewCard("Mystery", Logos, AnyType, Common), 0)
		if _, err := g.playCardFromZone(
			0,
			top,
			func() { g.State.Deck[0].removeAt(0) },
			playCardOptions{},
		); err != ErrWrongType {
			t.Fatalf("playCardFromZone err = %v, want %v", err, ErrWrongType)
		}
		if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != top {
			t.Errorf("unknown type should stay in deck, got %v", g.State.Deck[0].slice())
		}
	})
}

func TestDiscardTopOfEachDeckAmount(t *testing.T) {
	g := started(t)
	source := g.AddArtifact(NewCard("Rigged Lottery Source", Shadows, Tactic, Rare), 0)
	for range 3 {
		g.AddToDeck(NewCard("Shadows Card", Shadows, Creature, Common), 0)
		g.AddToDeck(NewCard("Mars Card", Mars, Creature, Common), 1)
	}
	ctx := &EffectContext{Resolver: g, Source: source, Controller: 0}

	e := DiscardTopOfEachDeck{Amount: 2}
	if got := e.Text(); got != "discard the top 2 cards of each player's deck" {
		t.Errorf("text = %q", got)
	}
	e.Resolve(ctx)

	if got := len(g.Discard(0)); got != 2 {
		t.Errorf("controller discard = %d, want 2", got)
	}
	if got := len(g.Discard(1)); got != 2 {
		t.Errorf("opponent discard = %d, want 2", got)
	}
	if got := len(ctx.Produced.Discarded); got != 4 {
		t.Errorf("recorded discards = %d, want 4", got)
	}
}

func TestBonkersComposition(t *testing.T) {
	g := started(t)
	source := g.AddArtifact(NewCard("Bonkers Killing Machine", Logos, Artifact, Rare), 0)
	p1Top := g.AddToDeck(NewCard("Mars Top", Mars, Tactic, Common), 0)
	p2Top := g.AddToDeck(NewCard("Dis Top", Dis, Tactic, Common), 1)
	marsCreature := g.AddToBattleline(
		NewCard("Mars Creature", Mars, Creature, Common, WithPower(4)),
		0,
	)
	disArtifact := g.AddArtifact(NewCard("Dis Artifact", Dis, Artifact, Common), 1)
	bystander := g.AddToBattleline(testCreature("bystander", 4), 1)
	ctx := &EffectContext{Resolver: g, Source: source, Controller: 0}

	effect := Sentences{Effects: []Effect{
		DiscardTopOfEachDeck{},
		ForEachDiscarded{
			Do: Destroy{
				Target: Target{Kind: TargetChosenCreatureOrArtifact}.OfContextualHouse(),
			},
		},
		Conditional{
			Cond: CardsDestroyedFewerThan{Amount: 2},
			Then: Destroy{Target: Target{Kind: TargetThisCreature}},
		},
	}}
	if got := effect.Text(); got != "discard the top card of each player's deck. For each card discarded this way, destroy a creature or artifact of that card's house. If fewer than 2 cards are destroyed this way, destroy {self}." {
		t.Errorf("text = %q", got)
	}

	effect.Resolve(ctx)

	if discard := g.Discard(
		0,
	); len(discard) != 2 || discard[0] != p1Top ||
		discard[1] != marsCreature {
		t.Errorf("controller discard = %v, want top card then Mars creature", discard)
	}
	if discard := g.Discard(
		1,
	); len(discard) != 2 || discard[0] != p2Top ||
		discard[1] != disArtifact {
		t.Errorf("opponent discard = %v, want top card then Dis artifact", discard)
	}
	if !g.inPlay(source) {
		t.Error("source should stay in play when two cards are destroyed")
	}
	if !g.inPlay(bystander) {
		t.Error("off-house bystander should stay in play")
	}
}

func TestBonkersCompositionSelfDestructs(t *testing.T) {
	g := started(t)
	source := g.AddArtifact(NewCard("Bonkers Killing Machine", Logos, Artifact, Rare), 0)
	top := g.AddToDeck(NewCard("Mars Top", Mars, Tactic, Common), 0)
	bystander := g.AddToBattleline(testCreature("bystander", 4), 1)
	ctx := &EffectContext{Resolver: g, Source: source, Controller: 0}

	Sentences{Effects: []Effect{
		DiscardTopOfEachDeck{},
		ForEachDiscarded{
			Do: Destroy{
				Target: Target{Kind: TargetChosenCreatureOrArtifact}.OfContextualHouse(),
			},
		},
		Conditional{
			Cond: CardsDestroyedFewerThan{Amount: 2},
			Then: Destroy{Target: Target{Kind: TargetThisCreature}},
		},
	}}.Resolve(ctx)

	if discard := g.Discard(0); len(discard) != 2 || discard[0] != top || discard[1] != source {
		t.Errorf("controller discard = %v, want discarded top card then source", discard)
	}
	if !g.inPlay(bystander) {
		t.Error("off-house bystander should stay in play")
	}
}

func TestForEachDiscardedAndContextualHouse(t *testing.T) {
	// validate surfaces a bad Do.
	if err := validateEffect(ForEachDiscarded{Do: Destroy{}}); err == nil {
		t.Error("ForEachDiscarded should reject a Do with no target")
	}
	if err := validateEffect(
		ForEachDiscarded{
			Do: Destroy{Target: Target{Kind: TargetChosenCreatureOrArtifact}.OfContextualHouse()},
		},
	); err != nil {
		t.Errorf("valid ForEachDiscarded = %v", err)
	}

	// Text renders the chosen-in-play noun and the contextual-house clause.
	target := Target{Kind: TargetChosenCreatureOrArtifact}.OfContextualHouse()
	if got := target.Text(); got != "a creature or artifact of that card's house" {
		t.Errorf("target text = %q", got)
	}

	// With no card in context, the contextual-house filter selects nothing.
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("c", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if got := target.Select(ctx); got != nil {
		t.Errorf("no-context select = %v, want nil", got)
	}
}

func TestResolverInPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	live := g.AddToBattleline(testCreature("live", 3), 0)
	art := g.AddArtifact(NewCard("Relic", Logos, Artifact, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if !resolverInPlay(ctx, live) {
		t.Error("a battleline creature should read in play")
	}
	if !resolverInPlay(ctx, art) {
		t.Error("an artifact should read in play")
	}
	if resolverInPlay(ctx, LocalID(200)) {
		t.Error("an absent id should not read in play")
	}
}

func TestEvasionSigilComposition(t *testing.T) {
	g := started(t) // Brobnar is active
	src := g.AddToBattleline(testCreature("attacker", 5), 0)
	top := g.AddToDeck(NewCard("Brobnar Top", Brobnar, Tactic, Common), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := Sentences{Effects: []Effect{
		DiscardTopOfDeck{},
		Conditional{Cond: ItIsOfHouse{House: TheActiveHouse}, Then: CancelFight{}},
	}}
	if got := e.Text(); got != "discard the top card of its controller's deck. If it is of the active house, the fight does not occur." {
		t.Errorf("text = %q", got)
	}

	e.Resolve(ctx)

	if discard := g.Discard(0); len(discard) != 1 || discard[0] != top {
		t.Errorf("discard = %v, want top card %d", discard, top)
	}
	if !g.State.FightCancelled {
		t.Error("active-house discard should cancel the current fight")
	}
}

func TestEvasionSigilCompositionMiss(t *testing.T) {
	g := started(t) // Brobnar is active
	src := g.AddToBattleline(testCreature("attacker", 5), 0)
	top := g.AddToDeck(NewCard("Mars Top", Mars, Tactic, Common), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	Sequence{Effects: []Effect{
		DiscardTopOfDeck{},
		Conditional{Cond: ItIsOfHouse{House: TheActiveHouse}, Then: CancelFight{}},
	}}.Resolve(ctx)

	if discard := g.Discard(0); len(discard) != 1 || discard[0] != top {
		t.Errorf("discard = %v, want top card %d", discard, top)
	}
	if g.State.FightCancelled {
		t.Error("off-house discard should not cancel the current fight")
	}
	// An empty deck puts no card in context, so the fight is not cancelled.
	g2 := started(t)
	src2 := g2.AddToBattleline(testCreature("attacker", 5), 0)
	Sequence{Effects: []Effect{
		DiscardTopOfDeck{},
		Conditional{Cond: ItIsOfHouse{House: TheActiveHouse}, Then: CancelFight{}},
	}}.Resolve(&EffectContext{Resolver: g2, Source: src2, Controller: 0})
	if g2.State.FightCancelled {
		t.Error("an empty deck should not cancel the fight")
	}
}

// TestDiscardDeckUntil covers the dig through the top of the deck: it stops at a
// card the filters admit, reports success so a Then can follow, and runs the deck
// out when nothing matches.
func TestDiscardDeckUntil(t *testing.T) {
	e := DiscardDeckUntil{Type: Creature, House: Brobnar}
	want := "discard cards from the top of your deck until you discard a Brobnar creature or run out of cards"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
	if got := (DiscardDeckUntil{Type: Artifact}).Text(); got !=
		"discard cards from the top of your deck until you discard an artifact or run out of cards" {
		t.Errorf("artifact text = %q", got)
	}
	if got := (DiscardDeckUntil{}).Text(); got !=
		"discard cards from the top of your deck until you discard a card or run out of cards" {
		t.Errorf("plain text = %q", got)
	}

	if got := (PutDiscardedIntoHand{}).Text(); got != "put the discarded card into your hand" {
		t.Errorf("tail text = %q", got)
	}
	if got := (PutDiscardedIntoHand{Type: Artifact}).Text(); got !=
		"put the discarded artifact into your hand" {
		t.Errorf("artifact text = %q", got)
	}
	if got := (PutDiscardedIntoHand{Type: Creature}).Text(); got !=
		"put the discarded creature into your hand" {
		t.Errorf("creature tail text = %q", got)
	}

	g := NewGame("A", "B", 1)
	skipped := g.AddToDeck(NewCard("Trick", Brobnar, Tactic, Common), 0)
	brute := g.AddToDeck(NewCard("Brute", Brobnar, Creature, Common, WithPower(5)), 0)
	buried := g.AddToDeck(NewCard("Buried", Brobnar, Creature, Common, WithPower(5)), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	Then{First: e, Result: PutDiscardedIntoHand{}}.Resolve(ctx)
	if len(g.Hand(0)) != 1 || g.Hand(0)[0] != brute {
		t.Errorf("hand = %v, want [%d]", g.Hand(0), brute)
	}
	if len(g.Discard(0)) != 1 || g.Discard(0)[0] != skipped {
		t.Errorf("discard = %v, want [%d]", g.Discard(0), skipped)
	}
	if len(g.Deck(0)) != 1 || g.Deck(0)[0] != buried {
		t.Errorf("deck = %v, want [%d]", g.Deck(0), buried)
	}

	// Nothing matching left: the dig empties the deck and the tail does nothing.
	Then{First: DiscardDeckUntil{Type: Artifact}, Result: PutDiscardedIntoHand{}}.Resolve(ctx)
	if len(g.Deck(0)) != 0 {
		t.Errorf("deck should be empty, got %v", g.Deck(0))
	}
	if len(g.Hand(0)) != 1 {
		t.Errorf("hand should be untouched, got %v", g.Hand(0))
	}
	// Resolved bare, the dig still runs; it just has no tail to gate.
	DiscardDeckUntil{Type: Artifact}.Resolve(ctx)
	PutDiscardedIntoHand{}.Resolve(ctx)
}

// TestRevealDeckUntilHouse covers the reveal-and-archive dig: it archives each
// card revealed until one of the named house turns up (reporting success so a Then
// can follow) or the controller stops, and runs safely on an empty deck.
func TestRevealDeckUntilHouse(t *testing.T) {
	want := "reveal cards from the top of your deck until you reveal a Brobnar card" +
		" or choose to stop, archiving each card revealed this way"
	if got := (RevealDeckUntilHouse{House: Brobnar}).Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
	if err := (RevealDeckUntilHouse{}).validate(); err == nil {
		t.Error("want validate error for unset house")
	}
	if err := (RevealDeckUntilHouse{House: Brobnar}).validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}

	// Keep revealing until a Brobnar card turns up: both cards above it and the
	// Brobnar card itself are archived, the card below it stays on the deck.
	g := NewGame("A", "B", 1)
	g.SetChooser(0, optionPicker{idx: 0})
	g.AddToDeck(NewCard("Trick", Logos, Tactic, Common), 0)
	g.AddToDeck(NewCard("Brute", Brobnar, Creature, Common, WithPower(5)), 0)
	deep := g.AddToDeck(NewCard("Deep", Logos, Creature, Common), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if !(RevealDeckUntilHouse{House: Brobnar}).resolveGate(ctx) {
		t.Error("resolveGate = false, want true when a Brobnar card is revealed")
	}
	if len(g.Archives(0)) != 2 {
		t.Errorf("archives = %v, want 2 cards", g.Archives(0))
	}
	if len(g.Deck(0)) != 1 || g.Deck(0)[0] != deep {
		t.Errorf("deck = %v, want [%d]", g.Deck(0), deep)
	}

	// Choosing to stop archives only the card revealed so far and reports failure.
	g2 := NewGame("A", "B", 1)
	g2.SetChooser(0, optionPicker{idx: 1})
	g2.AddToDeck(NewCard("One", Logos, Creature, Common), 0)
	kept := g2.AddToDeck(NewCard("Two", Logos, Creature, Common), 0)
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	if (RevealDeckUntilHouse{House: Brobnar}).resolveGate(ctx2) {
		t.Error("resolveGate = true, want false when the controller stops")
	}
	if len(g2.Archives(0)) != 1 {
		t.Errorf("archives = %v, want 1 card", g2.Archives(0))
	}
	if len(g2.Deck(0)) != 1 || g2.Deck(0)[0] != kept {
		t.Errorf("deck = %v, want [%d]", g2.Deck(0), kept)
	}

	// An empty deck reports failure and moves nothing; the bare Resolve is a no-op.
	g3 := NewGame("A", "B", 1)
	ctx3 := &EffectContext{Resolver: g3, Controller: 0}
	if (RevealDeckUntilHouse{House: Brobnar}).resolveGate(ctx3) {
		t.Error("resolveGate = true, want false on an empty deck")
	}
	RevealDeckUntilHouse{House: Brobnar}.Resolve(ctx3)
}

// TestRevealPurgeShuffleDeck covers the reveal-top-N/purge-one/shuffle dig: it
// purges one of the revealed top cards and leaves the rest in the deck. With the
// default chooser it reveals from the controller's own deck and purges the top
// card revealed.
func TestRevealPurgeShuffleDeck(t *testing.T) {
	want := "reveal the top 5 cards of a player's deck. Purge a card revealed this " +
		"way. Shuffle the other revealed cards into that deck"
	if got := (RevealPurgeShuffleDeck{Amount: 5}).Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
	if err := (RevealPurgeShuffleDeck{}).validate(); err == nil {
		t.Error("want validate error for unset amount")
	}
	if err := (RevealPurgeShuffleDeck{Amount: 5}).validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}

	g := NewGame("A", "B", 1)
	top := g.AddToDeck(NewCard("Top", Logos, Creature, Common), 0)
	g.AddToDeck(NewCard("Mid", Logos, Creature, Common), 0)
	g.AddToDeck(NewCard("Low", Logos, Creature, Common), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	RevealPurgeShuffleDeck{Amount: 3}.Resolve(ctx)
	if purge := g.Purge(0); len(purge) != 1 || purge[0] != top {
		t.Errorf("purge = %v, want [%d]", purge, top)
	}
	if len(g.Deck(0)) != 2 {
		t.Errorf("deck = %v, want 2 cards", g.Deck(0))
	}

	// An empty deck purges nothing and is a no-op.
	g2 := NewGame("A", "B", 1)
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	RevealPurgeShuffleDeck{Amount: 3}.Resolve(ctx2)
	if len(g2.Purge(0)) != 0 {
		t.Errorf("purge = %v, want empty", g2.Purge(0))
	}

	// Choosing the opponent's deck digs there instead.
	g3 := NewGame("A", "B", 1)
	oppTop := g3.AddToDeck(NewCard("OppTop", Logos, Creature, Common), 1)
	ctx3 := &EffectContext{Resolver: g3, Controller: 0}
	g3.SetChooser(0, optionPicker{idx: 1})
	RevealPurgeShuffleDeck{Amount: 3}.Resolve(ctx3)
	if purge := g3.Purge(1); len(purge) != 1 || purge[0] != oppTop {
		t.Errorf("opponent purge = %v, want [%d]", purge, oppTop)
	}
}

// TestLookAtTop covers the "look at the top N, keep one, discard the rest" dig:
// the chosen card goes to hand, the others to the discard pile, and short or empty
// decks are handled.
func TestLookAtTop(t *testing.T) {
	if got := (LookAtTop{Amount: 3}).Text(); got !=
		"look at the top 3 cards of your deck, put 1 into your hand, and discard the others" {
		t.Errorf("Text() = %q", got)
	}
	if (LookAtTop{}).validate() == nil {
		t.Error("a Count of 0 should be rejected")
	}
	if err := (LookAtTop{Amount: 3}).validate(); err != nil {
		t.Errorf("validate() = %v", err)
	}

	t.Run("keeps the chosen card and discards the rest", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToDeck(NewCard("A Card", Logos, Creature, Common, WithPower(2)), 0)
		b := g.AddToDeck(NewCard("B Card", Logos, Tactic, Common), 0)
		c := g.AddToDeck(NewCard("C Card", Logos, Artifact, Common), 0)
		bottom := g.AddToDeck(NewCard("Bottom", Logos, Creature, Common, WithPower(1)), 0)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{b}})
		LookAtTop{Amount: 3}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Hand(0); len(got) != 1 || got[0] != b {
			t.Errorf("hand = %v, want [%d]", got, b)
		}
		if got := g.Discard(0); !slices.Contains(got, a) || !slices.Contains(got, c) {
			t.Errorf("discard = %v, want to contain %d and %d", got, a, c)
		}
		if got := g.Deck(0); len(got) != 1 || got[0] != bottom {
			t.Errorf("deck = %v, want [%d]", got, bottom)
		}
	})

	t.Run("looks at as many as remain", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToDeck(NewCard("A Card", Logos, Creature, Common, WithPower(2)), 0)
		b := g.AddToDeck(NewCard("B Card", Logos, Tactic, Common), 0)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{a}})
		LookAtTop{Amount: 3}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Hand(0); len(got) != 1 || got[0] != a {
			t.Errorf("hand = %v, want [%d]", got, a)
		}
		if got := g.Discard(0); len(got) != 1 || got[0] != b {
			t.Errorf("discard = %v, want [%d]", got, b)
		}
	})

	t.Run("an empty deck does nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		LookAtTop{Amount: 3}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Hand(0); len(got) != 0 {
			t.Errorf("hand = %v, want empty", got)
		}
	})

	t.Run("a declined choice keeps everything in the deck", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.AddToDeck(NewCard("A Card", Logos, Creature, Common, WithPower(2)), 0)
		g.AddToDeck(NewCard("B Card", Logos, Tactic, Common), 0)
		g.SetChooser(0, orderRejectChooser{})
		LookAtTop{Amount: 3}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Hand(0); len(got) != 0 {
			t.Errorf("hand = %v, want empty", got)
		}
		if got := g.Deck(0); len(got) != 2 {
			t.Errorf("deck = %v, want 2 cards", got)
		}
	})
}

// TestReorderTop covers "look at the top N and put them back in any order": the
// controller's chosen order becomes the new top, short decks are left alone, and a
// declined choice leaves the order untouched.
func TestReorderTop(t *testing.T) {
	if got := (ReorderTop{Amount: 3}).Text(); got !=
		"look at the top 3 cards of your deck and put them back in any order" {
		t.Errorf("Text() = %q", got)
	}
	if (ReorderTop{}).validate() == nil {
		t.Error("a Count of 0 should be rejected")
	}
	if err := (ReorderTop{Amount: 3}).validate(); err != nil {
		t.Errorf("validate() = %v", err)
	}

	t.Run("places the chosen order on top, leaving the rest", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToDeck(NewCard("A Card", Logos, Creature, Common, WithPower(2)), 0)
		b := g.AddToDeck(NewCard("B Card", Logos, Tactic, Common), 0)
		c := g.AddToDeck(NewCard("C Card", Logos, Artifact, Common), 0)
		bottom := g.AddToDeck(NewCard("Bottom", Logos, Creature, Common, WithPower(1)), 0)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{c, a}})
		ReorderTop{Amount: 3}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Deck(0); len(got) != 4 ||
			got[0] != c || got[1] != a || got[2] != b || got[3] != bottom {
			t.Errorf("deck = %v, want [%d %d %d %d]", got, c, a, b, bottom)
		}
	})

	t.Run("reorders as many as remain", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToDeck(NewCard("A Card", Logos, Creature, Common, WithPower(2)), 0)
		b := g.AddToDeck(NewCard("B Card", Logos, Tactic, Common), 0)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{b}})
		ReorderTop{Amount: 3}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Deck(0); len(got) != 2 || got[0] != b || got[1] != a {
			t.Errorf("deck = %v, want [%d %d]", got, b, a)
		}
	})

	t.Run("fewer than two cards does nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		only := g.AddToDeck(NewCard("Only", Logos, Creature, Common, WithPower(2)), 0)
		ReorderTop{Amount: 3}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Deck(0); len(got) != 1 || got[0] != only {
			t.Errorf("deck = %v, want [%d]", got, only)
		}
	})

	t.Run("a declined choice keeps the original order", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToDeck(NewCard("A Card", Logos, Creature, Common, WithPower(2)), 0)
		b := g.AddToDeck(NewCard("B Card", Logos, Tactic, Common), 0)
		g.SetChooser(0, orderRejectChooser{})
		ReorderTop{Amount: 3}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Deck(0); len(got) != 2 || got[0] != a || got[1] != b {
			t.Errorf("deck = %v, want [%d %d]", got, a, b)
		}
	})
}

func TestDiscardTopAndForEachDiscardedHouseFilter(t *testing.T) {
	// DiscardTop validate and text.
	if err := (DiscardTop{Player: Controller, Amount: 0}).validate(); err == nil {
		t.Error("DiscardTop with Amount 0 should fail validation")
	}
	if err := (DiscardTop{Player: Controller, Amount: 2}).validate(); err != nil {
		t.Errorf("valid DiscardTop = %v", err)
	}
	if got := (DiscardTop{Player: Controller, Amount: 1}).Text(); got != "discard the top 1 card of your deck" {
		t.Errorf("singular text = %q", got)
	}
	if got := (DiscardTop{Player: Opponent, Amount: 2}).Text(); got != "discard the top 2 cards of your opponent's deck" {
		t.Errorf("opponent text = %q", got)
	}

	if got := (ForEachDiscarded{Do: GainAember{Player: Controller, Amount: 1}}).Text(); got != "for each card discarded this way, gain 1 Æmber" {
		t.Errorf("unfiltered text = %q", got)
	}
	if got := (ForEachDiscarded{House: Logos, Do: GainAember{Player: Controller, Amount: 1}}).Text(); got != "for each Logos card discarded this way, gain 1 Æmber" {
		t.Errorf("filtered text = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.AddToDeck(NewCard("L1", Logos, Tactic, Common), 0)
	g.AddToDeck(NewCard("M", Mars, Tactic, Common), 0)
	g.AddToDeck(NewCard("L2", Logos, Tactic, Common), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	Sentences{Effects: []Effect{
		DiscardTop{Player: Controller, Amount: 3},
		ForEachDiscarded{House: Logos, Do: GainAember{Player: Controller, Amount: 1}},
	}}.Resolve(ctx)

	if g.Aember(0) != 2 {
		t.Errorf("gained %d Æmber, want 2 (one per Logos card)", g.Aember(0))
	}
}
