package engine

import (
	"slices"
	"testing"
)

// These tests exercise destruction — the KeyForge simultaneous "Destroyed:"
// timing — and where a card and its upgrades and Æmber go as it leaves play.

// recordEnemyCount is a test effect that records how many enemy creatures are in
// play at the moment it resolves.
type recordEnemyCount struct{ got *int }

func (recordEnemyCount) Text() string { return "record enemy count" }
func (e recordEnemyCount) Resolve(ctx *EffectContext) {
	*e.got = len(ctx.Resolver.Battleline(ctx.Opponent()))
}

// TestDestroyOrderByCreature destroys two creatures that each carry one Destroyed
// ability: the controller resolves them one creature at a time (a lone remaining
// creature is forced), and both abilities fire.
func TestDestroyOrderByCreature(t *testing.T) {
	g := started(t)
	gain := GainAember{Player: Controller, Amount: 1}
	a := g.AddToBattleline(testCreature("a", 3, WithAbility(TriggerDestroyed, gain)), 0)
	b := g.AddToBattleline(testCreature("b", 3, WithAbility(TriggerDestroyed, gain)), 0)
	before := g.Aember(0)

	g.DestroyEach(0, []LocalID{a, b})

	if g.Aember(0) != before+2 {
		t.Errorf("aember = %d, want %d (both Destroyed abilities resolve)", g.Aember(0), before+2)
	}
}

// TestSaveFromDestruction covers a creature's own "Destroyed:" replacement: it
// marks itself saved so the discard step leaves it in play, and its Do resolves
// (Reassembling Automaton).
func TestSaveFromDestruction(t *testing.T) {
	e := SaveFromDestruction{Do: GainAember{Player: Controller, Amount: 1}}
	if got := e.Text(); got != "instead of destroying "+SelfName+", gain 1 Æmber" {
		t.Errorf("text = %q", got)
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v", err)
	}

	g := started(t)
	before := g.Aember(0)
	saved := g.AddToBattleline(testCreature("automaton", 3,
		WithAbility(TriggerDestroyed, e)), 0)

	g.DestroyEach(0, []LocalID{saved})

	if !slices.Contains(g.Battleline(0), saved) {
		t.Error("the saved creature should remain in play")
	}
	if g.Aember(0) != before+1 {
		t.Errorf("aember = %d, want %d (the save's Do resolved)", g.Aember(0), before+1)
	}
}

// orderMark is a test Destroyed effect that appends its tag to a shared log when
// it resolves, so a test can read back the resolution order. Two marks with the
// same tag render identical text, so orderTriggered treats them as one ability.
type orderMark struct {
	log *[]string
	tag string
}

func (m orderMark) Text() string             { return "mark " + m.tag }
func (m orderMark) Resolve(_ *EffectContext) { *m.log = append(*m.log, m.tag) }

// countingReverseChooser records how many ordering picks it is asked to make and
// picks the last reaction each time, reversing the window, so a test can prove
// identical abilities are never prompted while distinct ones are ordered in full.
type countingReverseChooser struct{ picks int }

func (c *countingReverseChooser) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	return cands[len(cands)-1], true
}

func (c *countingReverseChooser) ChooseReaction(_ string, reactions []OrderableReaction) int {
	c.picks++
	return len(reactions) - 1
}

// TestOrderTriggeredOrdersEveryDistinctWindow checks that a window carrying a
// distinct ability lets the player order it against the rest, and stops prompting
// the moment only identical abilities remain — since those resolve the same in any
// order. Two creatures carry the same mark and a third a different one; the
// reversing chooser resolves the distinct mark first, after which the identical
// pair is auto-resolved with no further prompt (one pick in all).
func TestOrderTriggeredOrdersEveryDistinctWindow(t *testing.T) {
	g := started(t)
	var log []string
	ch := &countingReverseChooser{}
	g.SetChooser(0, ch)
	a := g.AddToBattleline(
		testCreature("a", 3, WithAbility(TriggerDestroyed, orderMark{&log, "same"})),
		0,
	)
	b := g.AddToBattleline(
		testCreature("b", 3, WithAbility(TriggerDestroyed, orderMark{&log, "same"})),
		0,
	)
	c := g.AddToBattleline(
		testCreature("c", 3, WithAbility(TriggerDestroyed, orderMark{&log, "diff"})),
		0,
	)

	g.DestroyEach(0, []LocalID{a, b, c})

	if ch.picks != 1 {
		t.Errorf(
			"ordering picks = %d, want 1 (the distinct ability is ordered; the identical remainder is not prompted)",
			ch.picks,
		)
	}
	// The chooser reverses, so the distinct mark resolves first, then the identical
	// pair in their gathered order.
	want := []string{"diff", "same", "same"}
	if !slices.Equal(log, want) {
		t.Errorf("resolution order = %v, want %v", log, want)
	}
}

