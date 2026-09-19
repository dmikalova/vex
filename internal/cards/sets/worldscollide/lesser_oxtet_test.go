package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lesser Oxtet
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Demon
//
//	Elusive.
//	Play: Purge each card from your hand.
//	Reap: Keys cost +3 Æmber during your opponent's next turn.
func TestLesserOxtet(t *testing.T) {
	t.Run("play purges each card in hand", func(t *testing.T) {
		var oxtet, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand: ct.Cards(
					ct.Bind(&oxtet, LesserOxtet),
					ct.Bind(&other, ct.Creature()),
				),
			},
		})

		h.P1.Play(LesserOxtet)

		h.Expect(other).At(ct.Purge)
	})

	t.Run("reap taxes the opponent's next turn", func(t *testing.T) {
		var oxtet ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&oxtet, LesserOxtet)),
			},
			P2: ct.Side{
				House: card.House.Dis,
				Amber: 8,
			},
		})

		h.P1.Reap(oxtet)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Dis)

		// 8 Æmber does not cover a key at 6 + 3.
		h.P2.ExpectKeys(0)
		h.P2.ExpectAmber(8)
	})
}
