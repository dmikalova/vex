package engine

import "testing"

// These tests exercise USING a card — reaping, action abilities, and the checks
// that gate them — plus the machinery that fires triggered abilities, including
// those granted by upgrades. Combat lives in game_combat_test.go and destruction
// in game_destroy_test.go. Tests for specific released cards live with those
// cards in the set packages under cards/.

// isFixedPredicate and fixedNarrowing recognize exactly the conditions decided by
// facts fixed for the window — the played, used, or discarded card and its board
// position, or whose turn it is — so a reaction gated on one is narrowed before
// resolution instead of joining every ordering window (Dexus's flank, Dark Æmber
// Vault's friendly Mutant creature, Pile of Skulls' enemy-during-your-turn).
func TestSubjectPredicateNarrowing(t *testing.T) {
	fixed := []Condition{
		ItIsFriendly{},
		ItIsEnemy{},
		ItIsYourTurn{},
		ItIsOfTrait{Trait: Giant},
		OnFlank{OfIt: true, Where: RightFlank},
		And{Conditions: []Condition{ItIsFriendly{}, ItIsOfTrait{Trait: Giant}}},
		And{Conditions: []Condition{ItIsEnemy{}, ItIsYourTurn{}}},
	}
	for _, c := range fixed {
		if !isFixedPredicate(c) {
			t.Errorf("%T should be a fixed predicate", c)
		}
		if _, ok := fixedNarrowing(Conditional{Cond: c, Then: Draw{Amount: 1}}); !ok {
			t.Errorf("Conditional{%T} should narrow", c)
		}
	}
	notFixed := []Condition{
		OnFlank{OfIt: false, Where: RightFlank},
		PoolAember{Player: Opponent, Is: AtLeast, Amount: 1},
		And{},
		And{
			Conditions: []Condition{
				ItIsFriendly{},
				PoolAember{Player: Opponent, Is: AtLeast, Amount: 1},
			},
		},
	}
	for _, c := range notFixed {
		if isFixedPredicate(c) {
			t.Errorf("%T should not be a fixed predicate", c)
		}
	}
	// A Conditional with an Else, over a board condition, or a non-Conditional
	// effect does not narrow.
	if _, ok := fixedNarrowing(Conditional{
		Cond: ItIsFriendly{},
		Then: Draw{Amount: 1},
		Else: Draw{Amount: 1},
	}); ok {
		t.Error("a Conditional with an Else should not narrow")
	}
	if _, ok := fixedNarrowing(Conditional{
		Cond: PoolAember{Player: Opponent, Is: AtLeast, Amount: 1},
		Then: Draw{Amount: 1},
	}); ok {
		t.Error("a board-gated Conditional should not narrow")
	}
	if _, ok := fixedNarrowing(Draw{Amount: 1}); ok {
		t.Error("a non-Conditional effect should not narrow")
	}
}

func TestReapAndActionAbility(t *testing.T) {
	g := started(t)
	reaper := testCreature(
		"reaper",
		2,
		WithAbility(TriggerAfterReap, GainAember{Player: Controller, Amount: 1}),
	)
	id := g.AddToBattleline(reaper, 0)
	if err := g.Reap(0, id); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	// +1 from reap, +1 from the reap ability.
	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2", g.Aember(0))
	}
	if !g.Exhausted(id) {
		t.Error("reaper should be exhausted")
	}

	// Action ability.
	actor := testCreature(
		"actor",
		2,
		WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 3}),
	)
	aid := g.AddToBattleline(actor, 0)
	if err := g.UseAction(0, aid); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if g.Aember(0) != 5 {
		t.Errorf("aember = %d, want 5", g.Aember(0))
	}
	// Re-using an exhausted creature fails via canUse.
	if err := g.UseAction(0, aid); err != ErrCardExhausted {
		t.Errorf("second UseAction err = %v, want ErrCardExhausted", err)
	}
	// A creature without an action ability cannot use one.
	noAction := g.AddToBattleline(testCreature("plain", 2), 0)
	if err := g.UseAction(0, noAction); err != ErrWrongType {
		t.Errorf("UseAction err = %v, want ErrWrongType", err)
	}
}

func TestTimesUsedThisTurn(t *testing.T) {
	g := started(t)
	user := g.AddToBattleline(
		testCreature(
			"user",
			6,
			WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 1}),
		),
		0,
	)
	defender := g.AddToBattleline(testCreature("defender", 1), 1)

	if err := g.Reap(0, user); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if got := g.TimesUsedThisTurn(user); got != 1 {
		t.Fatalf("after reap, times used = %d, want 1", got)
	}

	g.State.Cards[user].Exhausted = false
	if err := g.Fight(0, user, defender); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if got := g.TimesUsedThisTurn(user); got != 2 {
		t.Fatalf("after fight, times used = %d, want 2", got)
	}

	g.State.Cards[user].Exhausted = false
	if err := g.UseAction(0, user); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if got := g.TimesUsedThisTurn(user); got != 3 {
		t.Fatalf("after action, times used = %d, want 3", got)
	}

	g.StartTurn(1)
	if got := g.TimesUsedThisTurn(user); got != 0 {
		t.Fatalf("StartTurn reset times used = %d, want 0", got)
	}
}

