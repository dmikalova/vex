package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Lumilu
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Beast • Cat
//
//	Reap: For each other friendly Beast creature, gain 1 Æmber.
func TestLumilu(t *testing.T) {
	t.Run("gains 1 Æmber for each other friendly Beast creature", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					Lumilu,
					ct.Creature(ct.Traits(card.Traits.Beast)),
					ct.Creature(ct.Power(3)),
				),
			},
		})

		h.P1.Reap(Lumilu)

		// 1 Æmber from reaping, plus 1 for the single other Beast (the
		// non-Beast creature is not counted).
		h.P1.ExpectAmber(2)
	})

	t.Run("counts only Beasts other than itself", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(Lumilu),
			},
		})

		h.P1.Reap(Lumilu)

		// Only the reap Æmber: Lumilu does not count itself.
		h.P1.ExpectAmber(1)
	})
}
