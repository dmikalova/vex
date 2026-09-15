package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Haedroth's Wall
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Location
//
//	Each friendly flank creature gains +2 power.
func TestHaedrothsWall(t *testing.T) {
	t.Run("gives friendly flank creatures +2 power while in play", func(t *testing.T) {
		var left, middle, right, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					HaedrothsWall,
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
					ct.Bind(&middle, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
				),
			},
		})

		h.Expect(left).Power(6)   // flank creature buffed
		h.Expect(right).Power(6)  // flank creature buffed
		h.Expect(middle).Power(4) // not on a flank
		h.Expect(enemy).Power(4)  // buffs only friendly creatures
	})
}