func TestCanUseErrors(t *testing.T) {
	g := started(t)
	// Enemy creature (owner mismatch).
	enemy := g.AddToBattleline(testCreature("enemy", 2), 1)
	if err := g.Reap(0, enemy); err != ErrWrongType {
		t.Errorf("reap enemy err = %v, want ErrWrongType", err)
	}
	// Exhausted.
	own := g.AddToBattleline(testCreature("own", 2), 0)
	g.State.Cards[own].Exhausted = true
	if err := g.CanUse(0, own); err != ErrCardExhausted {
		t.Errorf("CanUse exhausted = %v, want ErrCardExhausted", err)
	}
	if err := g.Reap(0, own); err != ErrCardExhausted {
		t.Errorf("err = %v, want ErrCardExhausted", err)
	}
	g.State.Cards[own].Exhausted = false
	// A creature played this turn enters exhausted, so it cannot be used yet.
	sick := g.AddToHand(testCreature("sick", 2), 0)
	if _, err := g.PlayCreature(0, handIdxByID(g, 0, sick), false); err != nil {
		t.Fatal(err)
	}
	if !g.Exhausted(sick) {
		t.Error("a freshly played creature should be exhausted")
	}
	if err := g.Reap(0, sick); err != ErrCardExhausted {
		t.Errorf("err = %v, want ErrCardExhausted", err)
	}
	// Wrong house.
	disC := NewCard("Dis Creature", Dis, Creature, Common, WithPower(2))
	disID := g.AddToBattleline(disC, 0)
	if err := g.Reap(0, disID); err != ErrWrongHouse {
		t.Errorf("err = %v, want ErrWrongHouse", err)
	}
	// An artifact is in play but is not a creature, so it cannot reap or fight.
	relic := g.AddArtifact(NewCard("Relic", Brobnar, Artifact, Common), 0)
	if err := g.Reap(0, relic); err != ErrWrongType {
		t.Errorf("reap artifact err = %v, want ErrWrongType", err)
	}
	// Not active player.
	if err := g.Reap(1, own); err != ErrNotActivePlayer {
		t.Errorf("err = %v, want ErrNotActivePlayer", err)
	}
	// Game over.
	g.State.Winner = 0
	if err := g.Reap(0, own); err != ErrGameOver {
		t.Errorf("err = %v, want ErrGameOver", err)
	}
}

// TestFightGrantForgivesWrongHouseForFightOnly verifies that a fight grant
// (Brothers in Battle) lets an off-house creature be used to fight — so a UI that
// asks CanUseTo before offering Fight offers it — while leaving reap barred.
func TestFightGrantForgivesWrongHouseForFightOnly(t *testing.T) {
	g := started(t) // Brobnar active
	off := g.AddToBattleline(NewCard("off", Sanctum, Creature, Common, WithPower(3)), 0)
	g.AddToBattleline(testCreature("foe", 3), 1)

	// Without the grant, fighting out of house is barred.
	if err := g.CanUseTo(0, off, FightUse); err != ErrWrongHouse {
		t.Errorf("ungranted fight = %v, want ErrWrongHouse", err)
	}

	g.State.MayFightHouse[0] = Sanctum
	if err := g.CanUseTo(0, off, FightUse); err != nil {
		t.Errorf("granted fight = %v, want nil", err)
	}
	// The grant is fight-only: reap stays barred out of house.
	if err := g.CanUseTo(0, off, ReapUse); err != ErrWrongHouse {
		t.Errorf("granted reap = %v, want ErrWrongHouse", err)
	}
}

func TestCanUseArtifact(t *testing.T) {
	g := started(t)
	relic := g.AddArtifact(NewCard("Relic", Brobnar, Artifact, Common), 0)
	if err := g.CanUseArtifact(0, relic); err != ErrCannotUse {
		t.Errorf("artifact without an action ability = %v, want ErrCannotUse", err)
	}

	actor := g.AddArtifact(
		NewCard(
			"Actor",
			Brobnar,
			Artifact,
			Common,
			WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 1}),
		),
		0,
	)
	if err := g.CanUseArtifact(0, actor); err != nil {
		t.Errorf("actor in the active house = %v, want nil", err)
	}

	// Wrong house, no Versatile: blocked.
	g.State.ActiveHouse = Dis
	if err := g.CanUseArtifact(0, actor); err != ErrWrongHouse {
		t.Errorf("actor out of house = %v, want ErrWrongHouse", err)
	}

	// Versatile relaxes the house check for an artifact same as a creature.
	versatile := g.AddArtifact(
		NewCard(
			"Versatile Actor", Brobnar, Artifact, Common,
			WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 1}),
			WithKeywords(Versatile),
		),
		0,
	)
	if err := g.CanUseArtifact(0, versatile); err != nil {
		t.Errorf("Versatile actor out of house = %v, want nil", err)
	}
	g.State.ActiveHouse = Brobnar

	// A creature is not an artifact.
	creature := g.AddToBattleline(testCreature("creature", 2), 0)
	if err := g.CanUseArtifact(0, creature); err != ErrWrongType {
		t.Errorf("CanUseArtifact on a creature = %v, want ErrWrongType", err)
	}

	// usable()'s own checks (e.g. wrong player) surface unchanged.
	if err := g.CanUseArtifact(1, actor); err != ErrNotActivePlayer {
		t.Errorf("CanUseArtifact wrong player = %v, want ErrNotActivePlayer", err)
	}
}

