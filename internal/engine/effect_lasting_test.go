package engine

import "testing"

func TestLastingActionOf(t *testing.T) {
	cases := []struct {
		do   Effect
		want lastingAction
	}{
		{GainAember{Amount: 1}, actGainAember},
		{LoseAember{Amount: 1}, actLoseAember},
		{DealDamage{Amount: 2}, actDealDamage},
		{CaptureAember{Amount: 1}, actCapture},
		{Draw{Amount: 1}, actDraw},
		{Ready{}, actReadyPlayed},
	}
	for _, c := range cases {
		if got, _, ok := lastingActionOf(c.do); !ok || got != c.want {
			t.Errorf("lastingActionOf(%T) = %v, %v; want %v", c.do, got, ok, c.want)
		}
	}
	// An effect the registry cannot carry reports not-ok.
	if _, _, ok := lastingActionOf(Heal{Amount: 1}); ok {
		t.Error("an unsupported effect should report not-ok")
	}
}

func TestNextPlayed(t *testing.T) {
	e := NextPlayed{
		Of:         namedHouse(Mars),
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
		le.House != namedHouse(Mars) ||
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
		Of:         namedHouse(Mars),
		Type:       Creature,
		EntersPlay: GainAember{Player: Controller, Amount: 1},
	}).validate() == nil {
		t.Error("an unsupported EntersPlay effect should be rejected")
	}
	// A context-dependent house matcher cannot survive the turn boundary.
	if (NextPlayed{
		Of:         HouseMatcher{Kind: MatchChosenHouse},
		Type:       Creature,
		EntersPlay: Ready{Target: Target{Kind: TargetTriggeringCreature}},
	}).validate() == nil {
		t.Error("a context-dependent house matcher should be rejected")
	}
	// A named matcher missing its house is rejected.
	if (NextPlayed{
		Of:         HouseMatcher{Kind: MatchNamedHouse},
		Type:       Creature,
		EntersPlay: Ready{Target: Target{Kind: TargetTriggeringCreature}},
	}).validate() == nil {
		t.Error("a named house matcher with no house should be rejected")
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
		Do: Shuffle{Zones: []Zone{Discard}},
	}
	if err := unsupported.validate(); err == nil {
		t.Error("unsupported Do should fail")
	}
	// DealDamage must target an enemy creature.
	if err := (ForRemainderOfTurn{On: EventCreaturePlayed, Do: DealDamage{Amount: 2, Target: Target{Kind: TargetEachCreature}}}).validate(); err == nil {
		t.Error("DealDamage with a non-enemy target should fail")
	}
}

