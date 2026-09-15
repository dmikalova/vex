package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Harbinger of Doom
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Armor:  3
//	Traits: Demon
//
//	Destroyed: Destroy each creature.
func TestHarbingerOfDoom(t *testing.T) {
	t.Run("wipes the whole board when it is destroyed", func(t *testing.T) {
		var attacker, ally, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(5))),
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					HarbingerOfDoom,
					ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
				),
			},
		})

		h.P1.Fight(attacker, HarbingerOfDoom)

		h.Expect(HarbingerOfDoom).At(ct.Discard)
		h.Expect(attacker).At(ct.Discard)
		h.Expect(ally).At(ct.Discard)
		h.Expect(other).At(ct.Discard)
	})
}