// TestOrderTriggeredAllIdenticalNeverPrompts checks that a window whose abilities
// are all identical is auto-ordered with no pick at all.
func TestOrderTriggeredAllIdenticalNeverPrompts(t *testing.T) {
	g := started(t)
	var log []string
	ch := &countingReverseChooser{}
	g.SetChooser(0, ch)
	a := g.AddToBattleline(
		testCreature("a", 3, WithAbility(TriggerDestroyed, orderMark{&log, "same"})),
		0,
	)
	b := g.AddToBattleline(
		testCreature("b", 3, WithAbility(TriggerDestroyed, orderMark{&log, "same"})),
		0,
	)

	g.DestroyEach(0, []LocalID{a, b})

	if ch.picks != 0 {
		t.Errorf("ordering picks = %d, want 0 (identical abilities are never prompted)", ch.picks)
	}
	if len(log) != 2 {
		t.Errorf("resolved %d abilities, want 2", len(log))
	}
}

func TestDestroyTogetherResolvesBeforeDiscard(t *testing.T) {
	g := started(t)
	// A's "Destroyed:" ability records how many enemy creatures are in play when
	// it fires. KeyForge tags both for destruction and resolves the Destroyed
	// abilities before moving anything to the discard, so B is still present.
	var enemiesWhenADied int
	a := g.AddToBattleline(
		testCreature(
			"a",
			3,
			WithAbility(TriggerDestroyed, recordEnemyCount{got: &enemiesWhenADied}),
		),
		0,
	)
	b := g.AddToBattleline(testCreature("b", 3), 1)

	g.DestroyEach(0, []LocalID{a, b})

	if enemiesWhenADied != 1 {
		t.Errorf(
			"A's Destroyed ability saw %d enemy creatures; want 1 (still in play until discard)",
			enemiesWhenADied,
		)
	}
	if g.inPlay(a) || g.inPlay(b) {
		t.Error("both creatures should be in the discard after the event")
	}
}

func TestDestroyedRelocationSkipsDiscard(t *testing.T) {
	g := started(t)
	// A creature whose "Destroyed:" ability returns it to the top of its deck
	// leaves play during the event, so it is not also moved to the discard.
	c := g.AddToBattleline(
		testCreature(
			"wanderer",
			3,
			WithAbility(
				TriggerDestroyed,
				PutFromPlay{Target: Target{Kind: TargetThisCreature}, Destination: ToTopOfDeck},
			),
		),
		0,
	)
	g.DestroyEach(0, []LocalID{c})

	if g.inPlay(c) {
		t.Error("creature should have left play")
	}
	if len(g.Discard(0)) != 0 {
		t.Errorf("relocated creature should not be discarded; discard size = %d", len(g.Discard(0)))
	}
	if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != c {
		t.Error("creature should be on top of its owner's deck")
	}
}

func TestDestroyMovesUpgradesToDiscard(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 1), 0)
	g.AddToHand(exBruteStrength(), 0)
	if _, err := g.PlayUpgrade(0, 0); err != nil {
		t.Fatal(err)
	}
	// Host now has 6 power (1 + 5). Kill it with enough damage.
	g.DealDamage(0, []DamageTarget{{ID: host, Amount: 6}})
	if g.inPlay(host) {
		t.Error("host should be destroyed")
	}
	if len(g.Discard(0)) != 2 { // host + upgrade
		t.Errorf("discard size = %d, want 2", len(g.Discard(0)))
	}
}