func TestAfterCardPlayedTrigger(t *testing.T) {
	g := started(t)
	g.AddToBattleline(
		testCreature(
			"watcher",
			3,
			WithAbility(TriggerAfterCardPlayed, GainAember{Player: Controller, Amount: 1}),
		),
		0,
	)
	g.AddToHand(NewCard("relic", Brobnar, Artifact, Common), 0)
	if _, err := g.PlayArtifact(0, 0); err != nil {
		t.Fatalf("PlayArtifact: %v", err)
	}
	if g.Aember(0) != 1 {
		t.Errorf("aember = %d, want 1 (watcher triggered on card play)", g.Aember(0))
	}
}

// TestSubjectNarrowedReactionEntersWindowOnlyWhenMatched checks that a reaction
// narrowed to a card shape — "after you play an artifact" (a Conditional{ItIs}
// over its whole effect) — joins the after-play window only when the played card
// matches, so it never forces a spurious ordering prompt on a play it does not
// narrow to (Harmonia not ordering with a played Tactic). A reaction gated on a
// board "if" instead (Overwhelmed) is not a subject narrowing, so it always fires.
func TestSubjectNarrowedReactionEntersWindowOnlyWhenMatched(t *testing.T) {
	g := started(t)
	g.AddToBattleline(testCreature("narrowed", 3, WithAbility(
		TriggerAfterCardPlayed,
		Conditional{Cond: ItIs{Type: Artifact}, Then: GainAember{Player: Controller, Amount: 1}},
	)), 0)
	g.AddToBattleline(testCreature("board", 3, WithAbility(
		TriggerAfterCardPlayed,
		Conditional{Cond: Overwhelmed{}, Then: GainAember{Player: Controller, Amount: 1}},
	)), 0)

	// A Tactic does not match the artifact narrowing, so only the board-gated
	// reaction joins the window — no spurious prompt from the narrowed one.
	tactic := g.AddToHand(NewCard("scheme", Brobnar, Tactic, Common), 0)
	if w := g.afterPlayReactions(0, tactic); len(w) != 1 {
		t.Fatalf("Tactic play window = %d entries, want 1 (narrowed watcher excluded)", len(w))
	}

	// An artifact matches, so both reactions join the window.
	relic := g.AddArtifact(NewCard("relic", Brobnar, Artifact, Common), 0)
	if w := g.afterPlayReactions(0, relic); len(w) != 2 {
		t.Fatalf("artifact play window = %d entries, want 2 (narrowed watcher matches)", len(w))
	}
}

func TestUpgradeUseConditionGatesHost(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	host := g.AddToBattleline(NewCard("host", Brobnar, Creature, Common, WithPower(3)), 0)
	// An attached upgrade bars the host's use until its condition is met.
	attachUpgrade(g, host, NewCard(
		"bind",
		Untamed,
		Upgrade,
		Common,
		WithRestrictions(
			Restrictions{UseCondition: CardsDiscarded{Player: Controller, Amount: 1}},
		),
	))
	if err := g.usable(0, host); err != ErrCannotUse {
		t.Fatalf("host with unmet upgrade use-condition = %v, want ErrCannotUse", err)
	}
	c := g.AddToHand(NewCard("spare", Brobnar, Tactic, Common), 0)
	g.DiscardCardFromHand(0, c)
	if err := g.usable(0, host); err != nil {
		t.Fatalf("host after discard = %v, want nil", err)
	}
}

func TestGrantedKeywordsFromUpgrades(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	wall := g.AddToBattleline(testCreature("wall", 10), 1)
	// An upgrade granting elusive + skirmish. hasKeyword checks the printed
	// keywords first, then each granted keyword.
	attachUpgrade(g, host, NewCard("cloak", Shadows, Upgrade, Common,
		WithStatic(StaticModifier{Keywords: []Keyword{Elusive, Skirmish}})))
	if !g.hasKeyword(host, Skirmish) {
		t.Fatal("granted skirmish should be reported by hasKeyword")
	}
	if g.hasKeyword(host, Poison) {
		t.Fatal("hasKeyword should not report a keyword neither printed nor granted")
	}
	// Skirmish spares the attacker from the 10-power wall's retaliation.
	g.fight(host, wall)
	if !g.inPlay(host) {
		t.Error("granted skirmish should spare the attacker")
	}
	if g.Damage(host) != 0 {
		t.Errorf("attacker took %d damage, want 0 (granted skirmish)", g.Damage(host))
	}
}

