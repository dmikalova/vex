package engine

import "testing"

// These tests exercise combat: fighting, dealing damage, armor absorption, the
// Skirmish/Assault/Hazardous keywords, and the fight-timed triggers. The generic
// example blueprints come from helpers_test.go.

func TestFight(t *testing.T) {
	g := started(t)
	attacker := g.AddToBattleline(testCreature("att", 5), 0)
	defender := g.AddToBattleline(testCreature("def", 4), 1)
	if err := g.Fight(0, attacker, defender); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if len(g.Battleline(1)) != 0 {
		t.Error("defender should be destroyed")
	}
	if g.Damage(attacker) != 4 {
		t.Errorf("attacker damage = %d, want 4", g.Damage(attacker))
	}

	// Invalid target: a friendly creature.
	friend := g.AddToBattleline(testCreature("friend", 2), 0)
	att2 := g.AddToBattleline(testCreature("att2", 2), 0)
	if err := g.Fight(0, att2, friend); err != ErrNoTarget {
		t.Errorf("err = %v, want ErrNoTarget", err)
	}
	// Invalid attacker: exhausted (fails via canUse before target checks).
	g.State.Cards[att2].Exhausted = true
	if err := g.Fight(0, att2, friend); err != ErrCardExhausted {
		t.Errorf("err = %v, want ErrCardExhausted", err)
	}
}

func TestGrantFightForHouse(t *testing.T) {
	g := started(t) // Brobnar is the active house
	att := g.AddToBattleline(NewCard("wild", Untamed, Creature, Common, WithPower(5)), 0)
	def := g.AddToBattleline(testCreature("def", 3), 1)

	// Without a grant, an out-of-house creature cannot fight.
	if err := g.Fight(0, att, def); err != ErrWrongHouse {
		t.Fatalf("without grant: err = %v, want ErrWrongHouse", err)
	}

	// The grant lets creatures of that house fight this turn.
	g.GrantMayPlayOrUse(0, HouseSelector{Match: namedHouse(Untamed)}, GrantFight, 0, 0)
	if err := g.Fight(0, att, def); err != nil {
		t.Fatalf("with grant: %v", err)
	}

	// A creature of another house is still barred.
	att2 := g.AddToBattleline(NewCard("holy", Sanctum, Creature, Common, WithPower(5)), 0)
	def2 := g.AddToBattleline(testCreature("def2", 3), 1)
	if err := g.Fight(0, att2, def2); err != ErrWrongHouse {
		t.Errorf("other-house attacker: err = %v, want ErrWrongHouse", err)
	}

	// The ready phase clears the grant.
	g.EndPlayPhase(0)
	def3 := g.AddToBattleline(testCreature("def3", 3), 1)
	if err := g.Fight(0, att, def3); err != ErrWrongHouse {
		t.Errorf("after EndPlayPhase: err = %v, want ErrWrongHouse (grant cleared)", err)
	}
}

func TestFightSimultaneousDestruction(t *testing.T) {
	g := started(t)
	// Equal-power creatures with no armor deal lethal damage to each other; both
	// are judged destroyed from the same combat and removed together.
	att := g.AddToBattleline(testCreature("att", 3), 0)
	def := g.AddToBattleline(testCreature("def", 3), 1)
	if err := g.Fight(0, att, def); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if len(g.Battleline(0)) != 0 || len(g.Battleline(1)) != 0 {
		t.Errorf(
			"both fighters should be destroyed: friendly=%v enemy=%v",
			g.Battleline(0),
			g.Battleline(1),
		)
	}
}

func TestAttackDamage(t *testing.T) {
	valdr := func() CardDefinition {
		return NewCard("Valdr", Brobnar, Creature, Common, WithPower(6),
			WithAttackDamage(AttackDamage{Amount: 2, FlankOnly: true}))
	}

	// A flank bonus adds to the damage dealt when the defender is on a flank.
	g := started(t)
	att := g.AddToBattleline(valdr(), 0)
	flank := g.AddToBattleline(testCreature("flank", 10), 1)
	if err := g.Fight(0, att, flank); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.Damage(flank) != 8 { // 6 power + 2 flank bonus
		t.Errorf("flank defender damage = %d, want 8", g.Damage(flank))
	}

	// A creature in the middle of the line is not on a flank: no bonus.
	g2 := started(t)
	att2 := g2.AddToBattleline(valdr(), 0)
	g2.AddToBattleline(testCreature("left", 10), 1)
	middle := g2.AddToBattleline(testCreature("middle", 10), 1)
	g2.AddToBattleline(testCreature("right", 10), 1)
	if err := g2.Fight(0, att2, middle); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g2.Damage(middle) != 6 { // 6 power, no bonus
		t.Errorf("middle defender damage = %d, want 6", g2.Damage(middle))
	}

	// A Fixed amount replaces power entirely (Ether Spider deals none).
	g3 := started(t)
	spider := g3.AddToBattleline(NewCard("spider", Brobnar, Creature, Common, WithPower(7),
		WithAttackDamage(AttackDamage{Fixed: true, Amount: 0})), 0)
	foe := g3.AddToBattleline(testCreature("foe", 10), 1)
	if err := g3.Fight(0, spider, foe); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g3.Damage(foe) != 0 {
		t.Errorf("fixed-0 attacker dealt %d damage, want 0", g3.Damage(foe))
	}

	// A Fixed amount also replaces the retaliation damage a defender deals back to
	// an attacker (Shadow Self deals none when fought, not only when it fights).
	g4 := started(t)
	att4 := g4.AddToBattleline(testCreature("attacker", 4), 0)
	spider4 := g4.AddToBattleline(NewCard("spider", Brobnar, Creature, Common, WithPower(7),
		WithAttackDamage(AttackDamage{Fixed: true, Amount: 0})), 1)
	if err := g4.Fight(0, att4, spider4); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g4.Damage(att4) != 0 {
		t.Errorf("attacker took %d retaliation from a fixed-0 defender, want 0", g4.Damage(att4))
	}
	if g4.Damage(spider4) != 4 {
		t.Errorf("fixed-0 defender took %d damage, want 4", g4.Damage(spider4))
	}
}

