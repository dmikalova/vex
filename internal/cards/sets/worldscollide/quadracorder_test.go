package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Quadracorder
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains, "Your opponent's keys cost +1 Æmber for each house represented among friendly creatures (to a maximum of 3)."
func TestQuadracorder(t *testing.T) {
	t.Run(
		"charges the opponent 1 more per house among friendly creatures, capped at 3",
		func(t *testing.T) {
			var mars, logos, shadows, brobnar ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.StarAlliance,
					InPlay: ct.Cards(
						ct.Upgraded(
							ct.Bind(&mars, ct.Creature(ct.OfHouse(card.House.Mars))),
							Quadracorder,
						),
						ct.Bind(&logos, ct.Creature(ct.OfHouse(card.House.Logos))),
						ct.Bind(&shadows, ct.Creature(ct.OfHouse(card.House.Shadows))),
						ct.Bind(&brobnar, ct.Creature(ct.OfHouse(card.House.Brobnar))),
					),
				},
			})

			// Mars, Logos, Shadows, Brobnar = 4 houses, capped at 3, so +3.
			if got := h.Game().CurrentKeyCost(1); got != 9 {
				t.Errorf("opponent key cost with four houses = %d, want 9", got)
			}
			if got := h.Game().CurrentKeyCost(0); got != 6 {
				t.Errorf("the controller's own key cost = %d, want 6", got)
			}
		},
	)

	t.Run("scales below the cap with fewer houses", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Creature(ct.OfHouse(card.House.Mars)),
						Quadracorder,
					),
					ct.Creature(ct.OfHouse(card.House.Logos)),
				),
			},
		})

		// Mars, Logos = 2 houses, so +2.
		if got := h.Game().CurrentKeyCost(1); got != 8 {
			t.Errorf("opponent key cost with two houses = %d, want 8", got)
		}
	})
}
