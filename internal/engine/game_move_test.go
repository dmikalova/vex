package engine

import "testing"

func TestMoveCardRefusesImpossibleMoves(t *testing.T) {
	t.Run("refuses a zone that holds no pile", func(t *testing.T) {
		g := started(t)
		id := g.Register(testCreature("held", 1), 0)
		g.State.Hand[0].add(id)
		before := g.State.Hand[0].Count

		if g.moveCard(id,
			zoneRef{Player: 0, Zone: Hand},
			zoneRef{Player: 0, Zone: inPlay},
			CardMoved{Player: 0, Card: id, From: Hand, To: inPlay}) {
			t.Fatal("moveCard into play reported a move")
		}
		if g.State.Hand[0].Count != before {
			t.Errorf("hand = %d cards, want %d", g.State.Hand[0].Count, before)
		}
	})

	t.Run("refuses a card the source zone does not hold", func(t *testing.T) {
		g := started(t)
		id := g.Register(testCreature("held", 1), 0)
		g.State.Hand[0].add(id)

		if g.moveCard(id,
			zoneRef{Player: 0, Zone: Discard},
			zoneRef{Player: 0, Zone: purged},
			CardMoved{Player: 0, Card: id, From: Discard, To: purged}) {
			t.Fatal("moveCard reported moving a card that was not there")
		}
		if g.State.Purge[0].Count != 0 {
			t.Errorf("purge = %d cards, want 0", g.State.Purge[0].Count)
		}
	})
}
