package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Xanthyx Harvester
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	You cannot use this card unless it has no non-Mars neighbor.
//	Reap: Gain 1 Æmber.
func TestXanthyxHarvester(t *testing.T) {
	t.Run("cannot be used while it has a non-Mars neighbor", func(t *testing.T) {
		var xan ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Bind(&xan, XanthyxHarvester),
					ct.Creature(ct.OfHouse(card.House.Sanctum)),
				),
			},
		})

		h.P1.ExpectCannotUse(xan)
	})

	t.Run("reaps for 1 aember when all neighbors are Mars", func(t *testing.T) {
		var xan ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Bind(&xan, XanthyxHarvester),
					ct.Creature(ct.OfHouse(card.House.Mars)),
				),
			},
		})

		h.P1.Reap(xan)

		// 1 Æmber from reaping, plus 1 from Xanthyx's ability.
		h.P1.ExpectAmber(2)
	})
}
