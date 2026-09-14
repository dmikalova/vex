package engine

import (
	"slices"
	"strings"
	"testing"
)

func TestChaosPortalComposition(t *testing.T) {
	effect := ChooseHouseThen{Then: Sentences{Effects: []Effect{
		RevealTopOfDeck{Amount: 1},
		Conditional{Cond: ItIs{House: chosenHouse}, Then: PlayRevealedCard{}},
	}}}
	if got := effect.Text(); got != "choose a house - reveal the top card of your deck. If it is of the chosen house, play it." {
		t.Errorf("text = %q", got)
	}

	g := started(t) // Brobnar is active; the chosen house can still play a Logos card.
	top := g.AddToDeck(NewCard("Portal Scout", Logos, Creature, Common, WithPower(2)), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0, ChosenHouse: Logos}

	Sequence{
		Effects: []Effect{
			RevealTopOfDeck{Amount: 1},
			Conditional{Cond: ItIs{House: chosenHouse}, Then: PlayRevealedCard{}},
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
			RevealTopOfDeck{Amount: 1},
			Conditional{Cond: ItIs{House: chosenHouse}, Then: PlayRevealedCard{}},
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

func TestPutRevealedCard(t *testing.T) {
	texts := map[DeckDest]string{
		IntoArchives: "archive it",
		IntoDiscard:  "discard it",
		IntoHand:     "put it into your hand",
		IntoPurge:    "purge it",
	}
	for dest, want := range texts {
		if got := (PutRevealedCard{To: dest}).Text(); got != want {
			t.Errorf("Text(%d) = %q, want %q", dest, got, want)
		}
	}
	if err := (PutRevealedCard{To: DeckDest(99)}).validate(); err == nil {
		t.Error("an unknown destination should not validate")
	}
	if err := (PutRevealedCard{To: IntoArchives}).validate(); err != nil {
		t.Errorf("archives should validate, got %v", err)
	}

	// Each destination routes the revealed top card out of the deck. Cards are
	// revealed in the order they were added, so they line up with the routes below.
	g := started(t)
	archiveTop := g.AddToDeck(NewCard("Archived", Logos, Creature, Common, WithPower(1)), 0)
	discardTop := g.AddToDeck(NewCard("Discarded", Logos, Creature, Common, WithPower(1)), 0)
	handTop := g.AddToDeck(NewCard("Handed", Logos, Creature, Common, WithPower(1)), 0)
	g.AddToDeck(NewCard("Purged", Logos, Creature, Common, WithPower(1)), 0)
	route := func(dest DeckDest) {
		Sequence{Effects: []Effect{RevealTopOfDeck{Amount: 1}, PutRevealedCard{To: dest}}}.
			Resolve(&EffectContext{Resolver: g, Controller: 0})
	}
	route(IntoArchives)
	if a := g.Archives(0); len(a) != 1 || a[0] != archiveTop {
		t.Errorf("archives = %v, want [%d]", a, archiveTop)
	}
	route(IntoDiscard)
	if d := g.Discard(0); len(d) != 1 || d[0] != discardTop {
		t.Errorf("discard = %v, want [%d]", d, discardTop)
	}
	route(IntoHand)
	if slices.Index(g.Hand(0), handTop) < 0 {
		t.Errorf("hand = %v, want to contain %d", g.Hand(0), handTop)
	}
	route(IntoPurge)
	if g.State.Deck[0].Count != 0 {
		t.Errorf("deck count = %d, want 0 after purging the last card", g.State.Deck[0].Count)
	}

	// With no card in context it is a no-op.
	before := len(g.Archives(0))
	PutRevealedCard{To: IntoArchives}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if after := len(g.Archives(0)); after != before {
		t.Errorf("no-context resolve moved cards: %d -> %d", before, after)
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

	e := DiscardTop{Player: EachPlayer, Amount: 2}
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
		DiscardTop{Player: EachPlayer},
		ForEachDiscarded{
			Do: Destroy{
				Target: Target{Kind: TargetChosenCreatureOrArtifact}.House(contextualHouse),
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
		DiscardTop{Player: EachPlayer},
		ForEachDiscarded{
			Do: Destroy{
				Target: Target{Kind: TargetChosenCreatureOrArtifact}.House(contextualHouse),
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
			Do: Destroy{
				Target: Target{Kind: TargetChosenCreatureOrArtifact}.House(contextualHouse),
			},
		},
	); err != nil {
		t.Errorf("valid ForEachDiscarded = %v", err)
	}

	// Text renders the chosen-in-play noun and the contextual-house clause.
	target := Target{Kind: TargetChosenCreatureOrArtifact}.House(contextualHouse)
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
		DiscardTop{},
		Conditional{Cond: ItIs{House: activeHouse}, Then: CancelFight{}},
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
		DiscardTop{},
		Conditional{Cond: ItIs{House: activeHouse}, Then: CancelFight{}},
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
		DiscardTop{},
		Conditional{Cond: ItIs{House: activeHouse}, Then: CancelFight{}},
	}}.Resolve(&EffectContext{Resolver: g2, Source: src2, Controller: 0})
	if g2.State.FightCancelled {
		t.Error("an empty deck should not cancel the fight")
	}
}

// TestDiscardUntil covers the dig through the top of the deck: it stops at
// a card the filters admit, reports success so a Then can follow, records the run it
// discarded, and runs the deck out when nothing matches.
func TestDiscardUntil(t *testing.T) {
	e := DiscardUntil{Type: Creature, House: Brobnar}
	want := "discard cards from the top of your deck until you discard a Brobnar creature or run out of cards"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
	if got := (DiscardUntil{Type: Artifact}).Text(); got !=
		"discard cards from the top of your deck until you discard an artifact or run out of cards" {
		t.Errorf("artifact text = %q", got)
	}
	if got := (DiscardUntil{}).Text(); got !=
		"discard cards from the top of your deck until you discard a card or run out of cards" {
		t.Errorf("plain text = %q", got)
	}
	if got := (DiscardUntil{House: Brobnar, MayStop: true}).Text(); got !=
		"discard cards from the top of your deck until you discard a Brobnar card or choose to stop" {
		t.Errorf("may-stop text = %q", got)
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
	// The whole discarded run is recorded, matching card last.
	if got := ctx.Produced.Discarded; len(got) != 2 || got[0] != skipped || got[1] != brute {
		t.Errorf("recorded run = %v, want [%d %d]", got, skipped, brute)
	}

	// Nothing matching left: the dig empties the deck and the tail does nothing.
	Then{First: DiscardUntil{Type: Artifact}, Result: PutDiscardedIntoHand{}}.Resolve(ctx)
	if len(g.Deck(0)) != 0 {
		t.Errorf("deck should be empty, got %v", g.Deck(0))
	}
	if len(g.Hand(0)) != 1 {
		t.Errorf("hand should be untouched, got %v", g.Hand(0))
	}
	// Resolved bare, the dig still runs; it just has no tail to gate.
	DiscardUntil{Type: Artifact}.Resolve(ctx)
	PutDiscardedIntoHand{}.Resolve(ctx)
}

// TestDiscardUntilMayStop covers the optional stop: the controller may end
// the dig before a match, which reports failure and leaves the rest of the deck,
// while still recording the run discarded so far.
func TestDiscardUntilMayStop(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetChooser(0, optionPicker{idx: 1})
	first := g.AddToDeck(NewCard("One", Logos, Creature, Common), 0)
	kept := g.AddToDeck(NewCard("Two", Logos, Creature, Common), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if (DiscardUntil{House: Brobnar, MayStop: true}).resolveGate(ctx) {
		t.Error("resolveGate = true, want false when the controller stops")
	}
	if len(g.Discard(0)) != 1 || g.Discard(0)[0] != first {
		t.Errorf("discard = %v, want [%d]", g.Discard(0), first)
	}
	if len(g.Deck(0)) != 1 || g.Deck(0)[0] != kept {
		t.Errorf("deck = %v, want [%d]", g.Deck(0), kept)
	}
	if got := ctx.Produced.Discarded; len(got) != 1 || got[0] != first {
		t.Errorf("recorded run = %v, want [%d]", got, first)
	}

	// An empty deck reports failure and moves nothing; the bare Resolve is a no-op.
	g2 := NewGame("A", "B", 1)
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	if (DiscardUntil{House: Brobnar, MayStop: true}).resolveGate(ctx2) {
		t.Error("resolveGate = true, want false on an empty deck")
	}
	DiscardUntil{House: Brobnar, MayStop: true}.Resolve(ctx2)
}

// TestArchiveDiscardedThisWay covers archiving the run a preceding dig discarded:
// every recorded card moves from the discard pile to archives, and an empty run
// archives nothing.
func TestArchiveDiscardedThisWay(t *testing.T) {
	if got := (ArchiveDiscardedThisWay{}).Text(); got != "archive each card discarded this way" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	a := g.Register(testCreature("a", 1), 0)
	b := g.Register(testCreature("b", 1), 0)
	g.State.Discard[0].add(a)
	g.State.Discard[0].add(b)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
		Produced:   Produced{Discarded: []LocalID{a, b}},
	}

	ArchiveDiscardedThisWay{}.Resolve(ctx)
	if g.State.Archives[0].Count != 2 {
		t.Errorf("archives = %d, want 2", g.State.Archives[0].Count)
	}
	if g.State.Discard[0].Count != 0 {
		t.Errorf("discard = %d, want 0", g.State.Discard[0].Count)
	}

	// An empty run archives nothing.
	ctx.Produced.Discarded = nil
	ArchiveDiscardedThisWay{}.Resolve(ctx)
	if g.State.Archives[0].Count != 2 {
		t.Errorf("archives after empty run = %d, want 2", g.State.Archives[0].Count)
	}
}

// TestRevealTopOfDeckRouting covers the reveal-and-route node: the text it renders,
// the validation that rejects a bad Amount, a bad step, or a Shuffle terminal that
// is not last, and a resolve that reveals the top cards of a chosen deck, purges one,
// and shuffles the deck — from the controller's own deck, the opponent's deck, an
// empty deck, and with a declined purge.
func TestRevealTopOfDeckRouting(t *testing.T) {
	borrNit := RevealTopOfDeck{Amount: 5, ChooseWhoseDeck: true, Then: []TopAct{
		ChooseAndMove{Count: 1, Dest: IntoPurge}, Shuffle{},
	}}

	t.Run("text", func(t *testing.T) {
		want := "reveal the top 5 cards of a player's deck. Purge a card revealed this " +
			"way. Shuffle that deck"
		if got := borrNit.Text(); got != want {
			t.Errorf("text = %q, want %q", got, want)
		}
		if got := (ChooseAndMove{Count: 2, Dest: IntoPurge}).clause(); got !=
			"purge 2 cards revealed this way" {
			t.Errorf("plural purge clause = %q", got)
		}
	})

	t.Run("validate", func(t *testing.T) {
		if (RevealTopOfDeck{}).validate() == nil {
			t.Error("an Amount of 0 should be rejected")
		}
		if err := borrNit.validate(); err != nil {
			t.Errorf("validate() = %v", err)
		}
		if (ChooseAndMove{Dest: IntoPurge}).validate() == nil {
			t.Error("a Count of 0 should be rejected")
		}
		if err := (Shuffle{}).validate(); err != nil {
			t.Errorf("Shuffle validate() = %v", err)
		}
		badStep := RevealTopOfDeck{Amount: 5, ChooseWhoseDeck: true, Then: []TopAct{
			ChooseAndMove{Dest: IntoPurge},
		}}
		if badStep.validate() == nil {
			t.Error("a step with a bad Count should be rejected")
		}
		notLast := RevealTopOfDeck{Amount: 5, ChooseWhoseDeck: true, Then: []TopAct{
			Shuffle{}, ChooseAndMove{Count: 1, Dest: IntoPurge},
		}}
		if notLast.validate() == nil {
			t.Error("a Shuffle terminal that is not last should be rejected")
		}
	})

	t.Run("purges one of the revealed cards and shuffles the rest", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		top := g.AddToDeck(NewCard("Top", Logos, Creature, Common), 0)
		g.AddToDeck(NewCard("Mid", Logos, Creature, Common), 0)
		g.AddToDeck(NewCard("Low", Logos, Creature, Common), 0)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		RevealTopOfDeck{Amount: 3, ChooseWhoseDeck: true, Then: []TopAct{
			ChooseAndMove{Count: 1, Dest: IntoPurge}, Shuffle{},
		}}.Resolve(ctx)
		if purge := g.Purge(0); len(purge) != 1 || purge[0] != top {
			t.Errorf("purge = %v, want [%d]", purge, top)
		}
		if len(g.Deck(0)) != 2 {
			t.Errorf("deck = %v, want 2 cards", g.Deck(0))
		}
	})

	t.Run("an empty deck is a no-op", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		RevealTopOfDeck{Amount: 3, ChooseWhoseDeck: true, Then: []TopAct{
			ChooseAndMove{Count: 1, Dest: IntoPurge}, Shuffle{},
		}}.Resolve(ctx)
		if len(g.Purge(0)) != 0 {
			t.Errorf("purge = %v, want empty", g.Purge(0))
		}
	})

	t.Run("choosing the opponent's deck digs there", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		oppTop := g.AddToDeck(NewCard("OppTop", Logos, Creature, Common), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		g.SetChooser(0, optionPicker{idx: 1})
		RevealTopOfDeck{Amount: 3, ChooseWhoseDeck: true, Then: []TopAct{
			ChooseAndMove{Count: 1, Dest: IntoPurge}, Shuffle{},
		}}.Resolve(ctx)
		if purge := g.Purge(1); len(purge) != 1 || purge[0] != oppTop {
			t.Errorf("opponent purge = %v, want [%d]", purge, oppTop)
		}
	})

	t.Run("a declined purge purges nothing but still shuffles", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.AddToDeck(NewCard("Top", Logos, Creature, Common), 0)
		g.AddToDeck(NewCard("Mid", Logos, Creature, Common), 0)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		g.SetChooser(0, orderRejectChooser{})
		RevealTopOfDeck{Amount: 3, ChooseWhoseDeck: true, Then: []TopAct{
			ChooseAndMove{Count: 1, Dest: IntoPurge}, Shuffle{},
		}}.Resolve(ctx)
		if len(g.Purge(0)) != 0 {
			t.Errorf("purge = %v, want empty", g.Purge(0))
		}
		if len(g.Deck(0)) != 2 {
			t.Errorf("deck = %v, want 2 cards", g.Deck(0))
		}
	})

	t.Run("purging more than remain stops early", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		only := g.AddToDeck(NewCard("Only", Logos, Creature, Common), 0)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		RevealTopOfDeck{Amount: 3, ChooseWhoseDeck: true, Then: []TopAct{
			ChooseAndMove{Count: 2, Dest: IntoPurge},
		}}.Resolve(ctx)
		if purge := g.Purge(0); len(purge) != 1 || purge[0] != only {
			t.Errorf("purge = %v, want [%d]", purge, only)
		}
	})
}

// TestLookAtTopOfDeck covers the peek-and-route node: a pure peek that puts the
// looked-at cards back untouched, the routing steps that draw, archive, and discard
// them, a closing reorder, short and empty decks, declined choices, and the
// validation that rejects a bad Amount or a ReorderRest that is not last.
func TestLookAtTopOfDeck(t *testing.T) {
	t.Run("text", func(t *testing.T) {
		peek := LookAtTopOfDeck{Amount: 3}
		if got := peek.Text(); got != "look at the top 3 cards of your deck" {
			t.Errorf("peek Text() = %q", got)
		}
		if got := (LookAtTopOfDeck{Amount: 1}).Text(); got !=
			"look at the top 1 card of your deck" {
			t.Errorf("singular Text() = %q", got)
		}
		eyegor := LookAtTopOfDeck{Amount: 3, Then: []TopAct{
			ChooseAndMove{Count: 1, Dest: IntoHand}, ChooseAndMove{Count: 2, Dest: IntoDiscard},
		}}
		if got := eyegor.Text(); got !=
			"look at the top 3 cards of your deck, put 1 into your hand, and discard 2" {
			t.Errorf("Eyegor Text() = %q", got)
		}
		philo := LookAtTopOfDeck{Amount: 3, Then: []TopAct{
			ChooseAndMove{Count: 1, Dest: IntoArchives},
			ChooseAndMove{Count: 1, Dest: IntoHand},
			ChooseAndMove{Count: 1, Dest: IntoDiscard},
		}}
		if got := philo.Text(); got !=
			"look at the top 3 cards of your deck, archive 1, put 1 into your hand, and discard 1" {
			t.Errorf("Philophosaurus Text() = %q", got)
		}
		reorder := LookAtTopOfDeck{Amount: 3, Then: []TopAct{ReorderRest{}}}
		if got := reorder.Text(); got !=
			"look at the top 3 cards of your deck and put them back in any order" {
			t.Errorf("reorder Text() = %q", got)
		}
	})

	t.Run("validate", func(t *testing.T) {
		if (LookAtTopOfDeck{}).validate() == nil {
			t.Error("an Amount of 0 should be rejected")
		}
		if err := (LookAtTopOfDeck{Amount: 3}).validate(); err != nil {
			t.Errorf("validate() = %v", err)
		}
		if (ChooseAndMove{}).validate() == nil {
			t.Error("a Count of 0 should be rejected")
		}
		if err := (ReorderRest{}).validate(); err != nil {
			t.Errorf("ReorderRest validate() = %v", err)
		}
		badStep := LookAtTopOfDeck{Amount: 3, Then: []TopAct{ChooseAndMove{}}}
		if badStep.validate() == nil {
			t.Error("a step with a bad Count should be rejected")
		}
		notLast := LookAtTopOfDeck{Amount: 3, Then: []TopAct{
			ReorderRest{}, ChooseAndMove{Count: 1, Dest: IntoHand},
		}}
		if notLast.validate() == nil {
			t.Error("a ReorderRest that is not last should be rejected")
		}
		last := LookAtTopOfDeck{Amount: 3, Then: []TopAct{
			ChooseAndMove{Count: 1, Dest: IntoHand}, ReorderRest{},
		}}
		if err := last.validate(); err != nil {
			t.Errorf("a trailing ReorderRest = %v", err)
		}
	})

	t.Run("a pure peek leaves the deck untouched", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToDeck(NewCard("A Card", Logos, Creature, Common, WithPower(2)), 0)
		b := g.AddToDeck(NewCard("B Card", Logos, Tactic, Common), 0)
		LookAtTopOfDeck{Amount: 3}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Deck(0); len(got) != 2 || got[0] != a || got[1] != b {
			t.Errorf("deck = %v, want [%d %d]", got, a, b)
		}
	})

	t.Run("draws one and discards the rest", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToDeck(NewCard("A Card", Logos, Creature, Common, WithPower(2)), 0)
		b := g.AddToDeck(NewCard("B Card", Logos, Tactic, Common), 0)
		c := g.AddToDeck(NewCard("C Card", Logos, Artifact, Common), 0)
		bottom := g.AddToDeck(NewCard("Bottom", Logos, Creature, Common, WithPower(1)), 0)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{b}})
		LookAtTopOfDeck{Amount: 3, Then: []TopAct{
			ChooseAndMove{Count: 1, Dest: IntoHand}, ChooseAndMove{Count: 2, Dest: IntoDiscard},
		}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
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

	t.Run("sorts the top three into three piles", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToDeck(NewCard("A Card", Logos, Creature, Common, WithPower(2)), 0)
		b := g.AddToDeck(NewCard("B Card", Logos, Tactic, Common), 0)
		c := g.AddToDeck(NewCard("C Card", Logos, Artifact, Common), 0)
		bottom := g.AddToDeck(NewCard("Bottom", Logos, Creature, Common, WithPower(1)), 0)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{a, b}})
		LookAtTopOfDeck{Amount: 3, Then: []TopAct{
			ChooseAndMove{Count: 1, Dest: IntoArchives},
			ChooseAndMove{Count: 1, Dest: IntoHand},
			ChooseAndMove{Count: 1, Dest: IntoDiscard},
		}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Archives(0); len(got) != 1 || got[0] != a {
			t.Errorf("archives = %v, want [%d]", got, a)
		}
		if got := g.Hand(0); len(got) != 1 || got[0] != b {
			t.Errorf("hand = %v, want [%d]", got, b)
		}
		if got := g.Discard(0); len(got) != 1 || got[0] != c {
			t.Errorf("discard = %v, want [%d]", got, c)
		}
		if got := g.Deck(0); len(got) != 1 || got[0] != bottom {
			t.Errorf("deck = %v, want [%d]", got, bottom)
		}
	})

	t.Run("routes as many as remain", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToDeck(NewCard("A Card", Logos, Creature, Common, WithPower(2)), 0)
		b := g.AddToDeck(NewCard("B Card", Logos, Tactic, Common), 0)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{a}})
		LookAtTopOfDeck{Amount: 3, Then: []TopAct{
			ChooseAndMove{Count: 1, Dest: IntoHand}, ChooseAndMove{Count: 2, Dest: IntoDiscard},
		}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Hand(0); len(got) != 1 || got[0] != a {
			t.Errorf("hand = %v, want [%d]", got, a)
		}
		if got := g.Discard(0); len(got) != 1 || got[0] != b {
			t.Errorf("discard = %v, want [%d]", got, b)
		}
	})

	t.Run("an empty deck does nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		LookAtTopOfDeck{Amount: 3, Then: []TopAct{
			ChooseAndMove{Count: 1, Dest: IntoHand},
		}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Hand(0); len(got) != 0 {
			t.Errorf("hand = %v, want empty", got)
		}
	})

	t.Run("a declined draw keeps everything in the deck", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.AddToDeck(NewCard("A Card", Logos, Creature, Common, WithPower(2)), 0)
		g.AddToDeck(NewCard("B Card", Logos, Tactic, Common), 0)
		g.SetChooser(0, orderRejectChooser{})
		LookAtTopOfDeck{Amount: 3, Then: []TopAct{
			ChooseAndMove{Count: 1, Dest: IntoHand}, ChooseAndMove{Count: 2, Dest: IntoDiscard},
		}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Hand(0); len(got) != 0 {
			t.Errorf("hand = %v, want empty", got)
		}
		if got := g.Deck(0); len(got) != 2 {
			t.Errorf("deck = %v, want 2 cards", got)
		}
	})

	t.Run("reorder puts the first pick deepest and the leftover on top", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToDeck(NewCard("A Card", Logos, Creature, Common, WithPower(2)), 0)
		b := g.AddToDeck(NewCard("B Card", Logos, Tactic, Common), 0)
		c := g.AddToDeck(NewCard("C Card", Logos, Artifact, Common), 0)
		bottom := g.AddToDeck(NewCard("Bottom", Logos, Creature, Common, WithPower(1)), 0)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{c, a}})
		LookAtTopOfDeck{Amount: 3, Then: []TopAct{ReorderRest{}}}.
			Resolve(&EffectContext{Resolver: g, Controller: 0})
		// c is placed first (ends deepest of the three), a next, and b — never
		// picked — rides on top, so the draw that follows takes b.
		if got := g.Deck(0); len(got) != 4 ||
			got[0] != b || got[1] != a || got[2] != c || got[3] != bottom {
			t.Errorf("deck = %v, want [%d %d %d %d]", got, b, a, c, bottom)
		}
	})

	t.Run("reorder with fewer than two cards does nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		only := g.AddToDeck(NewCard("Only", Logos, Creature, Common, WithPower(2)), 0)
		LookAtTopOfDeck{Amount: 3, Then: []TopAct{ReorderRest{}}}.
			Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Deck(0); len(got) != 1 || got[0] != only {
			t.Errorf("deck = %v, want [%d]", got, only)
		}
	})

	t.Run("a declined reorder keeps the original order", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToDeck(NewCard("A Card", Logos, Creature, Common, WithPower(2)), 0)
		b := g.AddToDeck(NewCard("B Card", Logos, Tactic, Common), 0)
		g.SetChooser(0, orderRejectChooser{})
		LookAtTopOfDeck{Amount: 3, Then: []TopAct{ReorderRest{}}}.
			Resolve(&EffectContext{Resolver: g, Controller: 0})
		if got := g.Deck(0); len(got) != 2 || got[0] != a || got[1] != b {
			t.Errorf("deck = %v, want [%d %d]", got, a, b)
		}
	})
}

func TestDiscardTopAndForEachDiscardedHouseFilter(t *testing.T) {
	// DiscardTop text across counts and perspectives. Amount 0 folds to one card,
	// and the granted default (unset Player) names "its controller's deck".
	if got := (DiscardTop{}).Text(); got != "discard the top card of its controller's deck" {
		t.Errorf("granted default text = %q", got)
	}
	if got := (DiscardTop{Player: Controller}).Text(); got != "discard the top card of your deck" {
		t.Errorf("singular text = %q", got)
	}
	if got := (DiscardTop{Player: Opponent, Amount: 2}).Text(); got != "discard the top 2 cards of your opponent's deck" {
		t.Errorf("opponent text = %q", got)
	}
	if got := (DiscardTop{Player: EachPlayer, Amount: 2}).Text(); got != "discard the top 2 cards of each player's deck" {
		t.Errorf("each-player text = %q", got)
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
