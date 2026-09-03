package engine

import "testing"

func TestNextPlayed(t *testing.T) {
	e := NextPlayed{
		Of:         Mars,
		Type:       Creature,
		EntersPlay: Ready{Target: Target{Kind: TargetTriggeringCreature}},
	}
	if e.Text() != "the next Mars creature you play this turn enters play ready" {
		t.Errorf("text = %q", e.Text())
	}
	if e.validate() != nil {
		t.Errorf("a Ready EntersPlay should validate, got %v", e.validate())
	}
	g := NewGame("A", "B", 1)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if le := g.State.Lasting[0]; g.State.LastingCount != 1 || le.Do != actReadyPlayed ||
		le.House != Mars ||
		le.Type != Creature ||
		!le.Once {
		t.Errorf(
			"resolve should register a one-shot Mars ready reaction; got %+v (count %d)",
			le,
			g.State.LastingCount,
		)
	}

	// AnyType covers both types that stay in play, and needs no house.
	anyCard := NextPlayed{
		Type:       AnyType,
		EntersPlay: Ready{Target: Target{Kind: TargetTriggeringCreature}},
	}
	want := "the next creature or artifact you play this turn enters play ready"
	if got := anyCard.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}

	// A card type is required: the effect must say what it is waiting for.
	if (NextPlayed{EntersPlay: Ready{Target: Target{Kind: TargetTriggeringCreature}}}).
		validate() == nil {
		t.Error("an unset card type should be rejected")
	}
	// An EntersPlay effect the flat registry cannot carry is rejected.
	if (NextPlayed{
		Of:         Mars,
		Type:       Creature,
		EntersPlay: GainAember{Player: Controller, Amount: 1},
	}).validate() == nil {
		t.Error("an unsupported EntersPlay effect should be rejected")
	}
}

func TestForRemainderOfTurnText(t *testing.T) {
	cases := []struct {
		e    ForRemainderOfTurn
		want string
	}{
		{
			ForRemainderOfTurn{
				On: EventCreaturePlayed,
				Do: GainAember{Player: Controller, Amount: 1},
			},
			"for the remainder of the turn, each time you play a creature, gain 1 Æmber",
		},
		{
			ForRemainderOfTurn{On: EventReap, Do: GainAember{Player: Controller, Amount: 1}},
			"for the remainder of the turn, after a creature reaps, gain 1 Æmber",
		},
		{
			ForRemainderOfTurn{
				On: EventCreaturePlayed,
				Do: DealDamage{Amount: 2, Target: Target{Kind: TargetChosenEnemyCreature}},
			},
			"for the remainder of the turn, each time you play a creature, deal 2 damage to an enemy creature",
		},
	}
	for _, c := range cases {
		if got := c.e.Text(); got != c.want {
			t.Errorf("text = %q, want %q", got, c.want)
		}
	}
}

func TestForRemainderOfTurnValidate(t *testing.T) {
	ok := ForRemainderOfTurn{On: EventCreaturePlayed, Do: GainAember{Player: Controller, Amount: 1}}
	if err := ok.validate(); err != nil {
		t.Errorf("valid reaction should pass: %v", err)
	}
	// A replacement event is not a reaction.
	if err := (ForRemainderOfTurn{On: EventReapAember, Do: GainAember{Player: Controller, Amount: 1}}).validate(); err == nil {
		t.Error("non-reaction event should fail")
	}
	// Draw is a supported Do (Library Access).
	if err := (ForRemainderOfTurn{On: EventCardPlayed, Do: Draw{Amount: 1}}).validate(); err != nil {
		t.Errorf("Draw should be a supported Do: %v", err)
	}
	// An unsupported Do effect.
	unsupported := ForRemainderOfTurn{
		On: EventCreaturePlayed,
		Do: ShuffleIntoDeck{Zones: []Zone{Discard}},
	}
	if err := unsupported.validate(); err == nil {
		t.Error("unsupported Do should fail")
	}
	// DealDamage must target an enemy creature.
	if err := (ForRemainderOfTurn{On: EventCreaturePlayed, Do: DealDamage{Amount: 2, Target: Target{Kind: TargetEachCreature}}}).validate(); err == nil {
		t.Error("DealDamage with a non-enemy target should fail")
	}
}

func TestForRemainderOfTurnGainsOnPlay(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	ForRemainderOfTurn{On: EventCreaturePlayed, Do: GainAember{Player: Controller, Amount: 1}}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.LastingCount != 1 {
		t.Fatalf("lasting count = %d, want 1", g.State.LastingCount)
	}
	g.AddToHand(testCreature("c", 3), 0)
	before := g.Aember(0)
	if _, err := g.PlayCreature(0, handIdx(g, 0, "c"), false); err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if got := g.Aember(0) - before; got != 1 {
		t.Errorf("Æmber gained on play = %d, want 1", got)
	}
	g.EndPlayPhase(0)
	if g.State.LastingCount != 0 {
		t.Error("the ready phase should clear the reaction")
	}
}

