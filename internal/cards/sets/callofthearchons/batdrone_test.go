package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Batdrone
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Robot
//
//	Skirmish.
//	Fight: Steal 1 Æmber.
func TestBatdrone(t *testing.T) {
	t.Run("steals 1 Æmber when it fights and Skirmish spares it return damage", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, InPlay: ct.Cards(Batdrone)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
				Amber: 2,
			},
		})

		h.Expect(Batdrone).Power(2)

		h.P1.Fight(Batdrone, enemy)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
		h.Expect(Batdrone).Damage(0) // Skirmish: no return damage
	})
}
