package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Francis the "Economist"
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Elf • Thief
//
//	Skirmish.
//	Fight: Each player gains 1 Æmber.
func TestFrancisTheEconomist(t *testing.T) {
	t.Run("each player gains 1 Æmber when it fights", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(FrancisTheEconomist),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(2))),
				),
			},
		})

		h.P1.Fight(FrancisTheEconomist, enemy)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})
}
