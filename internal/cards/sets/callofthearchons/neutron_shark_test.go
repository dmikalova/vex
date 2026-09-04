package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

func TestNeutronShark(t *testing.T) {
	t.Run("stops once the discarded card is a Logos card", func(t *testing.T) {
		var friend, foe, logos ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(NeutronShark),
				InPlay: ct.Cards(
					ct.Bind(&friend, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
				Deck: ct.Cards(ct.Bind(&logos, Batdrone)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			)},
		})

		// The enemy side has a single creature, so that destruction resolves without
		// a prompt; the friendly side offers Neutron Shark itself alongside its
		// neighbour. The discard then turns up a Logos card, ending the loop.
		h.P1.Play(NeutronShark)
		h.P1.ClickCard(friend)

		h.Expect(foe).At(ct.Discard)
		h.Expect(friend).At(ct.Discard)
		h.Expect(logos).At(ct.Discard)
	})
}