func TestStunBehavior(t *testing.T) {
	g := started(t)

	// Reap: a stunned creature recovers instead of reaping (no Æmber gained).
	r := g.AddToBattleline(testCreature("stunreap", 2), 0)
	g.State.Cards[r].Stunned = true
	before := g.State.Aember[0]
	if err := g.Reap(0, r); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	if g.State.Cards[r].Stunned {
		t.Error("reaping a stunned creature should clear its stun")
	}
	if g.State.Aember[0] != before {
		t.Error("reaping a stunned creature should not gain Æmber")
	}
	if !g.State.Cards[r].Exhausted {
		t.Error("recovering from stun should exhaust the creature")
	}

	// Fight: a stunned attacker recovers instead of fighting (defender unharmed).
	a := g.AddToBattleline(testCreature("stunfight", 5), 0)
	d := g.AddToBattleline(testCreature("enemy", 4), 1)
	g.State.Cards[a].Stunned = true
	if err := g.Fight(0, a, d); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.State.Cards[a].Stunned {
		t.Error("fighting with a stunned creature should clear its stun")
	}
	if g.Damage(d) != 0 {
		t.Errorf("stunned attacker should deal no damage; defender took %d", g.Damage(d))
	}

	// Action: a stunned creature recovers instead of using its action ability.
	act := g.AddToBattleline(
		testCreature(
			"stunact",
			3,
			WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 5}),
		),
		0,
	)
	g.State.Cards[act].Stunned = true
	before = g.State.Aember[0]
	if err := g.UseAction(0, act); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if g.State.Cards[act].Stunned {
		t.Error("using a stunned creature's action should clear its stun")
	}
	if g.State.Aember[0] != before {
		t.Error("a stunned action ability should not resolve")
	}
}

// TestMustFightWhenEnraged checks that an enraged creature must fight when used if
// it can: with an enemy to fight it may not reap or act, but with nothing to fight
// enrage cannot force it, so it may reap.
func TestMustFightWhenEnraged(t *testing.T) {
	// With an enemy present, an enraged creature must fight: reap and action are
	// barred, but fighting is still allowed.
	g := started(t)
	g.AddToBattleline(testCreature("enemy", 3), 1)
	withFoe := g.AddToBattleline(testCreature("enraged-foe", 3), 0)
	g.State.Cards[withFoe].Enraged = true
	if g.canUseTo(0, withFoe, ReapUse) != ErrCannotUse {
		t.Error("enraged creature with an enemy should not be able to reap")
	}
	if g.canUseTo(0, withFoe, ActionUse) != ErrCannotUse {
		t.Error("enraged creature with an enemy should not be able to use an Action")
	}
	if err := g.canUseTo(0, withFoe, FightUse); err != nil {
		t.Errorf("enraged creature should still be able to fight: %v", err)
	}

	// With no enemy to fight, enrage cannot force a fight, so the creature may reap.
	g2 := started(t)
	alone := g2.AddToBattleline(testCreature("enraged-alone", 3), 0)
	g2.State.Cards[alone].Enraged = true
	if err := g2.canUseTo(0, alone, ReapUse); err != nil {
		t.Errorf("enraged creature with nothing to fight should be able to reap: %v", err)
	}
}

// TestUnstun checks that Unstun spends a stunned, otherwise-usable creature's
// use shaking off the stun instead of reaping/fighting/acting, under the same
// checks Reap/Fight/an action would apply — including the active-house one,
// since a stunned creature only gets this choice when it is otherwise usable.
func TestUnstun(t *testing.T) {
	g := started(t)
	r := g.AddToBattleline(testCreature("stunreap", 2), 0)
	g.State.Cards[r].Stunned = true
	if err := g.Unstun(0, r); err != nil {
		t.Fatalf("Unstun: %v", err)
	}
	if g.State.Cards[r].Stunned {
		t.Error("Unstun should clear the stun")
	}
	if !g.State.Cards[r].Exhausted {
		t.Error("Unstun should exhaust the creature")
	}
	// Out of house, a stunned creature has no use to spend at all — Unstun is
	// blocked exactly like Reap, not a house-independent escape hatch.
	outsider := g.AddToBattleline(NewCard("outsider", Dis, Creature, Common, WithPower(2)), 0)
	g.State.Cards[outsider].Stunned = true
	if err := g.Unstun(0, outsider); err != ErrWrongHouse {
		t.Errorf("Unstun out of house = %v, want ErrWrongHouse", err)
	}
	// A creature that is not stunned has nothing for Unstun to do.
	fine := g.AddToBattleline(testCreature("fine", 2), 0)
	if err := g.Unstun(0, fine); err != ErrCannotUse {
		t.Errorf("Unstun an unstunned creature = %v, want ErrCannotUse", err)
	}
	// The checks Unstun shares with every other use still apply: an
	// already-exhausted creature has nothing left to spend.
	spent := g.AddToBattleline(testCreature("spent", 2), 0)
	g.State.Cards[spent].Stunned = true
	g.State.Cards[spent].Exhausted = true
	if err := g.Unstun(0, spent); err != ErrCardExhausted {
		t.Errorf("Unstun an exhausted creature = %v, want ErrCardExhausted", err)
	}
}