func TestAttackKeywordsPoison(t *testing.T) {
	spyyyder := func() CardDefinition {
		return NewCard("Spyyyder", Dis, Creature, Common, WithPower(2),
			WithKeywords(Skirmish),
			WithAttackKeywords(AttackKeywords{Keywords: []Keyword{Poison}, FlankOnly: true}))
	}

	// Against a flank creature that survives the fight damage, poison destroys it.
	g := started(t)
	att := g.AddToBattleline(spyyyder(), 0)
	flank := g.AddToBattleline(testCreature("flank", 10), 1)
	g.fight(att, flank)
	if g.inPlay(flank) {
		t.Error("flank defender should be destroyed by attack poison")
	}

	// A creature in the middle of the line is not on a flank: no poison, it lives.
	g2 := started(t)
	att2 := g2.AddToBattleline(spyyyder(), 0)
	g2.AddToBattleline(testCreature("left", 10), 1)
	middle := g2.AddToBattleline(testCreature("middle", 10), 1)
	g2.AddToBattleline(testCreature("right", 10), 1)
	g2.fight(att2, middle)
	if !g2.inPlay(middle) {
		t.Error("middle defender should survive: no attack poison off a flank")
	}

	// Poison needs damage to actually land: armor that absorbs it all spares the
	// creature.
	g3 := started(t)
	att3 := g3.AddToBattleline(spyyyder(), 0)
	armored := g3.AddToBattleline(
		NewCard("armored", Logos, Creature, Common, WithPower(10), WithArmor(5)), 1)
	g3.fight(att3, armored)
	if !g3.inPlay(armored) {
		t.Error("fully-absorbed defender should survive: poison needs damage dealt")
	}
}

func TestBeforeFightTrigger(t *testing.T) {
	g := started(t)
	spit := DealDamage{Amount: 1, Target: Target{Kind: TargetEachEnemyCreature}}

	// Defender survives: before-fight damages every enemy, then combat proceeds.
	att := g.AddToBattleline(testCreature("spitter", 3, WithAbility(TriggerBeforeFight, spit)), 0)
	weak := g.AddToBattleline(testCreature("weak", 1), 1)
	tough := g.AddToBattleline(testCreature("tough", 10), 1)
	if err := g.Fight(0, att, tough); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.inPlay(weak) {
		t.Error("before-fight damage should destroy the 1-power enemy")
	}
	if g.Damage(tough) != 4 { // 1 before-fight + 3 combat
		t.Errorf("tough damage = %d, want 4", g.Damage(tough))
	}
	if g.inPlay(att) {
		t.Error("attacker should die to the 10-power defender's retaliation")
	}

	// Defender destroyed before combat: the attack still happens but deals no
	// combat damage, so the attacker takes no retaliation.
	att2 := g.AddToBattleline(testCreature("spitter2", 3, WithAbility(TriggerBeforeFight, spit)), 0)
	frail := g.AddToBattleline(testCreature("frail", 1), 1)
	if err := g.Fight(0, att2, frail); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.inPlay(frail) {
		t.Error("frail defender should be destroyed before combat")
	}
	if g.Damage(att2) != 0 {
		t.Errorf("attacker took %d damage, want 0 (no combat)", g.Damage(att2))
	}
	if !g.State.Cards[att2].Exhausted {
		t.Error("attacker should be exhausted from attacking")
	}
}

// TestCamouflageBlocksNonFlankAttackers pins Camouflage: only a creature on a
// flank of its own battleline may be used to fight the camouflaged host; an
// interior attacker cannot, from either the fight verb or the target list.
func TestCamouflageBlocksNonFlankAttackers(t *testing.T) {
	g := started(t)
	left := g.AddToBattleline(testCreature("left", 4), 0)
	mid := g.AddToBattleline(testCreature("mid", 4), 0)
	g.AddToBattleline(testCreature("right", 4), 0)
	defender := g.AddToBattleline(testCreature("def", 3), 1)
	attachUpgrade(g, defender, NewCard("camo", Untamed, Upgrade, Uncommon,
		WithStatic(StaticModifier{ProtectsFromNonFlank: true})))

	// A non-flank (interior) attacker cannot fight the camouflaged defender.
	if got := g.FightTargets(0, mid); len(got) != 0 {
		t.Errorf("FightTargets for interior attacker = %v, want none", got)
	}
	if err := g.Fight(0, mid, defender); err != ErrNoTarget {
		t.Errorf("interior attacker Fight err = %v, want ErrNoTarget", err)
	}

	// A flank attacker may fight it.
	if got := g.FightTargets(0, left); len(got) != 1 || got[0] != defender {
		t.Errorf("FightTargets for flank attacker = %v, want [%d]", got, defender)
	}
	if err := g.Fight(0, left, defender); err != nil {
		t.Fatalf("flank attacker Fight: %v", err)
	}
}