// TestForRemainderOfTurnDrawsOnCardPlayed covers Library Access: the reaction
// draws for every later card played, but not for the play that armed it.
func TestForRemainderOfTurnDrawsOnCardPlayed(t *testing.T) {
	e := ForRemainderOfTurn{On: EventCardPlayed, Do: Draw{Amount: 1}}
	want := "for the remainder of the turn, each time you play another card, draw a card"
	if got := e.Text(); got != want {
		t.Errorf("text = %q", got)
	}

	g := started(t)
	source := g.AddToBattleline(testCreature("source", 3), 0)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: source})
	g.AddToDeck(testCreature("drawn", 1), 0)
	g.AddToHand(testCreature("c", 3), 0)
	before := len(g.Hand(0))

	if _, err := g.PlayCreature(0, handIdx(g, 0, "c"), false); err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	// One card left the hand and one was drawn for it.
	if got := len(g.Hand(0)); got != before {
		t.Errorf("hand = %d, want %d (played one, drew one)", got, before)
	}
}

// TestForRemainderOfTurnExceptsItsOwnPlay checks the arming card is skipped, so
// "each time you play another card" never counts the play that installed it.
func TestForRemainderOfTurnExceptsItsOwnPlay(t *testing.T) {
	g := started(t)
	armer := g.AddToBattleline(testCreature("armer", 3), 0)
	ForRemainderOfTurn{On: EventCardPlayed, Do: Draw{Amount: 1}}.
		Resolve(&EffectContext{Resolver: g, Controller: 0, Source: armer})
	g.AddToDeck(testCreature("drawn", 1), 0)
	before := len(g.Hand(0))

	g.emitLasting(EventCardPlayed, 0, armer)

	if got := len(g.Hand(0)); got != before {
		t.Errorf("hand = %d, want %d (the arming card draws nothing for itself)", got, before)
	}
}

func TestForRemainderOfTurnGainsOnReap(t *testing.T) {
	g := started(t)
	creature := g.AddToBattleline(testCreature("c", 3), 0)
	ForRemainderOfTurn{On: EventReap, Do: GainAember{Player: Controller, Amount: 1}}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})
	before := g.Aember(0)
	g.reapWith(creature)
	if got := g.Aember(0) - before; got != 2 {
		t.Errorf("Æmber gained on reap = %d, want 2 (1 base + 1 reaction)", got)
	}
}

func TestForRemainderOfTurnDamageOnPlay(t *testing.T) {
	g := started(t)
	foe := g.AddToBattleline(testCreature("foe", 5), 1)
	ForRemainderOfTurn{
		On: EventCreaturePlayed,
		Do: DealDamage{Amount: 2, Target: Target{Kind: TargetChosenEnemyCreature}},
	}.
		Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
	g.AddToHand(testCreature("minion", 4), 0)
	if _, err := g.PlayCreature(0, handIdx(g, 0, "minion"), false); err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.Damage(foe) != 2 {
		t.Errorf("foe damage = %d, want 2", g.Damage(foe))
	}
}

func TestInstead(t *testing.T) {
	if got := (Instead{Of: EventReapAember, With: Steal}).Text(); got != "for the remainder of the turn, instead of gaining Æmber from reaping, steal the same amount" {
		t.Errorf("text = %q", got)
	}
	if err := (Instead{Of: EventReapAember, With: Steal}).validate(); err != nil {
		t.Errorf("valid replacement should pass: %v", err)
	}
	if err := (Instead{Of: EventCreaturePlayed, With: Steal}).validate(); err == nil {
		t.Error("a reaction event should fail as a replacement")
	}
	if err := (Instead{Of: EventReapAember}).validate(); err == nil {
		t.Error("an unset replacement should fail")
	}
	if err := (Instead{Of: EventAemberAddedToPool, With: Capture}).validate(); err == nil {
		t.Error("a pool event without a Player should fail")
	}
	if err := (Instead{Of: EventAemberAddedToPool, With: Capture, Player: Opponent}).validate(); err != nil {
		t.Errorf("a scoped pool replacement should pass: %v", err)
	}

	// Reaping steals instead of gaining.
	g := started(t)
	creature := g.AddToBattleline(testCreature("c", 3), 0)
	g.State.Aember[1] = 2
	Instead{Of: EventReapAember, With: Steal}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	g.reapWith(creature)
	if g.Aember(0) != 1 || g.Aember(1) != 1 {
		t.Errorf("after steal: p0=%d p1=%d, want 1/1", g.Aember(0), g.Aember(1))
	}

	// With the opponent at zero there is nothing to steal.
	g.State.Cards[creature].Exhausted = false
	g.State.Aember[1] = 0
	before := g.Aember(0)
	g.reapWith(creature)
	if g.Aember(0) != before {
		t.Errorf("no Æmber to steal: p0 = %d, want %d", g.Aember(0), before)
	}

	// Two copies replace the same single reap payout, rather than letting a reap
	// steal twice. The first active replacement consumes the whole event.
	g2 := started(t)
	creature2 := g2.AddToBattleline(testCreature("c", 3), 0)
	g2.State.Aember[1] = 3
	Instead{Of: EventReapAember, With: Steal}.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	Instead{Of: EventReapAember, With: Steal}.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	g2.reapWith(creature2)
	if g2.Aember(0) != 1 || g2.Aember(1) != 2 {
		t.Errorf("two replacements: p0=%d p1=%d, want 1/2", g2.Aember(0), g2.Aember(1))
	}
}