// TestForOpponentNextTurnValidate checks the opponent-turn reaction shares the same
// gate: a reaction event with a supported Do passes, a replacement event fails.
func TestForOpponentNextTurnValidate(t *testing.T) {
	if err := (ForOpponentNextTurn{On: EventForgeKey, Do: GiveAember{All: true}}).validate(); err != nil {
		t.Errorf("valid reaction should pass: %v", err)
	}
	if err := (ForOpponentNextTurn{On: EventReapAember, Do: GiveAember{All: true}}).validate(); err == nil {
		t.Error("non-reaction event should fail")
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

	g.resolveLastingWindow(EventCardPlayed, 0, armer)

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
	if err := (Instead{Of: EventAemberTakenFromPool, With: FromCommonSupply}).validate(); err == nil {
		t.Error("a source pool event without a Player should fail")
	}
	if err := (Instead{Of: EventAemberTakenFromPool, With: FromCommonSupply, Player: Controller}).validate(); err != nil {
		t.Errorf("a scoped source replacement should pass: %v", err)
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

func TestReactionEventOf(t *testing.T) {
	if ev, ok := reactionEventOf(TriggerAfterReap); !ok || ev != EventReap {
		t.Errorf("reactionEventOf(Reap) = %v, %v; want EventReap, true", ev, ok)
	}
	if ev, ok := reactionEventOf(TriggerAfterFight); !ok || ev != EventFight {
		t.Errorf("reactionEventOf(Fight) = %v, %v; want EventFight, true", ev, ok)
	}
	if _, ok := reactionEventOf(TriggerAfterPlay); ok {
		t.Error("a non-reaction trigger should not map to a reaction event")
	}
}

func TestGainAbilityValidate(t *testing.T) {
	reap := func(e Effect) Ability { return Ability{Trigger: TriggerAfterReap, Effect: e} }

	if err := (GainAbility{Ability: reap(Draw{Amount: 1})}).validate(); err == nil {
		t.Error("an unset target should be rejected")
	}
	// A trigger the registry cannot hang a per-creature reaction on.
	if err := (GainAbility{
		Target:  Target{Kind: TargetTriggeringCreature},
		Ability: Ability{Trigger: TriggerAfterPlay, Effect: Draw{Amount: 1}},
	}).validate(); err == nil {
		t.Error("an unsupported trigger should be rejected")
	}
	// A Reap ability whose effect the flat registry cannot carry.
	if err := (GainAbility{
		Target:  Target{Kind: TargetTriggeringCreature},
		Ability: reap(ArchiveSource{}),
	}).validate(); err == nil {
		t.Error("an unsupported ability effect should be rejected")
	}
	if err := (GainAbility{
		Target:  Target{Kind: TargetTriggeringCreature},
		Ability: reap(Draw{Amount: 1}),
	}).validate(); err != nil {
		t.Errorf("valid GainAbility = %v", err)
	}
}

func TestGainAbilityText(t *testing.T) {
	e := GainAbility{
		Target:  Target{Kind: TargetTriggeringCreature},
		Ability: Ability{Trigger: TriggerAfterReap, Effect: Draw{Amount: 1}},
	}
	want := `it gains, "Reap: Draw a card."`
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

// A RemainderOfPlayerTurn grant names its window first, so a card that does not
// lean on a sibling effect to carry the duration still reads it (Adaptoid).
func TestGainAbilityRemainderOfTurnText(t *testing.T) {
	e := GainAbility{
		Target:   Target{Kind: TargetThisCreature},
		Duration: RemainderOfPlayerTurn,
		Ability:  Ability{Trigger: TriggerAfterFight, Effect: StealAember{Amount: 1}},
	}
	want := `for the remainder of the turn, ` + SelfName + ` gains, "Fight: Steal 1 Æmber."`
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

// A StartOfPlayerNextTurn grant of "Before Fight: Exalt this creature"
// (Diplomacy) validates, names its window first, and rejects the pairing with any
// other trigger.
func TestGainAbilityBeforeFightExaltValidateAndText(t *testing.T) {
	e := GainAbility{
		Target:   Target{Kind: TargetEachCreature},
		Duration: StartOfPlayerNextTurn,
		Ability: Ability{
			Trigger: TriggerBeforeFight,
			Effect:  Exalt{Target: Target{Kind: TargetThisCreature}, Amount: 1},
		},
	}
	if err := e.validate(); err != nil {
		t.Fatalf("valid Before Fight/Exalt StartOfPlayerNextTurn grant = %v", err)
	}
	want := `until the start of your next turn, each creature gains, "Before Fight: Exalt this creature."`
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}

	// A StartOfPlayerNextTurn grant is limited to a Before Fight ability.
	bad := GainAbility{
		Target:   Target{Kind: TargetEachCreature},
		Duration: StartOfPlayerNextTurn,
		Ability:  Ability{Trigger: TriggerAfterReap, Effect: Draw{Amount: 1}},
	}
	if err := bad.validate(); err == nil {
		t.Error(
			"a StartOfPlayerNextTurn grant on a non-Before-Fight trigger should be rejected",
		)
	}
}

// A StartOfPlayerNextTurn "Before Fight: Exalt this creature" grant is owned
// by the opponent so it clears at their ready phase, and it exalts whichever
// granted creature fights — friendly or enemy — regardless of whose turn it is
// (Diplomacy).
func TestGainAbilityBeforeFightExaltResolveAndFire(t *testing.T) {
	g := started(t)
	friendly := g.AddToBattleline(testCreature("friendly", 3), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 3), 1)

	GainAbility{
		Target:   Target{Kind: TargetEachCreature},
		Duration: StartOfPlayerNextTurn,
		Ability: Ability{
			Trigger: TriggerBeforeFight,
			Effect:  Exalt{Target: Target{Kind: TargetThisCreature}, Amount: 1},
		},
	}.Resolve(&EffectContext{Resolver: g, Source: friendly, Controller: 0})

	if g.State.LastingCount != 2 {
		t.Fatalf("lasting count = %d, want 2 (one per creature)", g.State.LastingCount)
	}
	for i := 0; i < int(g.State.LastingCount); i++ {
		le := g.State.Lasting[i]
		if le.On != EventBeforeFight || le.Do != actExalt || !le.HasSubject {
			t.Fatalf("entry %d = %+v, want a Subject-scoped before-fight exalt", i, le)
		}
		// Owned by the opponent (1) so it survives the caster's own turn end.
		if le.Controller != 1 {
			t.Errorf("entry %d controller = %d, want 1 (opponent-owned)", i, le.Controller)
		}
	}

	// The caster's ready phase does not clear an opponent-owned grant.
	g.clearLasting(0)
	if g.State.LastingCount != 2 {
		t.Fatalf("caster's ready cleared the grant, count = %d, want 2", g.State.LastingCount)
	}

	// Each granted creature exalts itself when it fights, whoever controls it.
	g.fireLastingBeforeFight(friendly)
	g.fireLastingBeforeFight(enemy)
	if got := g.AmberOn(friendly); got != 1 {
		t.Errorf("friendly Æmber-on-card = %d, want 1", got)
	}
	if got := g.AmberOn(enemy); got != 1 {
		t.Errorf("enemy Æmber-on-card = %d, want 1", got)
	}

	// The opponent's ready phase lifts the grant.
	g.clearLasting(1)
	if g.State.LastingCount != 0 {
		t.Fatalf("opponent's ready did not lift the grant, count = %d", g.State.LastingCount)
	}
}

// actExalt describes itself for the ordering prompt, even though a before-fight
// grant fires through fireLastingBeforeFight rather than resolveLastingWindow.
func TestActExaltDescribe(t *testing.T) {
	if got, want := actExalt.describe(), "exalt the creature"; got != want {
		t.Errorf("actExalt.describe() = %q, want %q", got, want)
	}
}

// GainAbility can grant a Fight reaction that readies the fighting creature, and
// renders the granted self-reference as "this creature" rather than the source
// card's name (Into the Fray).
func TestGainAbilityFightReady(t *testing.T) {
	e := GainAbility{
		Target: Target{Kind: TargetTriggeringCreature},
		Ability: Ability{
			Trigger: TriggerAfterFight,
			Effect:  Ready{Target: Target{Kind: TargetThisCreature}},
		},
	}
	if err := e.validate(); err != nil {
		t.Fatalf("valid Fight/Ready GainAbility = %v", err)
	}
	if got, want := e.Text(), `it gains, "Fight: Ready this creature."`; got != want {
		t.Errorf("text = %q, want %q", got, want)
	}

	g := started(t)
	granted := g.AddToBattleline(testCreature("granted", 3), 0)
	GainAbility{
		Target:  Target{Kind: TargetThisCreature},
		Ability: e.Ability,
	}.Resolve(&EffectContext{Resolver: g, Source: granted, Controller: 0})

	le := g.State.Lasting[0]
	if le.On != EventFight || le.Do != actReadyPlayed || le.Subject != granted {
		t.Fatalf("registered reaction = %+v, want a Subject-scoped fight ready", le)
	}

	g.State.Cards[granted].Exhausted = true
	g.resolveLastingWindow(EventFight, 0, granted)
	if g.State.Cards[granted].Exhausted {
		t.Error("the granted creature should be readied after it fights")
	}
}

// TestGainAbilityFightSteal grants a "Fight: Steal 1 Æmber" ability, the Adaptoid
// option: the creature registers a subject-scoped fight steal that moves Æmber
// from the opponent's pool to its controller's when it fights.
func TestGainAbilityFightSteal(t *testing.T) {
	e := GainAbility{
		Target: Target{Kind: TargetThisCreature},
		Ability: Ability{
			Trigger: TriggerAfterFight,
			Effect:  StealAember{Amount: 1},
		},
	}
	if err := e.validate(); err != nil {
		t.Fatalf("valid Fight/Steal GainAbility = %v", err)
	}
	if got, want := e.Text(), SelfName+` gains, "Fight: Steal 1 Æmber."`; got != want {
		t.Errorf("text = %q, want %q", got, want)
	}

	g := started(t)
	granted := g.AddToBattleline(testCreature("granted", 3), 0)
	g.SetAember(1, 2)
	e.Resolve(&EffectContext{Resolver: g, Source: granted, Controller: 0})

	le := g.State.Lasting[0]
	if le.On != EventFight || le.Do != actSteal || le.Subject != granted {
		t.Fatalf("registered reaction = %+v, want a Subject-scoped fight steal", le)
	}

	g.resolveLastingWindow(EventFight, 0, granted)
	if got := g.State.Aember[0]; got != 1 {
		t.Errorf("controller Æmber after the granted creature fights = %d, want 1", got)
	}
	if got := g.State.Aember[1]; got != 1 {
		t.Errorf("opponent Æmber after the steal = %d, want 1", got)
	}
}

// GainAbility scopes the granted reaction to the one creature it targets: that
// creature's reap draws, but another creature's reap does not.
func TestGainAbilityResolveSubjectScoped(t *testing.T) {
	g := started(t)
	granted := g.AddToBattleline(testCreature("granted", 3), 0)
	other := g.AddToBattleline(testCreature("other", 3), 0)
	g.AddToDeck(testCreature("d1", 1), 0)
	g.AddToDeck(testCreature("d2", 1), 0)

	GainAbility{
		Target:  Target{Kind: TargetThisCreature},
		Ability: Ability{Trigger: TriggerAfterReap, Effect: Draw{Amount: 1}},
	}.Resolve(&EffectContext{Resolver: g, Source: granted, Controller: 0})

	if g.State.LastingCount != 1 {
		t.Fatalf("lasting count = %d, want 1", g.State.LastingCount)
	}
	le := g.State.Lasting[0]
	if le.On != EventReap || le.Do != actDraw || !le.HasSubject || le.Subject != granted {
		t.Fatalf(
			"registered reaction = %+v, want a Subject-scoped reap draw on the granted creature",
			le,
		)
	}

	// The granted creature's reap draws.
	before := len(g.Hand(0))
	g.resolveLastingWindow(EventReap, 0, granted)
	if got := len(g.Hand(0)); got != before+1 {
		t.Errorf("hand after the granted creature reaps = %d, want %d", got, before+1)
	}

	// A different creature's reap does not (the Subject filter skips it).
	before = len(g.Hand(0))
	g.resolveLastingWindow(EventReap, 0, other)
	if got := len(g.Hand(0)); got != before {
		t.Errorf("hand after another creature reaps = %d, want %d (no draw)", got, before)
	}
}

func TestDamageOthersAfterUsingTraitText(t *testing.T) {
	want := "for the remainder of the turn, after you use a Dinosaur creature, " +
		"deal 1 damage to each non-Dinosaur creature"
	got := DamageOthersAfterUsingTrait{Trait: Dinosaur, Amount: 1}.Text()
	if got != want {
		t.Errorf("text = %q", got)
	}
}

func TestDamageOthersAfterUsingTraitValidate(t *testing.T) {
	if (DamageOthersAfterUsingTrait{Trait: Dinosaur, Amount: 1}).validate() != nil {
		t.Error("valid effect should not error")
	}
	if (DamageOthersAfterUsingTrait{Amount: 1}).validate() == nil {
		t.Error("missing trait should error")
	}
	if (DamageOthersAfterUsingTrait{Trait: Dinosaur}).validate() == nil {
		t.Error("non-positive amount should error")
	}
}

// marchBoard sets up a game with Legion's March armed for player 0 and a spread of
// Dinosaur and non-Dinosaur creatures on both battlelines, returning their ids.
func marchBoard(t *testing.T) (g *Game, dino, ally, allyDino, enemy, enemyDino LocalID) {
	t.Helper()
	g = started(t)
	dino = g.AddToBattleline(testCreature("dino", 6, WithTraits(Dinosaur)), 0)
	ally = g.AddToBattleline(testCreature("ally", 3), 0)
	allyDino = g.AddToBattleline(testCreature("allyDino", 3, WithTraits(Dinosaur)), 0)
	enemy = g.AddToBattleline(testCreature("enemy", 3), 1)
	enemyDino = g.AddToBattleline(testCreature("enemyDino", 3, WithTraits(Dinosaur)), 1)
	DamageOthersAfterUsingTrait{Trait: Dinosaur, Amount: 1}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, Source: dino},
	)
	return g, dino, ally, allyDino, enemy, enemyDino
}

// assertMarchFired checks the non-Dinosaurs on both sides took 1 damage and the
// Dinosaurs (including the one used) were spared.
func assertMarchFired(t *testing.T, g *Game, dino, ally, allyDino, enemy, enemyDino LocalID) {
	t.Helper()
	if got := g.Damage(ally); got != 1 {
		t.Errorf("friendly non-Dinosaur damage = %d, want 1", got)
	}
	if got := g.Damage(enemy); got != 1 {
		t.Errorf("enemy non-Dinosaur damage = %d, want 1", got)
	}
	if got := g.Damage(dino); got != 0 {
		t.Errorf("used Dinosaur damage = %d, want 0", got)
	}
	if got := g.Damage(allyDino); got != 0 {
		t.Errorf("friendly Dinosaur damage = %d, want 0", got)
	}
	if got := g.Damage(enemyDino); got != 0 {
		t.Errorf("enemy Dinosaur damage = %d, want 0", got)
	}
}

// Reaping a Dinosaur fires the march: non-Dinosaurs on both sides take 1.
func TestDamageOthersAfterUsingTraitReap(t *testing.T) {
	g, dino, ally, allyDino, enemy, enemyDino := marchBoard(t)
	if err := g.Reap(0, dino); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	assertMarchFired(t, g, dino, ally, allyDino, enemy, enemyDino)
}

// Fighting with a Dinosaur fires the march for the bystanders.
func TestDamageOthersAfterUsingTraitFight(t *testing.T) {
	g, dino, ally, allyDino, enemy, enemyDino := marchBoard(t)
	target := g.AddToBattleline(testCreature("target", 1, WithTraits(Dinosaur)), 1)
	if err := g.Fight(0, dino, target); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	// The combatant dino takes 1 combat damage from the target, not from the march;
	// the march spares it and every other Dinosaur, and hits the non-Dinosaurs.
	if got := g.Damage(ally); got != 1 {
		t.Errorf("friendly non-Dinosaur damage = %d, want 1", got)
	}
	if got := g.Damage(enemy); got != 1 {
		t.Errorf("enemy non-Dinosaur damage = %d, want 1", got)
	}
	if got := g.Damage(allyDino); got != 0 {
		t.Errorf("friendly Dinosaur damage = %d, want 0", got)
	}
	if got := g.Damage(enemyDino); got != 0 {
		t.Errorf("enemy Dinosaur bystander damage = %d, want 0", got)
	}
}

// Using a Dinosaur's action ability fires the march.
func TestDamageOthersAfterUsingTraitAction(t *testing.T) {
	g := started(t)
	actor := g.AddToBattleline(
		testCreature("actor", 6, WithTraits(Dinosaur),
			WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 1})),
		0,
	)
	ally := g.AddToBattleline(testCreature("ally", 3), 0)
	DamageOthersAfterUsingTrait{Trait: Dinosaur, Amount: 1}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, Source: actor},
	)
	if err := g.UseAction(0, actor); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if got := g.Damage(ally); got != 1 {
		t.Errorf("non-Dinosaur damage = %d, want 1", got)
	}
}

