package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lord Invidius
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Demon • Leader
//
//	Elusive.
//	While Lord Invidius is in the center of your battleline, it gains, "Reap: Take control of an enemy flank creature and exhaust it. While under your control, it belongs to house Dis."
func TestLordInvidius(t *testing.T) {
	t.Run(
		"centered, its reap seizes an enemy flank creature as an exhausted Dis creature",
		func(t *testing.T) {
			var foe ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Dis,
					InPlay: ct.Cards(
						ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(2)),
						LordInvidius,
						ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(2)),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
					),
				},
			})

			h.P1.Reap(LordInvidius)
			h.P1.ClickOption("left flank") // where the seized creature lands

			h.Expect(foe).Exhausted()
			if got := h.Game().Controller(foe.ID()); got != 0 {
				t.Errorf("seized creature controller = %d, want P1 (0)", got)
			}
			if got := h.Game().House(foe.ID()); got != card.House.Dis {
				t.Errorf("seized creature house = %v, want Dis", got)
			}
		},
	)

	t.Run("off center, it has no reap ability", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					LordInvidius,
					ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(2)),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(LordInvidius)

		if got := h.Game().Controller(foe.ID()); got != 1 {
			t.Errorf("seized creature controller = %d, want P2 (1) — no reap while off center", got)
		}
	})
}
