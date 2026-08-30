package engine

import "testing"

func TestNextCreaturePlayed(t *testing.T) {
	e := NextCreaturePlayed{Of: Mars, EntersPlay: Ready{Target: Target{Kind: TargetTriggeringCreature}}}
	if e.Text() != "the next Mars creature you play this turn enters play ready" {
		t.Errorf("text = %q", e.Text())
	}
	if e.validate() != nil {
		t.Errorf("a Ready EntersPlay should validate, got %v", e.validate())
	}
	g := NewGame("A", "B", 1)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if le := g.State.Lasting[0]; g.State.LastingCount != 1 || le.Do != actReadyPlayed || le.House != Mars || !le.Once {
		t.Errorf("resolve should register a one-shot Mars ready reaction; got %+v (count %d)", le, g.State.LastingCount)
	}

	// An EntersPlay effect the flat registry cannot carry is rejected.
	if (NextCreaturePlayed{Of: Mars, EntersPlay: GainAember{Player: Controller, Amount: 1}}).validate() == nil {
		t.Error("an unsupported EntersPlay effect should be rejected")
	}
}

func TestForRemainderOfTurnText(t *testing.T) {
	cases := []struct {
		e    ForRemainderOfTurn
		want string
	}{
		{
			ForRemainderOfTurn{On: EventCreaturePlayed, Do: GainAember{Player: Controller, Amount: 1}},
			"for the remainder of the turn, each time you play a creature, gain 1 Æmber",
		},
		{
			ForRemainderOfTurn{On: EventReap, Do: GainAember{Player: Controller, Amount: 1}},
			"for the remainder of the turn, after a creature reaps, gain 1 Æmber",
		},
		{
			ForRemainderOfTurn{On: EventCreaturePlayed, Do: DealDamage{Amount: 2, Target: Target{Kind: TargetChosenEnemyCreature}}},
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
	// An unsupported Do effect.
	if err := (ForRemainderOfTurn{On: EventCreaturePlayed, Do: Draw{Amount: 1}}).validate(); err == nil {
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
	g.EndTurn(0)
	if g.State.LastingCount != 0 {
		t.Error("EndTurn should clear the reaction")
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
	ForRemainderOfTurn{On: EventCreaturePlayed, Do: DealDamage{Amount: 2, Target: Target{Kind: TargetChosenEnemyCreature}}}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})
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