// Using a stunned Dinosaur spends the use to recover the stun, and that still
// counts as using it, so the march fires.
func TestDamageOthersAfterUsingTraitUnstun(t *testing.T) {
	g, dino, ally, allyDino, enemy, enemyDino := marchBoard(t)
	g.State.Cards[dino].Stunned = true
	if err := g.Reap(0, dino); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if g.State.Cards[dino].Stunned {
		t.Error("Dinosaur should have recovered from the stun")
	}
	assertMarchFired(t, g, dino, ally, allyDino, enemy, enemyDino)
}

// Using a non-Dinosaur does not fire the march.
func TestDamageOthersAfterUsingTraitTraitGate(t *testing.T) {
	g, _, ally, _, _, _ := marchBoard(t)
	plain := g.AddToBattleline(testCreature("plain", 3), 0)
	if err := g.Reap(0, plain); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if got := g.Damage(ally); got != 0 {
		t.Errorf("non-Dinosaur use should not fire the march, damage = %d", got)
	}
}

func TestReturnNextActionToHandText(t *testing.T) {
	want := "after you resolve your next tactic this turn, put it into your " +
		"hand instead of your discard pile"
	if got := (ReturnNextActionToHand{}).Text(); got != want {
		t.Errorf("text = %q", got)
	}
}

