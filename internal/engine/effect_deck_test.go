package engine

import (
	"strings"
	"testing"
)

func TestChaosPortalComposition(t *testing.T) {
	effect := ChooseHouseThen{Then: Sequence{Effects: []Effect{
		Sentence{Effect: RevealTopOfDeck{}},
		Conditional{Cond: ItIsOfHouse{House: TheChosenHouse}, Then: PlayRevealedCard{}},
	}}}
	if got := effect.Text(); got != "choose a house - reveal the top card of your deck. If it is of the chosen house, play it" {
		t.Errorf("text = %q", got)
	}

	g := started(t) // Brobnar is active; the chosen house can still play a Logos card.
	top := g.AddToDeck(NewCard("Portal Scout", Logos, Creature, Common, WithPower(2)), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0, ChosenHouse: Logos}

	Sequence{Effects: []Effect{RevealTopOfDeck{}, Conditional{Cond: ItIsOfHouse{House: TheChosenHouse}, Then: PlayRevealedCard{}}}}.Resolve(ctx)

	if g.State.Deck[0].Count != 0 {
		t.Errorf("deck count = %d, want 0", g.State.Deck[0].Count)
	}
	if got := g.Battleline(0); len(got) != 1 || got[0] != top {
		t.Errorf("battleline = %v, want [%d]", got, top)
	}
	if !g.State.Cards[top].Exhausted {
		t.Error("played creature should enter exhausted")
	}
	if g.State.CardsPlayedThisTurn[0] != 1 || g.State.CardsPlayedByHouseThisTurn[0][Logos] != 1 {
		t.Errorf("play counts = %d/%d, want 1/1",
			g.State.CardsPlayedThisTurn[0], g.State.CardsPlayedByHouseThisTurn[0][Logos])
	}
	if len(g.Log) == 0 || !strings.Contains(g.Log[len(g.Log)-1], "Portal Scout") {
		t.Errorf("log = %v, want the top card revealed and played", g.Log)
	}
}

func TestChaosPortalMissesAndGuards(t *testing.T) {
	g := started(t)
	top := g.AddToDeck(NewCard("Wrong House", Dis, Tactic, Common), 0)
	reveal := Sequence{Effects: []Effect{RevealTopOfDeck{}, Conditional{Cond: ItIsOfHouse{House: TheChosenHouse}, Then: PlayRevealedCard{}}}}
	reveal.Resolve(&EffectContext{Resolver: g, Controller: 0, ChosenHouse: Logos})
	if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != top {
		t.Errorf("non-matching top card moved: deck = %v, want [%d]", g.State.Deck[0].slice(), top)
	}
	if g.State.CardsPlayedThisTurn[0] != 0 {
		t.Errorf("non-matching card should not count as played, got %d", g.State.CardsPlayedThisTurn[0])
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
					g.AddToBattleline(NewCard("Imp", Brobnar, Creature, Common, WithPower(1),
						WithRestrictions(Restrictions{PlayCardLimit: PlayCardLimit{Player: Controller, Amount: 1}})), 0)
					g.State.CardsPlayedThisTurn[0] = 1
				},
				err: ErrCardPlayLimit,
			},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				g := started(t)
				top := g.AddToDeck(NewCard("Portal Scout", Logos, Creature, Common, WithPower(2)), 0)
				tc.setup(g)
				if _, err := g.playCardFromZone(0, top, func() { g.State.Deck[0].removeAt(0) }, playCardOptions{}); err != tc.err {
					t.Fatalf("playCardFromZone err = %v, want %v", err, tc.err)
				}
				if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != top {
					t.Errorf("guarded play should leave top card in deck, got %v", g.State.Deck[0].slice())
				}
			})
		}

	})

	t.Run("creature play restriction", func(t *testing.T) {
		g := started(t)
		g.AddToBattleline(NewCard("Blocker", Brobnar, Creature, Common, WithPower(1),
			WithRestrictions(Restrictions{CannotPlay: Creature})), 0)
		top := g.AddToDeck(NewCard("Portal Scout", Logos, Creature, Common, WithPower(2)), 0)
		if _, err := g.playCardFromZone(0, top, func() { g.State.Deck[0].removeAt(0) }, playCardOptions{}); err != ErrCannotPlayCreature {
			t.Fatalf("playCardFromZone err = %v, want %v", err, ErrCannotPlayCreature)
		}
		if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != top {
			t.Errorf("barred creature should stay in deck, got %v", g.State.Deck[0].slice())
		}
	})

	t.Run("unknown card type", func(t *testing.T) {
		g := started(t)
		top := g.AddToDeck(NewCard("Mystery", Logos, CardType("Mystery"), Common), 0)
		if _, err := g.playCardFromZone(0, top, func() { g.State.Deck[0].removeAt(0) }, playCardOptions{}); err != ErrWrongType {
			t.Fatalf("playCardFromZone err = %v, want %v", err, ErrWrongType)
		}
		if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != top {
			t.Errorf("unknown type should stay in deck, got %v", g.State.Deck[0].slice())
		}
	})
}

