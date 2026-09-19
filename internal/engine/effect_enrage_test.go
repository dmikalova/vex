package engine

import "testing"

func TestEnrageEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	foe1 := g.AddToBattleline(testCreature("foe1", 3), 1)
	foe2 := g.AddToBattleline(testCreature("foe2", 3), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := Enrage{Target: Target{Kind: TargetEachEnemyCreature}}
	if e.Text() != "enrage each enemy creature" {
		t.Errorf("enrage text = %q", e.Text())
	}
	e.Resolve(ctx)
	if !g.Enraged(foe1) || !g.Enraged(foe2) {
		t.Error("enrage should enrage each enemy creature")
	}
	if g.Enraged(src) {
		t.Error("enraging enemy creatures should not touch a friendly creature")
	}

	// An enrage that finds its target already enraged still logs the choice, just
	// without a state change.
	entries := len(g.Log)
	e.Resolve(ctx)
	if len(g.Log) == entries {
		t.Error("re-enraging an already-enraged creature should still log the choice")
	}
}

func TestEnrageValidate(t *testing.T) {
	if (Enrage{}).validate() == nil {
		t.Error("Enrage with no target should fail validation")
	}
	if (Enrage{Target: Target{Kind: TargetEachEnemyCreature}}).validate() != nil {
		t.Error("Enrage with a target should validate")
	}
}

// TestEnrageClearedByFight pins that only the creature used to fight loses its
// enrage; the defender it fought is not used, so it keeps its own.
func TestEnrageClearedByFight(t *testing.T) {
	g := NewGame("A", "B", 1)
	att := g.AddToBattleline(testCreature("att", 3, WithArmor(5)), 0)
	def := g.AddToBattleline(testCreature("def", 3, WithArmor(5)), 1)
	g.State.Cards[att].Enraged = true
	g.State.Cards[def].Enraged = true

	g.fight(att, def)

	if g.Enraged(att) {
		t.Error("a creature used to fight should lose its enrage")
	}
	if !g.Enraged(def) {
		t.Error("the defender is not used to fight, so it keeps its enrage")
	}
}

// TestEnrageClearedByFightingElusive is the reported case: an Elusive defender
// turns the fight's damage aside, but the fight still happened, so the attacker
// still loses its enrage (Gizelhart's Zealot fighting an elusive creature).
func TestEnrageClearedByFightingElusive(t *testing.T) {
	g := NewGame("A", "B", 1)
	att := g.AddToBattleline(testCreature("att", 4), 0)
	g.State.Cards[att].Enraged = true
	ghost := g.AddToBattleline(
		NewCard("ghost", Shadows, Creature, Common, WithPower(3), WithKeywords(Elusive)),
		1,
	)

	g.fight(att, ghost)

	if g.Damage(att) != 0 || g.Damage(ghost) != 0 {
		t.Fatalf("elusive should turn all fight damage aside: att=%d ghost=%d",
			g.Damage(att), g.Damage(ghost))
	}
	if g.Enraged(att) {
		t.Error("enrage is lost after fighting an elusive creature, damage or not")
	}
}

// TestEnrageKeptOnReap pins that reaping does not remove enrage — only fighting
// does (an enraged creature made to reap by Ghosthawk keeps its enrage).
func TestEnrageKeptOnReap(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	g.State.Cards[c].Enraged = true

	g.reapWith(c)

	if !g.Enraged(c) {
		t.Error("reaping should not remove enrage")
	}
}

// TestEnrageClearedBeforeFightAbilityReenrages pins the ordering: the used-to-fight
// removal runs before the post-fight window, so a "Fight:" ability that re-enrages
// the fighter (Gladiodontus) wins and it ends enraged.
func TestEnrageClearedBeforeFightAbilityReenrages(t *testing.T) {
	g := NewGame("A", "B", 1)
	att := g.AddToBattleline(
		testCreature("glad", 3, WithArmor(5),
			WithAbility(TriggerAfterFight, Enrage{Target: Target{Kind: TargetThisCreature}})),
		0,
	)
	def := g.AddToBattleline(testCreature("def", 3, WithArmor(5)), 1)
	g.State.Cards[att].Enraged = true

	g.fight(att, def)

	if !g.Enraged(att) {
		t.Error("a Fight: ability that re-enrages the fighter should leave it enraged")
	}
}
