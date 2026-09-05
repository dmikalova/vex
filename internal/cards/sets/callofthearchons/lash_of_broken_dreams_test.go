package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lash of Broken Dreams
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Common
//	Traits: Weapon
//
//	Action: Keys cost +3 Æmber during your opponent's next turn.
func TestLashOfBrokenDreams(t *testing.T) {
	t.Run("taxes the opponent's next turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(LashOfBrokenDreams),
			},
			P2: ct.Side{House: card.House.Dis, Amber: 8},
		})

		h.P1.UseAction(LashOfBrokenDreams)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Dis)

		// 8 Æmber does not cover a key at 6 + 3.
		h.P2.ExpectKeys(0)
		h.P2.ExpectAmber(8)
	})

	t.Run("the tax lifts after that turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(LashOfBrokenDreams),
			},
			P2: ct.Side{House: card.House.Dis, Amber: 6},
		})

		h.P1.UseAction(LashOfBrokenDreams)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Dis)
		h.P2.ExpectKeys(0)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Dis)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Dis)
		h.P2.ExpectKeys(1)
	})
}