// TestAfterNeighborFightsTrigger covers the trigger Little Niff uses: it fires
// for a creature used to fight beside the watcher, stealing 1 Æmber, but not for a
// creature that fights away from it.
func TestAfterNeighborFightsTrigger(t *testing.T) {
	steal := StealAember{Amount: 1}

	t.Run("fires when a neighbor is used to fight", func(t *testing.T) {
		g := started(t)
		g.AddToBattleline(
			testCreature("niff", 2, WithAbility(TriggerAfterNeighborFights, steal)), 0)
		attacker := g.AddToBattleline(testCreature("att", 4), 0)
		defender := g.AddToBattleline(testCreature("def", 3), 1)
		g.State.Aember[1] = 2
		if err := g.Fight(0, attacker, defender); err != nil {
			t.Fatalf("Fight: %v", err)
		}
		if g.Aember(0) != 1 {
			t.Errorf("controller aember = %d, want 1", g.Aember(0))
		}
		if g.Aember(1) != 1 {
			t.Errorf("opponent aember = %d, want 1", g.Aember(1))
		}
	})

	t.Run("does not fire for a fight away from the watcher", func(t *testing.T) {
		g := started(t)
		g.AddToBattleline(
			testCreature("niff", 2, WithAbility(TriggerAfterNeighborFights, steal)), 0)
		g.AddToBattleline(testCreature("buffer", 2), 0)
		attacker := g.AddToBattleline(testCreature("att", 4), 0)
		defender := g.AddToBattleline(testCreature("def", 3), 1)
		g.State.Aember[1] = 2
		if err := g.Fight(0, attacker, defender); err != nil {
			t.Fatalf("Fight: %v", err)
		}
		if g.Aember(0) != 0 {
			t.Errorf("controller aember = %d, want 0 (attacker not beside the watcher)",
				g.Aember(0))
		}
	})
}

func TestAfterDestroyedFightingTrigger(t *testing.T) {
	gain := GainAember{Player: Controller, Amount: 1}

	// The attacker survives and destroys the defender: its ability fires.
	g := started(t)
	hunter := g.AddToBattleline(
		testCreature("hunter", 6, WithAbility(TriggerAfterDestroyedFighting, gain)),
		0,
	)
	prey := g.AddToBattleline(testCreature("prey", 1), 1)
	before := g.Aember(0)
	if err := g.Fight(0, hunter, prey); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.inPlay(prey) {
		t.Fatal("prey should be destroyed by the fight")
	}
	if g.Aember(0) != before+1 {
		t.Errorf("survivor aember = %d, want %d", g.Aember(0), before+1)
	}

	// The defender survives while its Hazardous destroys the attacker: the
	// defender's ability fires (for its own controller).
	g2 := started(t)
	weak := g2.AddToBattleline(testCreature("weak", 1), 0)
	guard := g2.AddToBattleline(NewCard("guard", Sanctum, Creature, Common,
		WithPower(6), WithHazardous(5), WithAbility(TriggerAfterDestroyedFighting, gain)), 1)
	if got := g2.Hazardous(guard); got != 5 {
		t.Errorf("guard Hazardous = %d, want 5", got)
	}
	before2 := g2.Aember(1)
	if err := g2.Fight(0, weak, guard); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g2.inPlay(weak) {
		t.Fatal("weak attacker should be destroyed by Hazardous")
	}
	if g2.Aember(1) != before2+1 {
		t.Errorf("defender-controller aember = %d, want %d", g2.Aember(1), before2+1)
	}

	// Both combatants survive (Skirmish attacker, tougher defender): no ability.
	g3 := started(t)
	brawler := g3.AddToBattleline(
		testCreature(
			"brawler",
			3,
			WithKeywords(Skirmish),
			WithAbility(TriggerAfterDestroyedFighting, gain),
		),
		0,
	)
	wall := g3.AddToBattleline(testCreature("wall", 10), 1)
	before3 := g3.Aember(0)
	if err := g3.Fight(0, brawler, wall); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if !g3.inPlay(brawler) || !g3.inPlay(wall) {
		t.Fatal("both combatants should survive")
	}
	if g3.Aember(0) != before3 {
		t.Errorf("no destruction: aember = %d, want %d (unchanged)", g3.Aember(0), before3)
	}
}

func TestSkirmishAndArmorAndPoison(t *testing.T) {
	g := started(t)
	// Skirmish attacker takes no retaliation.
	skirm := g.AddToBattleline(testCreature("skirm", 3, WithKeywords(Skirmish)), 0)
	def := g.AddToBattleline(testCreature("def", 2), 1)
	if err := g.Fight(0, skirm, def); err != nil {
		t.Fatal(err)
	}
	if g.Damage(skirm) != 0 {
		t.Errorf("skirmisher took %d damage, want 0", g.Damage(skirm))
	}
	if want := "skirm is skirmish — it takes no damage in return"; !hasLogLine(g, want) {
		t.Errorf("log = %v, want a line %q", g.LogText(), want)
	}

	// Armor absorbs damage (raw application, no destruction).
	armored := g.AddToBattleline(testCreature("armored", 5, WithArmor(2)), 0)
	g.DealDamage(0, []DamageTarget{{ID: armored, Amount: 3}}) // 2 absorbed, 1 through
	if g.Damage(armored) != 1 {
		t.Errorf("armored damage = %d, want 1", g.Damage(armored))
	}
	g.DealDamage(0, []DamageTarget{{ID: armored, Amount: 0}}) // no-op path

	// Poison is offensive: a poison attacker destroys the creature it damages,
	// however much power that creature had left, while the poison attacker itself
	// survives (Skirmish spares it the return damage). The log names poison as the
	// killer.
	pois := g.AddToBattleline(testCreature("pois", 3, WithKeywords(Skirmish, Poison)), 0)
	prey := g.AddToBattleline(testCreature("prey", 10), 1)
	if err := g.Fight(0, pois, prey); err != nil {
		t.Fatal(err)
	}
	if g.inPlay(prey) {
		t.Error("a creature a poison attacker damaged should be destroyed")
	}
	if !g.inPlay(pois) {
		t.Error("a poison creature is not destroyed by dealing damage")
	}
	if want := "pois's poison is lethal to prey"; !hasLogLine(g, want) {
		t.Errorf("log = %v, want a line %q", g.LogText(), want)
	}
}