func TestGrantedAbilitiesFireFromUpgrades(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 5), 0)
	// An upgrade granting its host "Reap: Gain 1 Æmber".
	attachUpgrade(g, host, NewCard("charm", Shadows, Upgrade, Common,
		WithStatic(StaticModifier{Granted: []Ability{
			{Trigger: TriggerAfterReap, Effect: GainAember{Player: Controller, Amount: 1}},
		}})))

	// A non-matching trigger fires nothing granted.
	g.triggerAbilities(host, TriggerAfterFight, 0, false)
	if g.Aember(0) != 0 {
		t.Fatalf("granted Reap ability fired on a Fight trigger; aember = %d", g.Aember(0))
	}
	// The matching trigger fires the granted ability with the host as source.
	g.triggerAbilities(host, TriggerAfterReap, 0, false)
	if g.Aember(0) != 1 {
		t.Errorf("granted Reap ability aember = %d, want 1", g.Aember(0))
	}
}

func TestGrantedActionAbilitiesCanBeUsed(t *testing.T) {
	g := started(t)
	upgraded := g.AddToBattleline(testCreature("upgraded", 5), 0)
	attachUpgrade(g, upgraded, NewCard("controls", Logos, Upgrade, Rare,
		WithStatic(StaticModifier{Granted: []Ability{
			{Trigger: TriggerAction, Effect: GainAember{Player: Controller, Amount: 2}},
		}})))

	if err := g.UseAction(0, upgraded); err != nil {
		t.Fatalf("UseAction with upgrade-granted action: %v", err)
	}
	if g.Aember(0) != 2 {
		t.Fatalf("aember after upgrade-granted action = %d, want 2", g.Aember(0))
	}

	constantGranted := g.AddToBattleline(testCreature("constant-granted", 5), 0)
	g.AddArtifact(NewCard("ritual", Logos, Artifact, Rare,
		WithConstantAbility(ConstantAbility{
			Target: Target{Kind: TargetEachCreature},
			Granted: []Ability{
				{Trigger: TriggerAction, Effect: GainAember{Player: Controller, Amount: 3}},
			},
		})), 0)

	if err := g.UseAction(0, constantGranted); err != nil {
		t.Fatalf("UseAction with constant-granted action: %v", err)
	}
	if g.Aember(0) != 5 {
		t.Fatalf("aember after constant-granted action = %d, want 5", g.Aember(0))
	}
}

func TestGrantedAbilitiesFireFromConstants(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 5), 0)
	g.AddArtifact(NewCard("ritual", Dis, Artifact, Rare, WithConstantAbility(ConstantAbility{
		Target: Target{Kind: TargetEachCreature},
		Granted: []Ability{
			{Trigger: TriggerAfterReap, Effect: GainAember{Player: Controller, Amount: 1}},
		},
	})), 1)

	// The matching trigger fires the ability granted by an opponent's constant.
	g.triggerAbilities(host, TriggerAfterReap, 0, false)
	if g.Aember(0) != 1 {
		t.Errorf("constant-granted Reap ability aember = %d, want 1", g.Aember(0))
	}
}

// A constant ability with DisableTriggers stops the listed triggers from firing
// board-wide while its source stays in play: Purifier of Souls disables every
// Destroyed ability, so a creature destroyed alongside it gains its controller
// nothing.
func TestDisableTriggersStopsDestroyedAbilities(t *testing.T) {
	g := started(t)
	// A buffer with a plain constant ability (no DisableTriggers) is in play too,
	// so the disable scan skips a source that disables nothing.
	g.AddArtifact(NewCard("buffer", Dis, Artifact, Rare,
		WithConstantAbility(ConstantAbility{PowerBonus: 1})), 0)
	// The purifier lists a non-matching trigger before Destroyed, so the scan walks
	// past a trigger it does not disable before finding the one it does.
	g.AddArtifact(NewCard("purifier", Sanctum, Artifact, Rare,
		WithConstantAbility(ConstantAbility{
			DisableTriggers: []Trigger{TriggerAction, TriggerDestroyed},
		})), 0)
	victim := g.AddToBattleline(NewCard("v", Brobnar, Creature, Common, WithPower(3),
		WithAbility(TriggerDestroyed, GainAember{Player: Controller, Amount: 1})), 0)

	g.DestroyEach(0, []LocalID{victim})

	if g.Aember(0) != 0 {
		t.Errorf("aember = %d, want 0 (Destroyed ability must not fire)", g.Aember(0))
	}
}