func TestUpgradeCanPreventHostDestructionOnce(t *testing.T) {
	g := started(t)
	shield := NewCard(
		"Shield",
		Sanctum,
		Upgrade,
		Rare,
		WithStatic(
			StaticModifier{
				Replaces: Replace{When: EventCreatureDestroyed, With: Sequence{Effects: []Effect{
					Heal{Fully: true, Target: Target{Kind: TargetTriggeringCreature}},
					Destroy{Target: Target{Kind: TargetThisCreature}},
				}}},
			},
		),
	)
	host := g.AddToBattleline(testCreature("host", 3,
		WithAbility(TriggerDestroyed, GainAember{Player: Controller, Amount: 1})), 0)
	attachUpgrade(g, host, shield)
	upgrade := g.Upgrades(host)[0]
	g.State.Cards[host].Damage = 2

	g.DealDamage(0, []DamageTarget{{ID: host, Amount: 1}})

	if !g.inPlay(host) {
		t.Fatal("host should stay in play when its destruction is prevented")
	}
	if g.Damage(host) != 0 {
		t.Errorf("host damage = %d, want fully healed", g.Damage(host))
	}
	if got := len(g.Upgrades(host)); got != 0 {
		t.Errorf("host upgrade count = %d, want the shield consumed", got)
	}
	if got := g.Discard(0); len(got) != 1 || got[0] != upgrade {
		t.Errorf("discard = %v, want only the consumed upgrade", got)
	}
	if g.Aember(0) != 0 {
		t.Errorf("destroyed ability fired despite prevention; Æmber = %d, want 0", g.Aember(0))
	}

	g.DealDamage(0, []DamageTarget{{ID: host, Amount: 3}})

	if g.inPlay(host) {
		t.Fatal("host should be destroyed by the next lethal hit")
	}
	if got := g.Discard(0); len(got) != 2 || got[1] != host {
		t.Errorf("discard = %v, want consumed upgrade then host", got)
	}
}

func TestDestroyAttachedUpgradeIgnoresUnattached(t *testing.T) {
	g := started(t)
	up := g.Register(NewCard("Loose", Sanctum, Upgrade, Rare), 0)
	// An id that is not attached to any creature is a no-op, not a crash.
	g.destroyAttachedUpgrade(up)
	if _, ok := g.hostOf(up); ok {
		t.Error("an unattached upgrade should have no host")
	}
}

func TestDestroyGivesAmberToOpponent(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 3), 0) // owned by player 0
	g.State.Cards[c].Amber = 2
	if g.AmberOn(c) != 2 {
		t.Fatalf("AmberOn = %d, want 2", g.AmberOn(c))
	}
	g.DestroyEach(0, []LocalID{c})
	// The Æmber on the creature goes to its owner's opponent (player 1).
	if g.Aember(1) != 2 {
		t.Errorf("opponent aember = %d, want 2", g.Aember(1))
	}
	if g.AmberOn(c) != 0 {
		t.Errorf("amber not cleared: %d", g.AmberOn(c))
	}
}

// TestLeavePlayReleasesCapturedAember covers the non-destruction exits: a creature
// carrying captured Æmber that is bounced, archived, decked, shuffled, abducted, or
// grafted away releases that Æmber to its controller's opponent, just as destroying
// it does (Master Rulebook line 378).
func TestLeavePlayReleasesCapturedAember(t *testing.T) {
	exits := map[string]func(*Game, LocalID){
		"return to hand": func(g *Game, id LocalID) { g.putIntoHand(id) },
		"archive":        func(g *Game, id LocalID) { g.putIntoArchives(id) },
		"deck top":       func(g *Game, id LocalID) { g.putOnTopOfDeck(id) },
		"shuffle":        func(g *Game, id LocalID) { g.putIntoDeckShuffled(id) },
		"abduct":         func(g *Game, id LocalID) { g.PutIntoYourArchives(id, 1) },
		"graft":          func(g *Game, id LocalID) { g.GraftUnder(id, g.AddArtifact(NewCard("box", Mars, Artifact, Rare), 0)) },
	}
	for name, exit := range exits {
		t.Run(name, func(t *testing.T) {
			g := NewGame("A", "B", 1)
			c := g.AddToBattleline(testCreature("c", 3), 0) // owned by player 0
			g.State.Cards[c].Amber = 2
			exit(g, c)
			if g.Aember(1) != 2 {
				t.Errorf("opponent aember = %d, want 2", g.Aember(1))
			}
			if g.AmberOn(c) != 0 {
				t.Errorf("amber not cleared: %d", g.AmberOn(c))
			}
		})
	}
}

// TestNonCreatureAemberReturnsToSupply covers rule 927's second clause: a
// non-creature card leaving play with Æmber on it returns that Æmber to the common
// supply, so neither pool grows.
func TestNonCreatureAemberReturnsToSupply(t *testing.T) {
	g := NewGame("A", "B", 1)
	art := g.AddArtifact(NewCard("relic", Dis, Artifact, Rare), 0)
	g.State.Cards[art].Amber = 3
	g.putIntoHand(art)
	if g.Aember(0) != 0 || g.Aember(1) != 0 {
		t.Errorf("pools = %d/%d, want 0/0 (Æmber to the common supply)", g.Aember(0), g.Aember(1))
	}
	if g.AmberOn(art) != 0 {
		t.Errorf("amber not cleared: %d", g.AmberOn(art))
	}
}