// hasLogLine reports whether the game log holds a line exactly matching want.
func hasLogLine(g *Game, want string) bool {
	for _, line := range g.LogText() {
		if line == want {
			return true
		}
	}
	return false
}

func TestPoisonIsOffensive(t *testing.T) {
	// A poison defender's return damage destroys the attacker that fought into it.
	t.Run("return damage from a poison defender kills the attacker", func(t *testing.T) {
		g := started(t)
		attacker := g.AddToBattleline(testCreature("attacker", 8), 0)
		venom := g.AddToBattleline(testCreature("venom", 4, WithKeywords(Poison)), 1)
		if err := g.Fight(0, attacker, venom); err != nil {
			t.Fatal(err)
		}
		if g.inPlay(attacker) {
			t.Error("an attacker that fought into poison should be destroyed")
		}
		if want := "venom's poison is lethal to attacker"; !hasLogLine(g, want) {
			t.Errorf("log = %v, want a line %q", g.LogText(), want)
		}
	})

	// Armor that soaks the whole hit leaves no landing damage, so poison does not
	// destroy the creature.
	t.Run("armor that absorbs the whole hit spares the target", func(t *testing.T) {
		g := started(t)
		pois := g.AddToBattleline(testCreature("pois", 1, WithKeywords(Skirmish, Poison)), 0)
		walled := g.AddToBattleline(testCreature("walled", 5, WithArmor(3)), 1)
		if err := g.Fight(0, pois, walled); err != nil {
			t.Fatal(err)
		}
		if !g.inPlay(walled) {
			t.Error("armor soaked the whole hit, so poison should not destroy the creature")
		}
	})

	// Splash damage from a poison attacker destroys the neighbors it hits.
	t.Run("splash damage carries poison to neighbors", func(t *testing.T) {
		g := started(t)
		pois := g.AddToBattleline(
			testCreature("pois", 2, WithKeywords(Skirmish, Poison), WithSplashAttack(1)),
			0,
		)
		target := g.AddToBattleline(testCreature("target", 9), 1)
		bystander := g.AddToBattleline(testCreature("bystander", 9), 1)
		if err := g.Fight(0, pois, target); err != nil {
			t.Fatal(err)
		}
		if g.inPlay(bystander) {
			t.Error("a neighbor a poison creature splashed should be destroyed")
		}
	})
}

func TestFightDamageRedirect(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetChooser(0, orderLastChooser{}) // the Before Fight choice picks the last candidate
	gabos := g.AddToBattleline(
		testCreature(
			"gabos",
			5,
			WithAbility(
				TriggerBeforeFight,
				RedirectFightDamage{Target: Target{Kind: TargetChosenCreature}},
			),
		),
		0,
	)
	def := g.AddToBattleline(testCreature("defender", 3), 1)
	bystander := g.AddToBattleline(testCreature("bystander", 2), 1)

	g.fight(gabos, def)

	// Candidates are [gabos, defender, bystander]; the chooser takes the last, so
	// Gabos's fight damage lands on the bystander, not the defender it fights.
	if g.inPlay(bystander) {
		t.Error("bystander should take Gabos's 5 fight damage and be destroyed")
	}
	if g.Damage(def) != 0 {
		t.Errorf("defender should take no fight damage; got %d", g.Damage(def))
	}
	if g.Damage(gabos) != 3 {
		t.Errorf("Gabos should still take the defender's 3 damage back; got %d", g.Damage(gabos))
	}
	if g.State.FightDamageRedirect != 0 {
		t.Error("redirect should be cleared after the fight")
	}
}

func TestBeforeFightCanCancelFight(t *testing.T) {
	g := started(t) // Brobnar is active
	top := g.AddToDeck(NewCard("Brobnar Top", Brobnar, Tactic, Common), 0)
	attacker := g.AddToBattleline(NewCard("evader", Brobnar, Creature, Common,
		WithPower(5),
		WithAssault(2),
		WithAbility(TriggerBeforeFight, Sequence{Effects: []Effect{
			DiscardTop{},
			Conditional{Cond: ItIs{House: activeHouse}, Then: CancelFight{}},
		}}),
		WithAbility(TriggerAfterFight, GainAember{Player: Controller, Amount: 1}),
	), 0)
	defender := g.AddToBattleline(
		NewCard("hazard", Brobnar, Creature, Common, WithPower(8), WithHazardous(3)),
		1,
	)

	if err := g.Fight(0, attacker, defender); err != nil {
		t.Fatalf("Fight: %v", err)
	}

	if discard := g.Discard(0); len(discard) != 1 || discard[0] != top {
		t.Errorf("discard = %v, want discarded top card %d", discard, top)
	}
	if !g.State.Cards[attacker].Exhausted {
		t.Error("attacker should still be exhausted from being used to fight")
	}
	if g.Damage(attacker) != 0 {
		t.Errorf("attacker damage = %d, want 0 (Hazardous skipped)", g.Damage(attacker))
	}
	if g.Damage(defender) != 0 {
		t.Errorf(
			"defender damage = %d, want 0 (Assault and fight damage skipped)",
			g.Damage(defender),
		)
	}
	if g.Aember(0) != 0 {
		t.Errorf("after-fight ability fired: Æmber = %d, want 0", g.Aember(0))
	}
	if g.State.FightCancelled {
		t.Error("fight cancellation flag should be cleared after the fight")
	}
}