func TestDisableTriggersText(t *testing.T) {
	def := NewCard("Purifier", Sanctum, Creature, Rare, WithPower(5), WithArmor(2),
		WithConstantAbility(ConstantAbility{
			DisableTriggers: []Trigger{TriggerDestroyed},
		}))
	if got := constantText(&def); got != "Destroyed effects cannot trigger." {
		t.Errorf("constantText = %q, want %q", got, "Destroyed effects cannot trigger.")
	}
}

// exBoardWatcher is an artifact that draws a card whenever any creature is played,
// wherever it sits and whoever plays it (The Big One watches the whole board).
func exBoardWatcher() CardDefinition {
	return NewCard(
		"Board Watcher",
		Brobnar,
		Artifact,
		Rare,
		WithAbility(TriggerAfterCreaturePlayed, Draw{Amount: 1}),
	)
}

// exCreatureWatcher is a creature carrying the same global trigger, used to check
// that the trigger never fires on the very creature that was played.
func exCreatureWatcher() CardDefinition {
	return NewCard(
		"Creature Watcher",
		Brobnar,
		Creature,
		Rare,
		WithPower(3),
		WithAbility(TriggerAfterCreaturePlayed, Draw{Amount: 1}),
	)
}

// TestAfterCreaturePlayed covers the global "after a creature is played" trigger:
// it fires for a creature played anywhere on the board, for either player's play,
// but not on the played creature itself.
func TestAfterCreaturePlayed(t *testing.T) {
	t.Run("fires when the controller plays a creature away from the watcher", func(t *testing.T) {
		g := started(t)
		g.State.ActiveHouse = Brobnar
		g.AddArtifact(exBoardWatcher(), 0)
		g.AddToBattleline(testCreature("buffer", 2), 0)
		g.AddToDeck(testCreature("top", 2), 0)
		g.AddToHand(NewCard("newbie", Brobnar, Creature, Common, WithPower(3)), 0)
		before := g.State.Hand[0].Count

		if _, err := g.PlayCreature(0, handIdx(g, 0, "newbie"), false); err != nil {
			t.Fatalf("PlayCreature: %v", err)
		}

		// -1 for the creature played out of hand, +1 for the card the watcher drew.
		if got := g.State.Hand[0].Count; got != before {
			t.Errorf("hand = %d, want %d (played one, drew one)", got, before)
		}
	})

	t.Run("fires when the opponent plays a creature", func(t *testing.T) {
		g := started(t)
		g.State.ActiveHouse = Brobnar
		g.AddArtifact(exBoardWatcher(), 0)
		g.AddToDeck(testCreature("top", 2), 0)
		before := g.State.Deck[0].Count

		enemy := g.AddToBattleline(testCreature("enemy", 3), 1)
		g.emitCreaturePlayed(enemy)

		if got := g.State.Deck[0].Count; got != before-1 {
			t.Errorf("deck = %d, want %d (watcher drew for the enemy play)", got, before-1)
		}
	})

	t.Run("does not fire on the played creature itself", func(t *testing.T) {
		g := started(t)
		g.State.ActiveHouse = Brobnar
		g.AddToDeck(testCreature("top", 2), 0)
		before := g.State.Deck[0].Count

		self := g.AddToBattleline(exCreatureWatcher(), 0)
		g.emitCreaturePlayed(self)

		if got := g.State.Deck[0].Count; got != before {
			t.Errorf("deck = %d, want %d (a card must not fire on itself)", got, before)
		}
	})
}

// exWatcher is a creature that draws a card whenever another creature is played
// into a battleline position adjacent to it (Fila the Researcher).
func exWatcher() CardDefinition {
	return NewCard(
		"Watcher",
		Logos,
		Creature,
		Uncommon,
		WithPower(1),
		WithAbility(TriggerAfterCreaturePlayedAdjacent, Draw{Amount: 1}),
	)
}

