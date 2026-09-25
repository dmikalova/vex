package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Dark Faerie
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Mutant
//
//	Skirmish.
//	Fight: Gain 2 Æmber.
func TestDarkFaerie(t *testing.T) {
	t.Run("gains 2 aember when it fights", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(DarkFaerie),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(1))),
			)},
		})

		h.P1.Fight(DarkFaerie, enemy)

		h.P1.ExpectAmber(2)
	})
}