func TestBeforeFightCancelMissStillFights(t *testing.T) {
	g := started(t) // Brobnar is active
	top := g.AddToDeck(NewCard("Mars Top", Mars, Tactic, Common), 0)
	attacker := g.AddToBattleline(NewCard("evader", Brobnar, Creature, Common,
		WithPower(9),
		WithAbility(TriggerBeforeFight, Sequence{Effects: []Effect{
			DiscardTop{},
			Conditional{Cond: ItIs{House: activeHouse}, Then: CancelFight{}},
		}}),
		WithAbility(TriggerAfterFight, GainAember{Player: Controller, Amount: 1}),
	), 0)
	defender := g.AddToBattleline(testCreature("defender", 8), 1)

	if err := g.Fight(0, attacker, defender); err != nil {
		t.Fatalf("Fight: %v", err)
	}

	if discard := g.Discard(0); len(discard) != 1 || discard[0] != top {
		t.Errorf("discard = %v, want discarded top card %d", discard, top)
	}
	if g.Damage(attacker) != 8 {
		t.Errorf("attacker damage = %d, want 8", g.Damage(attacker))
	}
	if g.inPlay(defender) {
		t.Error("defender should be destroyed by the uncancelled fight")
	}
	if g.Aember(0) != 1 {
		t.Errorf("after-fight ability did not fire: Æmber = %d, want 1", g.Aember(0))
	}
}

func TestAssaultAndHazardous(t *testing.T) {
	// Assault: the attacker deals its Assault to the defender before fight damage.
	g := NewGame("A", "B", 1)
	att := g.AddToBattleline(
		NewCard("assailant", Brobnar, Creature, Common, WithPower(3), WithAssault(2)),
		0,
	)
	def := g.AddToBattleline(testCreature("beefy", 10), 1)
	g.fight(att, def)
	if g.Damage(def) != 5 { // 2 assault + 3 fight
		t.Errorf("assault: defender damage = %d, want 5 (2 assault + 3 fight)", g.Damage(def))
	}

	// Hazardous can destroy the attacker before any combat damage is exchanged.
	g2 := NewGame("A", "B", 1)
	frail := g2.AddToBattleline(testCreature("frail", 4), 0)
	thorns := g2.AddToBattleline(
		NewCard("thorns", Untamed, Creature, Common, WithPower(6), WithHazardous(5)),
		1,
	)
	g2.fight(frail, thorns)
	if g2.inPlay(frail) {
		t.Error("hazardous should destroy the frail attacker before combat")
	}
	if g2.Damage(thorns) != 0 {
		t.Errorf("no combat should occur; thorns damage = %d, want 0", g2.Damage(thorns))
	}

	// Assault granted by an upgrade stacks (the accessor sums upgrade bonuses).
	g3 := NewGame("A", "B", 1)
	host := g3.AddToBattleline(testCreature("host", 3), 0)
	wall := g3.AddToBattleline(testCreature("wall", 10), 1)
	attachUpgrade(
		g3,
		host,
		NewCard("bearway", Untamed, Upgrade, Common, WithStatic(StaticModifier{AssaultBonus: 2})),
	)
	g3.fight(host, wall)
	if g3.Damage(wall) != 5 {
		t.Errorf("upgraded assault: wall damage = %d, want 5", g3.Damage(wall))
	}

	// Hazardous granted by an upgrade stacks too.
	g4 := NewGame("A", "B", 1)
	weak := g4.AddToBattleline(testCreature("weak", 1), 0)
	guard := g4.AddToBattleline(testCreature("guard", 6), 1)
	attachUpgrade(
		g4,
		guard,
		NewCard("flames", Dis, Upgrade, Common, WithStatic(StaticModifier{HazardousBonus: 2})),
	)
	g4.fight(weak, guard)
	if g4.inPlay(weak) {
		t.Error("upgraded hazardous should destroy the weak attacker")
	}
}

// TestAssaultDestroysTrigger covers the pre-fight "after a creature is destroyed
// by this creature's assault damage" trigger (Skoll's Assault).
func TestAssaultDestroysTrigger(t *testing.T) {
	g := NewGame("A", "B", 1)
	att := g.AddToBattleline(
		NewCard("skoll", Brobnar, Creature, Common, WithPower(3), WithAssault(5),
			WithAbility(TriggerAfterAssaultDestroys, GainAember{Player: Controller, Amount: 1})),
		0,
	)
	def := g.AddToBattleline(testCreature("prey", 4), 1)
	g.fight(att, def)
	if g.inPlay(def) {
		t.Error("assault should destroy the defender before the fight")
	}
	if g.Aember(0) != 1 {
		t.Errorf("aember = %d, want 1 (assault-destroys trigger fired)", g.Aember(0))
	}
}

