package engine

import "testing"

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

func TestDestroyTogetherResolvesBeforeDiscard(t *testing.T) {
	g := started(t)
	// A's "Destroyed:" ability records how many enemy creatures are in play when
	// it fires. KeyForge tags both for destruction and resolves the Destroyed
	// abilities before moving anything to the discard, so B is still present.
	var enemiesWhenADied int
	a := g.AddToBattleline(testCreature("a", 3, WithAbility(TriggerDestroyed, recordEnemyCount{got: &enemiesWhenADied})), 0)
	b := g.AddToBattleline(testCreature("b", 3), 1)

	g.DestroyEach(0, []LocalID{a, b})

	if enemiesWhenADied != 1 {
		t.Errorf("A's Destroyed ability saw %d enemy creatures; want 1 (still in play until discard)", enemiesWhenADied)
	}
	if g.inPlay(a) || g.inPlay(b) {
		t.Error("both creatures should be in the discard after the event")
	}
}

func TestDestroyedRelocationSkipsDiscard(t *testing.T) {
	g := started(t)
	// A creature whose "Destroyed:" ability returns it to the top of its deck
	// leaves play during the event, so it is not also moved to the discard.
	c := g.AddToBattleline(testCreature("wanderer", 3, WithAbility(TriggerDestroyed, PutFromPlay{Target: Target{Kind: TargetThisCreature}, Destination: ToTopOfDeck})), 0)
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
	shield := NewCard("Shield", Sanctum, Upgrade, Rare,
		WithStatic(StaticModifier{Replaces: Replace{When: EventCreatureDestroyed, With: Sequence{Effects: []Effect{
			Heal{Fully: true, Target: Target{Kind: TargetTriggeringCreature}},
			Destroy{Target: Target{Kind: TargetThisCreature}},
		}}}}))
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
	if g.State.Cards[host].UpgradeCount != 0 {
		t.Errorf("host upgrade count = %d, want the shield consumed", g.State.Cards[host].UpgradeCount)
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

func TestPurgesDestroyed(t *testing.T) {
	g := started(t)
	g.AddArtifact(NewCard("ritual", Dis, Artifact, Rare, WithConstantAbility(ConstantAbility{
		Target:  Target{Kind: TargetEachCreature},
		Granted: []Ability{{Trigger: TriggerDestroyed, Effect: PurgeCreature{Target: Target{Kind: TargetThisCreature}}}},
	})), 0)
	victim := g.AddToBattleline(NewCard("v", Brobnar, Creature, Common, WithPower(3),
		WithAbility(TriggerDestroyed, GainAember{Player: Controller, Amount: 1})), 1)
	// The active player orders the ritual's granted ability before the victim's
	// printed one. Purging the victim stops its remaining Destroyed abilities.
	g.SetChooser(0, optionPicker{idx: 1})

	g.DestroyEach(0, []LocalID{victim})

	if got := g.Purge(1); len(got) != 1 || got[0] != victim {
		t.Errorf("purge = %v, want [victim]", got)
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
	victim := g.AddToBattleline(NewCard("v", Brobnar, Creature, Common, WithPower(3),
		WithAbility(TriggerDestroyed, GainAember{Player: Controller, Amount: 1}),
		WithAbility(TriggerAfterReap, GainAember{Player: Controller, Amount: 1})), 0)
	attachUpgrade(g, victim, NewCard("upgrade", Brobnar, Upgrade, Common,
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
		Target:  Target{Kind: TargetEachFriendlyCreature},
		Granted: []Ability{{Trigger: TriggerDestroyed, Effect: GainAember{Player: Controller, Amount: 1}}},
	})), 1)

	if got := g.destroyedAbilities([]LocalID{victim}); len(got) != 3 {
		t.Errorf("destroyed abilities = %d, want printed + upgrade + constant = 3", len(got))
	}
}