// TestAfterCreaturePlayedAdjacent covers the trigger Fila the Researcher uses: it
// fires for a creature played next to the watcher, but not for one played away
// from it.
func TestAfterCreaturePlayedAdjacent(t *testing.T) {
	t.Run("fires when a creature is played next to the watcher", func(t *testing.T) {
		g := started(t)
		g.State.ActiveHouse = Logos
		g.AddToBattleline(exWatcher(), 0)
		g.AddToDeck(testCreature("top", 2), 0)
		g.AddToHand(NewCard("newbie", Logos, Creature, Common, WithPower(3)), 0)
		before := g.State.Hand[0].Count

		if _, err := g.PlayCreature(0, handIdx(g, 0, "newbie"), false); err != nil {
			t.Fatalf("PlayCreature: %v", err)
		}

		// -1 for the creature played out of hand, +1 for the card the watcher drew.
		if got := g.State.Hand[0].Count; got != before {
			t.Errorf("hand = %d, want %d (played one, drew one)", got, before)
		}
	})

	t.Run("does not fire for a creature played away from the watcher", func(t *testing.T) {
		g := started(t)
		g.State.ActiveHouse = Logos
		g.AddToBattleline(exWatcher(), 0)
		g.AddToBattleline(testCreature("buffer", 2), 0)
		g.AddToDeck(testCreature("top", 2), 0)
		g.AddToHand(NewCard("newbie", Logos, Creature, Common, WithPower(3)), 0)
		before := g.State.Hand[0].Count

		// flankLeft: played on the far side, next to the buffer, not the watcher.
		if _, err := g.PlayCreature(0, handIdx(g, 0, "newbie"), false); err != nil {
			t.Fatalf("PlayCreature: %v", err)
		}

		// Only the played creature left the hand; the watcher never drew.
		if got := g.State.Hand[0].Count; got != before-1 {
			t.Errorf("hand = %d, want %d (played one, drew none)", got, before-1)
		}
	})
}

func TestEmitAfterEnemyDestroyed(t *testing.T) {
	pileDef := NewCard("Pile", Brobnar, Artifact, Rare,
		WithAbility(TriggerAfterCreatureDestroyed, Conditional{
			Cond: And{Conditions: []Condition{ItIsEnemy{}, ItIsYourTurn{}}},
			Then: CaptureAember{
				Amount: 1,
				Target: Target{Kind: TargetChosenFriendlyCreature},
				Source: Opponent,
			},
		}))

	t.Run("captures when an enemy creature is destroyed on your turn", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddArtifact(pileDef, 0)
		friendly := g.AddToBattleline(testCreature("f", 5), 0)
		enemy := g.AddToBattleline(testCreature("e", 2), 1)
		g.State.Aember[1] = 3

		g.destroyEach(0, []LocalID{enemy})

		if g.AmberOn(friendly) != 1 {
			t.Errorf("friendly Æmber = %d, want 1", g.AmberOn(friendly))
		}
		if g.Aember(1) != 2 {
			t.Errorf("opponent pool = %d, want 2", g.Aember(1))
		}
	})

	t.Run("does not fire when a friendly creature is destroyed", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddArtifact(pileDef, 0)
		friendly := g.AddToBattleline(testCreature("f", 5), 0)
		g.State.Aember[1] = 3

		g.destroyEach(0, []LocalID{friendly})

		if g.Aember(1) != 3 {
			t.Errorf("opponent pool = %d, want 3 (a friendly death does not trigger)", g.Aember(1))
		}
	})

	t.Run("does not fire on the opponent's turn", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 1
		g.AddArtifact(pileDef, 0)
		g.AddToBattleline(testCreature("f", 5), 0)
		mine := g.AddToBattleline(testCreature("m", 2), 0)
		g.State.Aember[1] = 3

		g.destroyEach(1, []LocalID{mine})

		if g.Aember(1) != 3 {
			t.Errorf("opponent pool = %d, want 3 (it is not your turn)", g.Aember(1))
		}
	})
}

func TestEmitAfterFriendlyDestroyed(t *testing.T) {
	watcherDef := NewCard("Watcher", Brobnar, Artifact, Rare,
		WithAbility(TriggerAfterCreatureDestroyed, Conditional{
			Cond: ItIsFriendly{},
			Then: GainAember{
				Player: Controller,
				Amount: 1,
			},
		}))

	t.Run("fires when a creature under the same control is destroyed", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddArtifact(watcherDef, 0)
		friendly := g.AddToBattleline(testCreature("f", 5), 0)

		g.destroyEach(0, []LocalID{friendly})

		if g.Aember(0) != 1 {
			t.Errorf("pool = %d, want 1", g.Aember(0))
		}
	})

	t.Run("does not fire when an enemy creature is destroyed", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddArtifact(watcherDef, 0)
		enemy := g.AddToBattleline(testCreature("e", 2), 1)

		g.destroyEach(0, []LocalID{enemy})

		if g.Aember(0) != 0 {
			t.Errorf("pool = %d, want 0 (an enemy death does not trigger)", g.Aember(0))
		}
	})
}

