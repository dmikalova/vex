package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Special Delivery
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	Versatile.
//	Action: Deal 3 damage to a flank creature. If this damage destroys that creature, purge it.
func TestSpecialDelivery(t *testing.T) {
	t.Run("purges a flank creature its damage destroys", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(SpecialDelivery)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
			},
		})

		h.P1.UseAction(SpecialDelivery)

		h.Expect(foe).At(ct.Purge)
	})
}
