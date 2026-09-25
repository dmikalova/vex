package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Scaethe
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Mutant
//
//	Destroyed: Destroy the least powerful enemy creature.
func TestScaethe(t *testing.T) {
	t.Run("destroys the least powerful enemy creature when destroyed", func(t *testing.T) {
		var scaethe, bigFoe, weakFoe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&scaethe, Scaethe)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&bigFoe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6))),
					ct.Bind(&weakFoe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Fight(scaethe, bigFoe)

		h.Expect(scaethe).At(ct.Discard)
		h.Expect(weakFoe).At(ct.Discard)
		h.Expect(bigFoe).At(ct.PlayArea)
	})
}
