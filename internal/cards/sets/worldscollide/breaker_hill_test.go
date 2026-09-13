package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Breaker Hill
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	Each neighboring Creature gains, "Action: Steal 1 Æmber."
func TestBreakerHill(t *testing.T) {
	t.Run("a neighbor may use the granted action to steal 1 Æmber", func(t *testing.T) {
		var neighbor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.Shadows))),
					BreakerHill,
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(neighbor)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("both flanking neighbors receive the granted action", func(t *testing.T) {
		var left, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Shadows))),
					BreakerHill,
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Shadows))),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(left)
		h.P1.UseAction(right)

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(1)
	})
}