func TestPurgesDestroyed(t *testing.T) {
	g := started(t)
	g.AddArtifact(NewCard("ritual", Dis, Artifact, Rare,
		WithConstantAbility(ConstantAbility{
			Target: Target{Kind: TargetEachCreature},
			Granted: []Ability{
				{
					Trigger: TriggerDestroyed,
					Effect:  PurgeCreature{Target: Target{Kind: TargetThisCreature}},
				},
			},
		})), 0)
	enemy := g.AddToBattleline(NewCard("v", Brobnar, Creature, Common, WithPower(3),
		WithAbility(TriggerDestroyed, GainAember{Player: Controller, Amount: 1})), 1)
	// The active player orders the ritual's granted purge before the creature's
	// printed gain; the reversing chooser picks the last-gathered ability first, and
	// the constant-granted purge is gathered after the printed gain. Purging the
	// creature stops its remaining Destroyed abilities.
	g.SetChooser(0, &countingReverseChooser{})

	g.DestroyEach(0, []LocalID{enemy})

	if got := g.Purge(1); len(got) != 1 || got[0] != enemy {
		t.Errorf("purge = %v, want [enemy]", got)
	}
	if len(g.Discard(1)) != 0 {
		t.Error("a purged creature must not be in the discard pile")
	}
	if g.Aember(1) != 0 {
		t.Error("a purged creature's remaining Destroyed abilities must not resolve")
	}
}

func TestDestroyedAbilitiesCollectEverySource(t *testing.T) {
	g := started(t)
	enemy := g.AddToBattleline(NewCard("v", Brobnar, Creature, Common, WithPower(3),
		WithAbility(TriggerDestroyed, GainAember{Player: Controller, Amount: 1}),
		WithAbility(TriggerAfterReap, GainAember{Player: Controller, Amount: 1})), 0)
	attachUpgrade(g, enemy, NewCard("upgrade", Brobnar, Upgrade, Common,
		WithStatic(StaticModifier{Granted: []Ability{
			{Trigger: TriggerDestroyed, Effect: GainAember{Player: Controller, Amount: 1}},
			{Trigger: TriggerAfterReap, Effect: GainAember{Player: Controller, Amount: 1}},
		}})))
	g.AddArtifact(NewCard("grantor", Dis, Artifact, Rare, WithConstantAbility(ConstantAbility{
		Target: Target{Kind: TargetEachCreature},
		Granted: []Ability{
			{Trigger: TriggerDestroyed, Effect: GainAember{Player: Controller, Amount: 1}},
			{Trigger: TriggerAfterReap, Effect: GainAember{Player: Controller, Amount: 1}},
		},
	})), 1)
	g.AddArtifact(NewCard("other", Dis, Artifact, Rare, WithConstantAbility(ConstantAbility{
		Target: Target{Kind: TargetEachFriendlyCreature},
		Granted: []Ability{
			{Trigger: TriggerDestroyed, Effect: GainAember{Player: Controller, Amount: 1}},
		},
	})), 1)

	if got := g.destroyedAbilities([]LocalID{enemy}); len(got) != 3 {
		t.Errorf("destroyed abilities = %d, want printed + upgrade + constant = 3", len(got))
	}
}

// TestDestroyedEachDoesNotReTriggerItself checks the destroying-window guard: a
// creature whose "Destroyed: Destroy each creature" ability re-selects itself
// (Harbinger of Doom) is already being destroyed by the enclosing batch, so its
// own Destroyed ability does not fire again forever and it is discarded once. Two
// such creatures in play prove the guard nests: destroying one wipes the whole
// board — including the other, whose identical ability adds nothing new — and the
// destruction terminates with every creature in the discard pile exactly once.
func TestDestroyedEachDoesNotReTriggerItself(t *testing.T) {
	g := started(t)
	wipe := Destroy{Target: Target{Kind: TargetEachCreature}}
	harb := g.AddToBattleline(testCreature("harb", 3, WithAbility(TriggerDestroyed, wipe)), 0)
	harb2 := g.AddToBattleline(testCreature("harb2", 3, WithAbility(TriggerDestroyed, wipe)), 0)
	bystander := g.AddToBattleline(testCreature("bystander", 3), 0)

	g.DestroyEach(0, []LocalID{harb})

	discard := g.Discard(0)
	for _, id := range []LocalID{harb, harb2, bystander} {
		if g.inPlay(id) {
			t.Errorf("creature %v still in play; the board wipe should have destroyed it", id)
		}
		n := 0
		for _, d := range discard {
			if d == id {
				n++
			}
		}
		if n != 1 {
			t.Errorf("creature %v appears %d times in discard, want exactly 1", id, n)
		}
	}
}
