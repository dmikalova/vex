package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Interdimensional Graft
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If an opponent forges a key on their next turn, they must give you their remaining Æmber.
func TestInterdimensionalGraft(t *testing.T) {
	t.Run("gives the opponent's remaining Æmber after they forge on their next turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(InterdimensionalGraft)},
			P2: ct.Side{House: card.House.Brobnar, Amber: 10},
		})

		h.P1.Play(InterdimensionalGraft)
		h.P1.EndTurn()

		h.P1.ExpectAmber(5) // 1 bonus Æmber plus the 4 remaining after P2 forges
		h.P2.ExpectAmber(0)
		h.P2.ExpectKeys(1)
	})

	t.Run("does not transfer if the opponent does not forge on their next turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(InterdimensionalGraft)},
			P2: ct.Side{House: card.House.Brobnar, Amber: 5},
		})

		h.P1.Play(InterdimensionalGraft)
		h.P1.EndTurn()

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(5)
		h.P2.ExpectKeys(0)
	})

	t.Run("expires after the opponent's immediate next turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(InterdimensionalGraft)},
			P2: ct.Side{House: card.House.Brobnar, Amber: 5},
		})

		h.P1.Play(InterdimensionalGraft)
		h.P1.EndTurn()
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Logos)
		h.Game().State.Aember[1] = 8
		h.P1.EndTurn()

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
		h.P2.ExpectKeys(1)
	})
}
