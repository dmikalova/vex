package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mother Northelle
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Human • Monk
//
//	Elusive.
//	Reap: Move 1 Æmber from a friendly creature to your pool.
func TestMotherNorthelle(t *testing.T) {
	t.Run("moves 1 aember from a friendly creature to your pool when reaping", func(t *testing.T) {
		var northelle, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Bind(&northelle, MotherNorthelle),
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
				),
			},
			P2: ct.Side{},
		})
		h.Game().State.Cards[ally.ID()].Amber = 1

		h.P1.Reap(northelle)

		h.P1.ExpectAmber(2)
		h.Expect(ally).AmberOn(0)
	})
}
