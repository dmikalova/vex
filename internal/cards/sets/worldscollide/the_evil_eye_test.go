package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Evil Eye
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Keys cost +3 Æmber during your opponent's next turn.
func TestTheEvilEye(t *testing.T) {
	t.Run("taxes the opponent's next turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(TheEvilEye),
			},
			P2: ct.Side{
				House: card.House.Dis,
				Amber: 8,
			},
		})

		h.P1.Play(TheEvilEye)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Dis)

		// 8 Æmber does not cover a key at 6 + 3.
		h.P2.ExpectKeys(0)
		h.P2.ExpectAmber(8)
	})
}
