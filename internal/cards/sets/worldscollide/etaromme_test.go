package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Etaromme
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Demon
//
//	Reap: Destroy a creature of the house with the most creatures in play.
func TestEtaromme(t *testing.T) {
	t.Run("destroys a chosen creature of the most populous house", func(t *testing.T) {
		var etaromme, brob, brob2, dis ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&etaromme, Etaromme)),
			},
			P2: ct.Side{
				// Brobnar leads with two creatures; Dis (Etaromme) has just one.
				InPlay: ct.Cards(
					ct.Bind(&brob, ct.Creature(ct.Power(3), ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&brob2, ct.Creature(ct.Power(3), ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&dis, ct.Creature(ct.Power(3), ct.OfHouse(card.House.Dis))),
				),
			},
		})

		h.P1.Reap(etaromme)
		h.P1.ClickCard(brob)

		h.Expect(brob).At(ct.Discard)
		h.Expect(brob2).At(ct.PlayArea)
		h.Expect(dis).At(ct.PlayArea)
		h.Expect(etaromme).At(ct.PlayArea)
	})
}