// TestPreFightDamageNarratesSource checks that Assault and Hazardous narrate the
// creature and keyword value that dealt the pre-fight damage, not a bare
// "takes N damage" line.
func TestPreFightDamageNarratesSource(t *testing.T) {
	g := NewGame("A", "B", 1)
	att := g.AddToBattleline(
		NewCard("imp", Dis, Creature, Common, WithPower(6), WithAssault(2)),
		0,
	)
	def := g.AddToBattleline(
		NewCard("director", Logos, Creature, Common, WithPower(6), WithHazardous(3)),
		1,
	)
	g.fight(att, def)

	assault := "imp's 2 Assault deals 2 damage to director"
	hazardous := "director's 3 Hazardous deals 3 damage to imp"
	var gotAssault, gotHazardous bool
	for _, line := range g.LogText() {
		switch line {
		case assault:
			gotAssault = true
		case hazardous:
			gotHazardous = true
		}
	}
	if !gotAssault {
		t.Errorf("log = %v, want a line %q", g.LogText(), assault)
	}
	if !gotHazardous {
		t.Errorf("log = %v, want a line %q", g.LogText(), hazardous)
	}
}

func TestSplashAttack(t *testing.T) {
	// Splash-attack hits each neighbor of the fought creature at the same time as
	// the fight damage the fought creature takes.
	g := NewGame("A", "B", 1)
	att := g.AddToBattleline(
		NewCard("splasher", Brobnar, Creature, Common, WithPower(5), WithSplashAttack(2)),
		0,
	)
	left := g.AddToBattleline(testCreature("left", 20), 1)
	mid := g.AddToBattleline(testCreature("mid", 20), 1)
	right := g.AddToBattleline(testCreature("right", 20), 1)
	g.fight(att, mid)
	if g.Damage(mid) != 5 {
		t.Errorf("fought creature damage = %d, want 5", g.Damage(mid))
	}
	if g.Damage(left) != 2 || g.Damage(right) != 2 {
		t.Errorf("neighbor damage = %d/%d, want 2/2", g.Damage(left), g.Damage(right))
	}

	// A fought flank creature has only one neighbor: only that one is splashed.
	g2 := NewGame("A", "B", 1)
	att2 := g2.AddToBattleline(
		NewCard("splasher2", Brobnar, Creature, Common, WithPower(5), WithSplashAttack(3)),
		0,
	)
	flank := g2.AddToBattleline(testCreature("flank", 20), 1)
	inner := g2.AddToBattleline(testCreature("inner", 20), 1)
	g2.fight(att2, flank)
	if g2.Damage(inner) != 3 {
		t.Errorf("flank neighbor damage = %d, want 3", g2.Damage(inner))
	}

	// Splash-attack granted by an upgrade stacks (the accessor sums upgrade bonuses).
	g3 := NewGame("A", "B", 1)
	host := g3.AddToBattleline(testCreature("host", 5), 0)
	tgt := g3.AddToBattleline(testCreature("tgt", 20), 1)
	nb := g3.AddToBattleline(testCreature("nb", 20), 1)
	attachUpgrade(
		g3,
		host,
		NewCard(
			"spray",
			Untamed,
			Upgrade,
			Common,
			WithStatic(StaticModifier{SplashAttackBonus: 2}),
		),
	)
	if got := g3.splashAttack(host); got != 2 {
		t.Errorf("upgraded splash-attack = %d, want 2", got)
	}
	g3.fight(host, tgt)
	if g3.Damage(nb) != 2 {
		t.Errorf("upgraded splash neighbor damage = %d, want 2", g3.Damage(nb))
	}
}

func TestElusive(t *testing.T) {
	// The first fight against an elusive creature deals no fight damage either way.
	g := NewGame("A", "B", 1)
	att := g.AddToBattleline(testCreature("att", 5), 0)
	ghost := g.AddToBattleline(
		NewCard("ghost", Shadows, Creature, Common, WithPower(3), WithKeywords(Elusive)),
		1,
	)
	g.fight(att, ghost)
	if g.Damage(ghost) != 0 {
		t.Errorf("elusive defender damage = %d, want 0", g.Damage(ghost))
	}
	if g.Damage(att) != 0 {
		t.Errorf("attacker damage = %d, want 0 (no fight damage at all)", g.Damage(att))
	}

	// Elusive is spent for the turn: a second fight resolves normally.
	att2 := g.AddToBattleline(testCreature("att2", 5), 0)
	g.fight(att2, ghost)
	if g.inPlay(ghost) {
		t.Error("second fight should destroy the no-longer-elusive defender")
	}
	if g.Damage(att2) != 3 {
		t.Errorf("second attacker damage = %d, want 3", g.Damage(att2))
	}

	// StartTurn refreshes it.
	g2 := NewGame("A", "B", 1)
	a1 := g2.AddToBattleline(testCreature("a1", 5), 0)
	dodger := g2.AddToBattleline(
		NewCard("dodger", Shadows, Creature, Common, WithPower(3), WithKeywords(Elusive)),
		1,
	)
	g2.fight(a1, dodger)
	g2.StartTurn(1)
	g2.StartTurn(0)
	a2 := g2.AddToBattleline(testCreature("a2", 5), 0)
	g2.fight(a2, dodger)
	if g2.Damage(dodger) != 0 {
		t.Errorf("after StartTurn, elusive damage = %d, want 0 (refreshed)", g2.Damage(dodger))
	}

	// Elusive stops only fight damage: Hazardous still hits the attacker, and the
	// keyword is spent even though the fight never happens.
	g3 := NewGame("A", "B", 1)
	frail := g3.AddToBattleline(testCreature("frail", 2), 0)
	thorns := g3.AddToBattleline(
		NewCard("thorns", Untamed, Creature, Common,
			WithPower(6), WithHazardous(5), WithKeywords(Elusive)),
		1,
	)
	g3.fight(frail, thorns)
	if g3.inPlay(frail) {
		t.Error("hazardous should still destroy the attacker of an elusive creature")
	}
	if !g3.ElusiveSpent(thorns) {
		t.Error("elusive should be spent even when the fight never happens")
	}

	// A fight that never gets past "Before Fight" does not spend Elusive: the
	// keyword replaces the fight, and there was no fight to replace.
	g4 := NewGame("A", "B", 1)
	sigil := g4.AddToBattleline(NewCard("sigil", Sanctum, Creature, Common,
		WithPower(3), WithAbility(TriggerBeforeFight, CancelFight{})), 0)
	ghoul := g4.AddToBattleline(
		NewCard("ghoul", Shadows, Creature, Common, WithPower(3), WithKeywords(Elusive)),
		1,
	)
	g4.fight(sigil, ghoul)
	if g4.ElusiveSpent(ghoul) {
		t.Error("a cancelled fight should not spend elusive")
	}
}

