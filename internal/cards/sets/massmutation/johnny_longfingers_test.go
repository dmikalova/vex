package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Johnny Longfingers
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Mutant • Thief
//
//	Each friendly Mutant creature gains, "Destroyed: Steal 1 Æmber."
func TestJohnnyLongfingers(t *testing.T) {
	t.Run("friendly Mutant steals 1 Æmber when destroyed", func(t *testing.T) {
		var mutant, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					JohnnyLongfingers,
					ct.Bind(&mutant, ct.Creature(
						ct.OfHouse(card.House.Shadows),
						ct.Traits(card.Traits.Mutant),
						ct.Power(2),
					)),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				),
				Amber: 2,
			},
		})

		h.P1.Fight(mutant, enemy)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})

	t.Run("friendly non-Mutant does not steal when destroyed", func(t *testing.T) {
		var beast, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					JohnnyLongfingers,
					ct.Bind(&beast, ct.Creature(
						ct.OfHouse(card.House.Shadows),
						ct.Traits(card.Traits.Beast),
						ct.Power(2),
					)),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				),
				Amber: 2,
			},
		})

		h.P1.Fight(beast, enemy)

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(2)
	})
}