func TestBonkersComposition(t *testing.T) {
	g := started(t)
	source := g.AddArtifact(NewCard("Bonkers Killing Machine", Logos, Artifact, Rare), 0)
	p1Top := g.AddToDeck(NewCard("Mars Top", Mars, Tactic, Common), 0)
	p2Top := g.AddToDeck(NewCard("Dis Top", Dis, Tactic, Common), 1)
	marsCreature := g.AddToBattleline(NewCard("Mars Creature", Mars, Creature, Common, WithPower(4)), 0)
	disArtifact := g.AddArtifact(NewCard("Dis Artifact", Dis, Artifact, Common), 1)
	bystander := g.AddToBattleline(testCreature("bystander", 4), 1)
	ctx := &EffectContext{Resolver: g, Source: source, Controller: 0}

	effect := Sequence{Effects: []Effect{
		Sentence{Effect: DiscardTopOfEachDeck{}},
		Sentence{Effect: ForEachDiscarded{Do: Destroy{Target: Target{Kind: TargetChosenInPlay}.OfContextualHouse()}}},
		Conditional{Cond: CardsDestroyedFewerThan{Amount: 2}, Then: Destroy{Target: Target{Kind: TargetThisCreature}}},
	}}
	if got := effect.Text(); got != "discard the top card of each player's deck. For each card discarded this way, destroy a creature or artifact of that card's house. If fewer than 2 cards are destroyed this way, destroy {self}" {
		t.Errorf("text = %q", got)
	}

	effect.Resolve(ctx)

	if discard := g.Discard(0); len(discard) != 2 || discard[0] != p1Top || discard[1] != marsCreature {
		t.Errorf("controller discard = %v, want top card then Mars creature", discard)
	}
	if discard := g.Discard(1); len(discard) != 2 || discard[0] != p2Top || discard[1] != disArtifact {
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

	Sequence{Effects: []Effect{
		Sentence{Effect: DiscardTopOfEachDeck{}},
		Sentence{Effect: ForEachDiscarded{Do: Destroy{Target: Target{Kind: TargetChosenInPlay}.OfContextualHouse()}}},
		Conditional{Cond: CardsDestroyedFewerThan{Amount: 2}, Then: Destroy{Target: Target{Kind: TargetThisCreature}}},
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
	if err := validateEffect(ForEachDiscarded{Do: Destroy{Target: Target{Kind: TargetChosenInPlay}.OfContextualHouse()}}); err != nil {
		t.Errorf("valid ForEachDiscarded = %v", err)
	}

	// Text renders the chosen-in-play noun and the contextual-house clause.
	target := Target{Kind: TargetChosenInPlay}.OfContextualHouse()
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
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if !resolverInPlay(ctx, live) {
		t.Error("a battleline creature should read in play")
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

	e := Sequence{Effects: []Effect{
		Sentence{Effect: DiscardTopOfDeck{}},
		Conditional{Cond: ItIsOfHouse{House: TheActiveHouse}, Then: CancelFight{}},
	}}
	if got := e.Text(); got != "discard the top card of its controller's deck. If it is of the active house, the fight does not occur" {
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
