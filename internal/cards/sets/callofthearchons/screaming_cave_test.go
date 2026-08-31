package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Screaming Cave
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Shuffle your hand and discard pile into your deck.
func TestScreamingCave(t *testing.T) {
	t.Run("shuffles hand and discard pile into the deck", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Dis,
				InPlay:  ct.Cards(ScreamingCave),
				Hand:    ct.Cards(ct.Creature(ct.OfHouse(card.House.Dis)), ct.Creature(ct.OfHouse(card.House.Dis))),
				Discard: ct.Cards(ct.Creature(ct.OfHouse(card.House.Dis))),
			},
		})

		h.P1.UseAction(ScreamingCave)

		g := h.Game()
		if g.State.Hand[0].Count != 0 {
			t.Errorf("hand should be empty, count = %d", g.State.Hand[0].Count)
		}
		if g.State.Discard[0].Count != 0 {
			t.Errorf("discard should be empty, count = %d", g.State.Discard[0].Count)
		}
		if g.State.Deck[0].Count != 3 {
			t.Errorf("deck should hold 3 cards, count = %d", g.State.Deck[0].Count)
		}
	})
}
