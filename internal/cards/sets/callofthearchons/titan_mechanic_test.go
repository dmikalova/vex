package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Titan Mechanic
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Cyborg • Scientist
//
//	While Titan Mechanic is on a flank, each player's keys cost -1 Æmber.
func TestTitanMechanic(t *testing.T) {
	t.Run("discounts both players' keys while on a flank", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(TitanMechanic),
			},
		})

		if got := h.Game().CurrentKeyCost(0); got != 5 {
			t.Errorf("controller key cost = %d, want 5", got)
		}
		if got := h.Game().CurrentKeyCost(1); got != 5 {
			t.Errorf("opponent key cost = %d, want 5", got)
		}
	})

	t.Run("does nothing while flanked on both sides", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3)),
					TitanMechanic,
					ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3)),
				),
			},
		})

		if got := h.Game().CurrentKeyCost(0); got != 6 {
			t.Errorf("key cost with Titan Mechanic off the flanks = %d, want 6", got)
		}
	})
}