func TestFightRestriction(t *testing.T) {
	g := started(t) // Brobnar active
	stunnedOnly := Target{Kind: TargetEachCreature}.Stunned()
	att := g.AddToBattleline(
		NewCard("twig", Brobnar, Creature, Common, WithPower(7), WithFightRestriction(stunnedOnly)),
		0,
	)
	def := g.AddToBattleline(testCreature("def", 3), 1)

	// An unstunned enemy is not a legal target.
	if err := g.Fight(0, att, def); err != ErrNoTarget {
		t.Errorf("fight unstunned = %v, want ErrNoTarget", err)
	}

	// Once the enemy is stunned it can be fought.
	g.State.Cards[def].Stunned = true
	if err := g.Fight(0, att, def); err != nil {
		t.Fatalf("fight stunned: %v", err)
	}
}

func TestFightTargets(t *testing.T) {
	// Normal: a usable attacker may fight every enemy creature.
	g := started(t)
	att := g.AddToBattleline(testCreature("att", 3), 0)
	a := g.AddToBattleline(testCreature("a", 3), 1)
	b := g.AddToBattleline(testCreature("b", 3), 1)
	if got := g.FightTargets(0, att); len(got) != 2 || got[0] != a || got[1] != b {
		t.Errorf("FightTargets = %v, want [%d %d]", got, a, b)
	}

	// No enemy creatures: nothing to fight.
	g2 := started(t)
	lone := g2.AddToBattleline(testCreature("lone", 3), 0)
	if got := g2.FightTargets(0, lone); got != nil {
		t.Errorf("FightTargets with no enemies = %v, want nil", got)
	}

	// Barred from fighting: no targets even with enemies present.
	g3 := started(t)
	barred := g3.AddToBattleline(testCreature("barred", 3), 0)
	g3.AddToBattleline(testCreature("foe", 3), 1)
	g3.State.CannotFight[0].Value = true
	if got := g3.FightTargets(0, barred); got != nil {
		t.Errorf("barred FightTargets = %v, want nil", got)
	}

	// Wrong house without a grant: not usable, so no targets.
	g4 := started(t) // Brobnar active
	off := g4.AddToBattleline(NewCard("off", Sanctum, Creature, Common, WithPower(3)), 0)
	g4.AddToBattleline(testCreature("foe", 3), 1)
	if got := g4.FightTargets(0, off); got != nil {
		t.Errorf("off-house FightTargets = %v, want nil", got)
	}

	// A fight grant forgives the wrong house: the off-house creature may fight.
	g4.State.MayFightHouse[0] = Sanctum
	if got := g4.FightTargets(0, off); len(got) != 1 {
		t.Errorf("granted FightTargets = %v, want one target", got)
	}

	// A fight restriction narrows the targets to the enemies it allows.
	g5 := started(t)
	picky := g5.AddToBattleline(
		NewCard("picky", Brobnar, Creature, Common, WithPower(7),
			WithFightRestriction(Target{Kind: TargetEachCreature}.Stunned())), 0)
	awake := g5.AddToBattleline(testCreature("awake", 3), 1)
	stunned := g5.AddToBattleline(testCreature("stunned", 3), 1)
	g5.State.Cards[stunned].Stunned = true
	got := g5.FightTargets(0, picky)
	if len(got) != 1 || got[0] != stunned {
		t.Errorf(
			"restricted FightTargets = %v, want [%d] (only the stunned enemy, not %d)",
			got,
			stunned,
			awake,
		)
	}
}

func TestTauntProtectsNeighbors(t *testing.T) {
	g := started(t)
	att := g.AddToBattleline(testCreature("att", 3), 0)
	left := g.AddToBattleline(testCreature("left", 3), 1)
	taunter := g.AddToBattleline(
		NewCard("taunter", Brobnar, Creature, Common, WithPower(3), WithKeywords(Taunt)), 1)
	right := g.AddToBattleline(testCreature("right", 3), 1)
	far := g.AddToBattleline(testCreature("far", 3), 1)

	if got := g.FightTargets(0, att); len(got) != 2 || got[0] != taunter || got[1] != far {
		t.Errorf("FightTargets = %v, want [%d %d]", got, taunter, far)
	}
	if err := g.Fight(0, att, left); err != ErrNoTarget {
		t.Errorf("fighting a taunt neighbor: err = %v, want ErrNoTarget", err)
	}
	if err := g.Fight(0, att, right); err != ErrNoTarget {
		t.Errorf("fighting a taunt neighbor: err = %v, want ErrNoTarget", err)
	}
	if err := g.Fight(0, att, taunter); err != nil {
		t.Fatalf("fighting the taunt creature itself: %v", err)
	}
}

