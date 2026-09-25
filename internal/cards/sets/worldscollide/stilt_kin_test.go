package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Stilt-Kin
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Goblin
//
//	Skirmish.
//	After a Giant creature is played adjacent to Stilt-Kin, ready and fight with Stilt-Kin.
func TestStiltKin(t *testing.T) {
	t.Run("a Giant played adjacent readies Stilt-Kin and fights", func(t *testing.T) {
		var stilt, giant, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&stilt, StiltKin)),
				Hand: ct.Cards(
					ct.Bind(&giant, ct.Creature(
						ct.OfHouse(card.House.Brobnar),
						ct.Traits(card.Traits.Giant),
						ct.Power(5),
					)),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20)))),
			},
		})

		stilt.Exhaust() // must be readied before it can fight

		h.P1.Play(giant)

		h.Expect(foe).Damage(2)   // Stilt-Kin was readied and fought for its power
		h.Expect(stilt).Damage(0) // Skirmish: no damage in return
	})

	t.Run("a non-Giant played adjacent does nothing", func(t *testing.T) {
		var stilt, goblin, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&stilt, StiltKin)),
				Hand: ct.Cards(
					ct.Bind(&goblin, ct.Creature(
						ct.OfHouse(card.House.Brobnar),
						ct.Traits(card.Traits.Goblin),
						ct.Power(5),
					)),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20)))),
			},
		})

		stilt.Exhaust()

		h.P1.Play(goblin)

		h.Expect(foe).Damage(0) // the trait gate is not met, so no fight
		h.Expect(stilt).Exhausted()
	})
}