// Resolve arms a one-shot redirect the action-play path later consumes.
func TestReturnNextActionToHandArms(t *testing.T) {
	g := started(t)
	ReturnNextActionToHand{}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if !g.consumeNextActionToHand(0) {
		t.Fatal("redirect was not armed")
	}
	if g.consumeNextActionToHand(0) {
		t.Fatal("redirect should be consumed after one read")
	}
}

// A redirect owned by one player is not consumed by the other.
func TestReturnNextActionToHandScopedToController(t *testing.T) {
	g := started(t)
	ReturnNextActionToHand{}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.consumeNextActionToHand(1) {
		t.Fatal("opponent should not consume the redirect")
	}
	if !g.consumeNextActionToHand(0) {
		t.Fatal("controller should consume its own redirect")
	}
}

// With the redirect armed, the next action card the controller plays returns to
// their hand instead of their discard pile, and only that one card.
func TestReturnNextActionRedirectsPlayedAction(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.StartTurn(0)
	first := g.AddToHand(NewCard("First Action", Sanctum, Tactic, Common), 0)
	second := g.AddToHand(NewCard("Second Action", Sanctum, Tactic, Common), 0)
	g.AddLasting(LastingEffect{On: EventNextActionToHand, Do: actReturnToHand, Once: true})

	if err := g.PlayAction(0, handIdxByID(g, 0, first)); err != nil {
		t.Fatalf("PlayAction first: %v", err)
	}
	if !g.State.Hand[0].contains(first) {
		t.Error("first action should have returned to hand")
	}
	if g.State.Discard[0].contains(first) {
		t.Error("first action should not be in the discard pile")
	}

	// The redirect fired once: the next action discards normally.
	if err := g.PlayAction(0, handIdxByID(g, 0, second)); err != nil {
		t.Fatalf("PlayAction second: %v", err)
	}
	if !g.State.Discard[0].contains(second) {
		t.Error("second action should discard normally after the redirect is spent")
	}
	if g.State.Hand[0].contains(second) {
		t.Error("second action should not return to hand")
	}
}

