package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

func TestLibraryAccess(t *testing.T) {
	t.Run("draws a card each time another card is played this turn", func(t *testing.T) {
		var c1, c2, d1, d2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					LibraryAccess,
					ct.Bind(&c1, ct.Creature(ct.OfHouse(card.House.Logos))),
					ct.Bind(&c2, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
				Deck: ct.Cards(
					ct.Bind(&d1, ct.Creature(ct.OfHouse(card.House.Logos))),
					ct.Bind(&d2, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		// Library Access itself is not "another card", so it draws nothing.
		h.P1.Play(LibraryAccess)
		h.Expect(d1).At(ct.Deck)
		h.Expect(d2).At(ct.Deck)

		h.P1.Play(c1)
		h.Expect(d1).At(ct.Hand)

		h.P1.Play(c2)
		h.Expect(d2).At(ct.Hand)
	})
}
