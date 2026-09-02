package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Key Abduction
func TestKeyAbduction(t *testing.T) {
	mars := ct.OfHouse(card.House.Mars)

	t.Run("returns each Mars creature and forges at the reduced cost", func(t *testing.T) {
		var mine, theirs, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				Hand:   ct.Cards(KeyAbduction),
				InPlay: ct.Cards(ct.Bind(&mine, ct.Creature(mars))),
				Amber:  14,
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&theirs, ct.Creature(mars)),
					ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Dis))),
				),
			},
		})

		h.P1.Play(KeyAbduction)
		h.Expect(mine).At(ct.Hand)
		h.Expect(theirs).At(ct.Hand)
		h.Expect(other).At(ct.PlayArea)

		// One returned card in hand cuts the +9 to +8, so the key costs 14 out of the
		// 15 Æmber the pool holds once the card's own bonus lands.
		h.P1.ExpectKeys(1)
		h.P1.ExpectAmber(1)
	})
}
