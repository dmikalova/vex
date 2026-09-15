package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bulwark
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Each neighboring creature gains +2 armor.
func TestBulwark(t *testing.T) {
	t.Run("gives each battleline neighbor +2 armor but not distant creatures", func(t *testing.T) {
		var left, right, far ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
					Bulwark,
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
					ct.Bind(&far, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
				),
			},
		})

		h.Expect(left).Armor(2)
		h.Expect(right).Armor(2)
		h.Expect(far).Armor(0) // not a neighbor
	})
}
