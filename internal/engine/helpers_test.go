package engine

import "testing"

// ---- shared test helpers ----

// started returns a game with Alice (player 0) active and Brobnar chosen.
func started(t *testing.T) *Game {
	t.Helper()
	g := NewGame("Alice", "Bob", 1)
	g.StartTurn(0)
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
	return NewCard(
		"Brobnar Giant",
		Brobnar,
		Creature,
		Rare,
		WithPower(5),
		WithArmor(0),
		WithTraits(Giant),
		WithAbility(
			TriggerAfterForgeKey,
			DealDamage{
				Amount: 2,
				Target: Target{Kind: TargetEachEnemyCreature},
			},
		),
	)
}

func exBruteStrength() CardDefinition {
	return NewCard("Brute Strength", Brobnar, Upgrade, Uncommon,
		WithBonus(BonusAember), WithStatic(StaticModifier{PowerBonus: 5}))
}

func exBattleFury() CardDefinition {
	return NewCard(
		"Battle Fury",
		Brobnar,
		Tactic,
		Common,
		WithBonus(BonusAember),
		WithAbility(
			TriggerAfterPlay,
			OnChooseCreature{
				Target: Target{Kind: TargetChosenFriendlyCreature},
				Verbs:  []CreatureVerb{ReadyVerb{}, FightVerb{}},
			},
		),
	)
}

func exAutocannon() CardDefinition {
	return NewCard(
		"Autocannon",
		Brobnar,
		Artifact,
		Rare,
		WithBonus(BonusAember),
		WithTraits(Weapon),
		WithAbility(
			TriggerAfterCreatureEnters,
			DealDamage{
				Amount: 1,
				Target: Target{Kind: TargetTriggeringCreature},
			},
		),
	)
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

// resolveLastingWindow resolves the registry reactions responding to event for
// subject as a standalone window — the test path for exercising a lasting reaction
// in isolation, routed through the same order-and-resolve machinery a card window
// folds them into (lastingReactions gathers, orderTriggered orders, resolveWindow
// resolves).
func (g *Game) resolveLastingWindow(event Event, actor int, subject LocalID) {
	g.resolveWindow(
		g.orderTriggered(actor, g.lastingReactions(event, actor, subject)),
	)
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

// declineOptionChooser answers every labeled option with the second choice (a "No"
// on a Yes/No prompt), while taking the first creature candidate — for testing the
// declined branch of an optional prompt.
type declineOptionChooser struct{}

func (declineOptionChooser) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	if len(cands) == 0 {
		return 0, false
	}
	return cands[0], true
}

func (declineOptionChooser) ChooseOption(_, _ string, _ []string) int { return 1 }

// idQueueChooser pops the next scripted id for each choice, falling back to the
// first candidate once the queue empties.
type idQueueChooser struct{ ids []LocalID }

func (c *idQueueChooser) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	if len(c.ids) > 0 {
		id := c.ids[0]
		c.ids = c.ids[1:]
		return id, true
	}
	return cands[0], true
}

// countingChooser takes the first candidate and counts how many times it is asked,
// so a test can assert a vacuous choice prompts for nothing.
type countingChooser struct{ calls int }

func (c *countingChooser) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	c.calls++
	return cands[0], true
}

// panicOnOptionChooser fails the test if any labeled option is offered, so a test
// can prove a vacuous prompt is never presented. It takes the first creature
// candidate for any non-option choice.
type panicOnOptionChooser struct{}

func (panicOnOptionChooser) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	if len(cands) == 0 {
		return 0, false
	}
	return cands[0], true
}

func (panicOnOptionChooser) ChooseOption(_, _ string, _ []string) int {
	panic("no option should be offered")
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
func attachUpgrade(g *Game, host LocalID, def CardDefinition) LocalID {
	up := g.Register(def, g.owner(host))
	g.AttachUpgrade(host, up)
	return up
}
