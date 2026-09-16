package engine

import "testing"

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
