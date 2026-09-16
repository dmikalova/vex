package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Titan Engineer
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Cyborg • Scientist
//
//	While Titan Engineer is not on a flank, each player's keys cost +1 Æmber.
func TestTitanEngineer(t *testing.T) {
	t.Run("does nothing while on a flank", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(TitanEngineer),
			},
		})

		if got := h.Game().CurrentKeyCost(0); got != 6 {
			t.Errorf("controller key cost = %d, want 6", got)
		}
		if got := h.Game().CurrentKeyCost(1); got != 6 {
			t.Errorf("opponent key cost = %d, want 6", got)
		}
	})

	t.Run("raises both players' keys while off the flanks", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3)),
					TitanEngineer,
					ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3)),
				),
			},
		})

		if got := h.Game().CurrentKeyCost(0); got != 7 {
			t.Errorf("controller key cost off the flanks = %d, want 7", got)
		}
		if got := h.Game().CurrentKeyCost(1); got != 7 {
			t.Errorf("opponent key cost off the flanks = %d, want 7", got)
		}
	})
}
