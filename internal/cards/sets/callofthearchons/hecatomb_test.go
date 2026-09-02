package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hecatomb
func TestHecatomb(t *testing.T) {
	dis := ct.OfHouse(card.House.Dis)

	t.Run("destroys each Dis creature and pays each player for their own dead", func(t *testing.T) {
		var mine1, mine2, mineOther, theirs ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Hecatomb),
				InPlay: ct.Cards(
					ct.Bind(&mine1, ct.Creature(dis)),
					ct.Bind(&mine2, ct.Creature(dis)),
					ct.Bind(&mineOther, ct.Creature(ct.OfHouse(card.House.Mars))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&theirs, ct.Creature(dis))),
			},
		})

		h.P1.Play(Hecatomb)
		h.Expect(mine1).At(ct.Discard)
		h.Expect(mine2).At(ct.Discard)
		h.Expect(theirs).At(ct.Discard)
		h.Expect(mineOther).At(ct.PlayArea)

		// Two dead of their own plus the card's own 1 Æmber bonus; the opponent
		// collects only for the one creature they controlled.
		h.P1.ExpectAmber(3)
		h.P2.ExpectAmber(1)
	})
}
