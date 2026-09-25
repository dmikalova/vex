package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Gluttony
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Traits: Demon • Sin
//
//	Play: For each friendly Sin creature, exalt Gluttony.
//	Reap: Move all Æmber from each friendly creature to your pool.
func TestGluttony(t *testing.T) {
	t.Run("exalts once for each friendly Sin creature on play", func(t *testing.T) {
		var gluttony ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&gluttony, Gluttony)),
				InPlay: ct.Cards(
					ct.Creature(
						ct.OfHouse(card.House.Dis),
						ct.Traits(card.Traits.Sin),
						ct.Power(3),
					),
				),
			},
		})

		h.P1.Play(Gluttony)

		// Two friendly Sin creatures once Gluttony is in play (itself and one other).
		h.Expect(gluttony).AmberOn(2)
	})

	t.Run("moves all Æmber from friendly creatures to your pool on reap", func(t *testing.T) {
		var a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					Gluttony,
					ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
					ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
				),
			},
		})
		h.Game().AddAmberOn(a.ID(), 2)
		h.Game().AddAmberOn(b.ID(), 3)

		h.P1.Reap(Gluttony)

		h.Expect(a).AmberOn(0)
		h.Expect(b).AmberOn(0)
		h.P1.ExpectAmber(6) // 2 + 3 moved, plus 1 from the reap itself
	})
}