func TestEmitCreatureReaped(t *testing.T) {
	// Orb-style: after any creature reaps, stun it.
	orbDef := NewCard(
		"Orb",
		Dis,
		Artifact,
		Rare,
		WithAbility(
			TriggerAfterCreatureReaps,
			Stun{Target: Target{Kind: TargetTriggeringCreature}},
		),
	)
	// Pip-style: after an enemy creature reaps, stun it.
	pipDef := NewCard(
		"Pip",
		Logos,
		Creature,
		Common,
		WithAbility(
			TriggerAfterCreatureReaps,
			Conditional{
				Cond: ItIsEnemy{},
				Then: Stun{Target: Target{Kind: TargetTriggeringCreature}},
			},
		),
	)

	t.Run("enemy Orb stuns the reaper", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddArtifact(orbDef, 1)
		reaper := g.AddToBattleline(testCreature("r", 3), 0)

		g.reapWith(reaper)

		if !g.State.Cards[reaper].Stunned {
			t.Error("reaper should be stunned by the enemy Orb")
		}
	})

	t.Run("friendly Orb stuns your own reaper too", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddArtifact(orbDef, 0)
		reaper := g.AddToBattleline(testCreature("r", 3), 0)

		g.reapWith(reaper)

		if !g.State.Cards[reaper].Stunned {
			t.Error("friendly Orb should stun your own reaper (any creature)")
		}
	})

	t.Run("enemy-reap reaction fires for the reaper's opponent", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddToBattleline(pipDef, 1)
		reaper := g.AddToBattleline(testCreature("r", 3), 0)

		g.reapWith(reaper)

		if !g.State.Cards[reaper].Stunned {
			t.Error("the enemy reaper should be stunned by Pip")
		}
	})

	t.Run("enemy-reap reaction does not fire on the controller's own reap", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddToBattleline(pipDef, 0)
		reaper := g.AddToBattleline(testCreature("r", 3), 0)

		g.reapWith(reaper)

		if g.State.Cards[reaper].Stunned {
			t.Error("your own reap should not fire your Pip's enemy-reap reaction")
		}
	})
}

// TestAfterUpgradeEnters covers the global "after an upgrade enters play"
// trigger: it fires on an in-play watcher when an upgrade is played onto any
// creature (Armory Officer Nel draws a card).
func TestAfterUpgradeEnters(t *testing.T) {
	watcher := func() CardDefinition {
		return NewCard("Nel", StarAlliance, Creature, Common, WithPower(4),
			WithAbility(TriggerAfterUpgradeEnters, Draw{Amount: 1}))
	}

	t.Run("fires when an upgrade is played", func(t *testing.T) {
		g := started(t)
		g.State.ActiveHouse = Brobnar
		g.AddToBattleline(watcher(), 0)
		g.AddToBattleline(testCreature("host", 3), 0)
		g.AddToDeck(testCreature("top", 2), 0)
		g.AddToHand(exBruteStrength(), 0)
		before := g.State.Deck[0].Count

		if _, err := g.PlayUpgrade(0, handIdx(g, 0, "Brute Strength")); err != nil {
			t.Fatalf("PlayUpgrade: %v", err)
		}

		if got := g.State.Deck[0].Count; got != before-1 {
			t.Errorf("deck = %d, want %d (watcher drew when the upgrade entered)",
				got, before-1)
		}
	})

	t.Run("fires for the opponent's upgrade too", func(t *testing.T) {
		g := started(t)
		g.AddToBattleline(watcher(), 0)
		g.AddToDeck(testCreature("top", 2), 0)
		before := g.State.Deck[0].Count

		up := g.Register(NewCard("chip", StarAlliance, Upgrade, Common), 1)
		g.emitUpgradeEntered(up)

		if got := g.State.Deck[0].Count; got != before-1 {
			t.Errorf("deck = %d, want %d (watcher drew for the enemy upgrade)",
				got, before-1)
		}
	})
}

// TestAfterDestroyedReactionsAreCreatureOnly pins that an "after a creature
// leaves play" reaction watches creatures only. An artifact and an upgrade both
// leave play the same way a creature does — and under the invalid-state sweep an
// orphaned upgrade leaves play in its own right — but neither is a creature, so
// neither wakes Pile of Skulls or Loot the Bodies.
func TestAfterDestroyedReactionsAreCreatureOnly(t *testing.T) {
	watcher := NewCard("Pile", Brobnar, Artifact, Rare,
		WithAbility(TriggerAfterCreatureDestroyed, Draw{Amount: 1}))

	g := NewGame("A", "B", 1)
	g.State.ActivePlayer = 0
	g.AddArtifact(watcher, 0)

	creature := g.AddToBattleline(testCreature("e", 2), 1)
	artifact := g.AddArtifact(NewCard("Relic", Brobnar, Artifact, Common), 1)
	upgrade := g.Register(NewCard("Boon", Brobnar, Upgrade, Common), 1)
	g.AttachUpgrade(creature, upgrade)

	for _, tc := range []struct {
		name string
		id   LocalID
		want int
	}{
		{"creature", creature, 1},
		{"artifact", artifact, 0},
		{"upgrade", upgrade, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := len(g.afterDestroyedReactions([]LocalID{tc.id})); got != tc.want {
				t.Errorf("reactions to a destroyed %s = %d, want %d", tc.name, got, tc.want)
			}
		})
	}
}