func TestTakesExtraDamageText(t *testing.T) {
	e := TakesExtraDamage{Target: Target{Kind: TargetChosenCreature}, Amount: 2}
	want := "for the remainder of the turn, whenever a creature takes damage, " +
		"it takes an additional 2 damage"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
}

func TestTakesExtraDamageValidate(t *testing.T) {
	if err := validateEffect(
		TakesExtraDamage{Target: Target{Kind: TargetChosenCreature}, Amount: 2},
	); err != nil {
		t.Errorf("valid effect rejected: %v", err)
	}
	if validateEffect(TakesExtraDamage{Amount: 2}) == nil {
		t.Error("want error for missing target")
	}
	if validateEffect(
		TakesExtraDamage{Target: Target{Kind: TargetChosenCreature}, Amount: 0},
	) == nil {
		t.Error("want error for non-positive amount")
	}
}

func TestTakesExtraDamageAugmentsDamage(t *testing.T) {
	g := started(t)
	victim := g.AddToBattleline(testCreature("victim", 10), 0)
	TakesExtraDamage{Target: Target{Kind: TargetThisCreature}, Amount: 2}.Resolve(
		&EffectContext{Resolver: g, Source: victim, Controller: 0},
	)
	// A single instance of 3 damage lands as 3 + 2 = 5.
	g.dealDamage(0, DamageTarget{ID: victim, Amount: 3})
	if got := g.Damage(victim); got != 5 {
		t.Errorf("damage = %d, want 5 (3 + 2 bonus)", got)
	}
	// A creature with no augmentation takes only what it is dealt.
	bare := g.AddToBattleline(testCreature("bare", 10), 0)
	g.dealDamage(0, DamageTarget{ID: bare, Amount: 3})
	if got := g.Damage(bare); got != 3 {
		t.Errorf("bare damage = %d, want 3", got)
	}
}

func TestLastingExtraDamageQuery(t *testing.T) {
	g := started(t)
	a := g.AddToBattleline(testCreature("a", 5), 0)
	b := g.AddToBattleline(testCreature("b", 5), 0)
	g.AddLasting(
		LastingEffect{
			On:         EventCreatureTakesDamage,
			Do:         actTakeExtraDamage,
			Amount:     2,
			Subject:    a,
			HasSubject: true,
		},
	)
	if got := g.lastingExtraDamage(a); got != 2 {
		t.Errorf("extra damage on a = %d, want 2", got)
	}
	if got := g.lastingExtraDamage(b); got != 0 {
		t.Errorf("extra damage on b = %d, want 0", got)
	}
}
