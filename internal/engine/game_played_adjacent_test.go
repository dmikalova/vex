package engine

import "testing"

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
