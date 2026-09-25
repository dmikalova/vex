package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Call to Action
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Ready each friendly Knight creature.
func TestCallToAction(t *testing.T) {
	t.Run("readies each friendly Knight creature", func(t *testing.T) {
		var knight, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(CallToAction),
				InPlay: ct.Cards(
					ct.Bind(&knight, ct.Creature(ct.Traits(card.Traits.Knight))),
					ct.Bind(&other, ct.Creature(ct.Traits(card.Traits.Beast))),
				),
			},
		})
		knight.Exhaust()
		other.Exhaust()

		h.P1.Play(CallToAction)

		h.Expect(knight).Ready()
		h.Expect(other).Exhausted()
	})
}
