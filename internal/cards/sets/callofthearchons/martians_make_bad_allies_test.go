package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

func TestMartiansMakeBadAllies(t *testing.T) {
	t.Run("purges each non-Mars creature in hand and gains 1 Æmber per purge", func(t *testing.T) {
		var martian, brobnar, untamed, tactic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand: ct.Cards(
					MartiansMakeBadAllies,
					ct.Bind(&martian, ct.Creature(ct.OfHouse(card.House.Mars))),
					ct.Bind(&brobnar, ct.Creature(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&untamed, ct.Creature(ct.OfHouse(card.House.Untamed))),
					ct.Bind(&tactic, ct.Tactic(ct.OfHouse(card.House.Shadows))),
				),
			},
		})

		h.P1.Play(MartiansMakeBadAllies)

		h.Expect(brobnar).At(ct.Purge)
		h.Expect(untamed).At(ct.Purge)
		h.Expect(martian).At(ct.Hand)
		h.Expect(tactic).At(ct.Hand)
		h.P1.ExpectAmber(2)
	})
}
