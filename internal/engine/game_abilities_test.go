package engine

import "testing"

// These tests exercise USING a card — reaping, action abilities, and the checks
// that gate them — plus the machinery that fires triggered abilities, including
// those granted by upgrades. Combat lives in game_combat_test.go and destruction
// in game_destroy_test.go. Tests for specific released cards live with those
// cards in the set packages under cards/.

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

	g.BeginTurn(1)
	if got := g.TimesUsedThisTurn(user); got != 0 {
		t.Fatalf("BeginTurn reset times used = %d, want 0", got)
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
