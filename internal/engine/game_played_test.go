package engine

import "testing"

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