// A creature with taunt of its own is not shielded by a taunt neighbor.
func TestTauntDoesNotProtectTauntNeighbors(t *testing.T) {
	g := started(t)
	att := g.AddToBattleline(testCreature("att", 3), 0)
	taunt := func(name string) CardDefinition {
		return NewCard(name, Brobnar, Creature, Common, WithPower(3), WithKeywords(Taunt))
	}
	a := g.AddToBattleline(taunt("a"), 1)
	b := g.AddToBattleline(taunt("b"), 1)
	if got := g.FightTargets(0, att); len(got) != 2 || got[0] != a || got[1] != b {
		t.Errorf("FightTargets = %v, want [%d %d]", got, a, b)
	}
}

func TestTauntShielded(t *testing.T) {
	g := started(t)
	left := g.AddToBattleline(testCreature("left", 3), 1)
	taunter := g.AddToBattleline(
		NewCard("taunter", Brobnar, Creature, Common, WithPower(3), WithKeywords(Taunt)), 1)
	right := g.AddToBattleline(testCreature("right", 3), 1)
	far := g.AddToBattleline(testCreature("far", 3), 1)

	if !g.TauntShielded(left) {
		t.Error("left neighbor of a taunter should be shielded")
	}
	if !g.TauntShielded(right) {
		t.Error("right neighbor of a taunter should be shielded")
	}
	if g.TauntShielded(taunter) {
		t.Error("a taunter is not shielded by itself")
	}
	if g.TauntShielded(far) {
		t.Error("a creature two away from the taunter should not be shielded")
	}
}

// A card that takes damage for other creatures absorbs the damage aimed at each
// creature its Target names, and its own damage is never redirected onward.
func TestTakesDamageFor(t *testing.T) {
	g := started(t)
	shield := NewCard("Shield", Shadows, Creature, Common, WithPower(9),
		WithTakesDamageFor(Target{Kind: TargetEachCreature}.Neighboring()))
	ward := g.AddToBattleline(testCreature("ward", 2), 0)
	sid := g.AddToBattleline(shield, 0)
	far := g.AddToBattleline(testCreature("far", 2), 0)
	// A second shield beside the first would bounce damage back and forth if a
	// redirect could chain; it must not.
	other := g.AddToBattleline(shield, 0)

	g.dealDamage(0, DamageTarget{ID: ward, Amount: 2})
	if g.Damage(ward) != 0 || g.Damage(sid) != 2 {
		t.Errorf("ward=%d shield=%d, want 0 and 2", g.Damage(ward), g.Damage(sid))
	}
	if len(g.Battleline(0)) != 4 {
		t.Error("the shielded creature should have survived")
	}

	// A creature the shield does not reach takes its own damage.
	g.dealDamage(0, DamageTarget{ID: far, Amount: 1})
	if g.Damage(far) != 0 || g.Damage(sid) != 3 {
		t.Errorf("far=%d shield=%d, want 0 and 3", g.Damage(far), g.Damage(sid))
	}
	if g.Damage(other) != 0 {
		t.Errorf("other shield=%d, want 0 (the nearer shield took it)", g.Damage(other))
	}

	// Damage aimed at a shield stays on it rather than hopping to its neighbor.
	g.dealDamage(0, DamageTarget{ID: sid, Amount: 3})
	if g.Damage(sid) != 6 {
		t.Errorf("shield=%d, want 6 (a redirect never chains)", g.Damage(sid))
	}

	// A creature already out of play is not redirected either; the damage is
	// simply dropped.
	g.destroyEach(0, []LocalID{ward})
	g.dealDamage(0, DamageTarget{ID: ward, Amount: 1})
	if g.Damage(sid) != 6 {
		t.Errorf("shield=%d, want 6 (damage to a card out of play goes nowhere)",
			g.Damage(sid))
	}
}

// Redirection is the last step of the damage chart, so both creatures' armor
// absorbs: the shielded creature's armor reduces the damage it is dealt, and the
// shield's armor reduces again from what it receives.
func TestTakesDamageForAfterArmor(t *testing.T) {
	g := started(t)
	shield := NewCard("Shield", Shadows, Creature, Common, WithPower(9), WithArmor(2),
		WithTakesDamageFor(Target{Kind: TargetEachCreature}.Neighboring()))
	ward := g.AddToBattleline(NewCard("ward", Shadows, Creature, Common,
		WithPower(6), WithArmor(2)), 0)
	sid := g.AddToBattleline(shield, 0)

	g.dealDamage(0, DamageTarget{ID: ward, Amount: 3})
	if g.Damage(ward) != 0 || g.Damage(sid) != 0 {
		t.Errorf("ward=%d shield=%d, want 0 and 0 (3 - 2 armor, then the last point"+
			" absorbed by the shield's own armor)", g.Damage(ward), g.Damage(sid))
	}

	// Both creatures' armor is spent down by what it absorbed, so the next hit is
	// only partly blunted.
	g.dealDamage(0, DamageTarget{ID: ward, Amount: 5})
	if g.Damage(ward) != 0 || g.Damage(sid) != 4 {
		t.Errorf("ward=%d shield=%d, want 0 and 4", g.Damage(ward), g.Damage(sid))
	}

	// Armor that fully absorbs leaves nothing to redirect, so the shield is untouched.
	tough := g.AddToBattleline(NewCard("tough", Shadows, Creature, Common,
		WithPower(6), WithArmor(3)), 0)
	g.dealDamage(0, DamageTarget{ID: tough, Amount: 3})
	if g.Damage(tough) != 0 || g.Damage(sid) != 4 {
		t.Errorf("tough=%d shield=%d, want 0 and 4 (its own armor swallowed it whole)",
			g.Damage(tough), g.Damage(sid))
	}
}
