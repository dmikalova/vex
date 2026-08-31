package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Staunch Knight
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Staunch Knight gains +2 power while it is on a flank.
func TestStaunchKnight(t *testing.T) {
	t.Run("gains +2 power while on a flank", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(StaunchKnight, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
			},
		})

		h.Expect(StaunchKnight).Power(6)
	})

	t.Run("no bonus while in the middle of the battleline", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3)),
					StaunchKnight,
					ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3)),
				),
			},
		})

		h.Expect(StaunchKnight).Power(4)
	})
}
