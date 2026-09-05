package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Drummernaut
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Giant
//
//	Play/Fight/Reap: Put another friendly Giant trait creature into its owner's hand.
func TestDrummernaut(t *testing.T) {
	t.Run("returns another friendly Giant to hand when played", func(t *testing.T) {
		var giant ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(Drummernaut),
				InPlay: ct.Cards(
					ct.Bind(&giant, ct.Creature(ct.Traits(card.Traits.Giant), ct.Power(6))),
				),
			},
		})

		h.P1.Play(Drummernaut)

		h.Expect(giant).At(ct.Hand)
	})
}
