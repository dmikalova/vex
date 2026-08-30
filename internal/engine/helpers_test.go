package engine

import "testing"

// ---- shared test helpers ----

// started returns a game with Alice (player 0) active and Brobnar chosen.
func started(t *testing.T) *Game {
	t.Helper()
	g := NewGame("Alice", "Bob", 1)
	g.BeginTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	return g
}

func testCreature(name string, power int, opts ...CardOption) CardDefinition {
	base := []CardOption{WithPower(power)}
	return NewCard(name, Brobnar, Creature, Common, append(base, opts...)...)
}

// The engine tests exercise the rules against a few example blueprints, defined
// here rather than imported from the card database so package engine stays free of
// any dependency on the card packages that import it.

func exGiant() CardDefinition {
	return NewCard("Brobnar Giant", Brobnar, Creature, Rare,
		WithPower(5), WithArmor(0), WithTraits("Giant"),
		WithAbility(TriggerAfterForgeKey, DealDamage{Amount: 2, Target: Target{Kind: TargetEachEnemyCreature}}))
}

func exBruteStrength() CardDefinition {
	return NewCard("Brute Strength", Brobnar, Upgrade, Uncommon,
		WithAemberBonus(1), WithStatic(StaticModifier{PowerBonus: 5}))
}

func exBattleFury() CardDefinition {
	return NewCard("Battle Fury", Brobnar, Tactic, Common,
		WithAemberBonus(1),
		WithAbility(TriggerAfterPlay, OnChooseCreature{Target: Target{Kind: TargetChosenFriendlyCreature}, Verbs: []CreatureVerb{ReadyVerb{}, FightVerb{}}}))
}

func exAutocannon() CardDefinition {
	return NewCard("Autocannon", Brobnar, Artifact, Rare,
		WithAemberBonus(1), WithTraits("Weapon"),
		WithAbility(TriggerAfterCreatureEnters, DealDamage{Amount: 1, Target: Target{Kind: TargetTriggeringCreature}}))
}

// ---- small test-only lookups ----

func handIdx(g *Game, player int, name string) int {
	for i, id := range g.Hand(player) {
		if g.Name(id) == name {
			return i
		}
	}
	return -1
}

func handIdxByID(g *Game, player int, want LocalID) int {
	for i, id := range g.Hand(player) {
		if id == want {
			return i
		}
	}
	return -1
}

// orderLastChooser always picks the last candidate, reversing an ordering.
type orderLastChooser struct{}

func (orderLastChooser) ChooseCreature(_, _ string, c []LocalID) (LocalID, bool) {
	return c[len(c)-1], true
}

// orderRejectChooser refuses to pick, so ordering falls back to the given order.
type orderRejectChooser struct{}

func (orderRejectChooser) ChooseCreature(_, _ string, _ []LocalID) (LocalID, bool) {
	return 0, false
}

// orderAllChooser implements Orderer, arranging ids in a single call (reversing
// them) instead of being asked to pick the next id repeatedly.
type orderAllChooser struct{}

func (orderAllChooser) ChooseCreature(_, _ string, c []LocalID) (LocalID, bool) {
	return c[0], true
}

func (orderAllChooser) OrderCreatures(_, _ string, ids []LocalID) []LocalID {
	out := append([]LocalID(nil), ids...)
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

// attachUpgrade registers an upgrade and attaches it to a host creature.
func attachUpgrade(g *Game, host LocalID, def CardDefinition) {
	up := g.Register(def, g.owner(host))
	core := &g.State.Cards[host]
	core.Upgrades[core.UpgradeCount] = up
	core.UpgradeCount++
}
