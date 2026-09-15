package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Earthshaker
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  7
//	Traits: Giant
//
//	Play: Destroy each creature with power 3 or lower.
func TestEarthshaker(t *testing.T) {
	t.Run("destroys each creature with power 3 or lower on Play", func(t *testing.T) {
		var weak, strong ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(Earthshaker),
				InPlay: ct.Cards(
					ct.Bind(&weak, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(2))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&strong, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				),
			},
		})

		h.P1.Play(Earthshaker)

		h.Expect(weak).At(ct.Discard)         // power <= 3 destroyed
		h.Expect(strong).At(ct.PlayArea)      // a stronger creature survives
		h.Expect(Earthshaker).At(ct.PlayArea) // 7 power, survives its own Play
	})
}
