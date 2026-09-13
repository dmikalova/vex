package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Grump Buggy
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Vehicle
//
//	Your opponent's keys cost +1 Æmber for each friendly Creature with power 5 or higher.
//	Your keys cost +1 Æmber for each enemy Creature with power 5 or higher.
func TestGrumpBuggy(t *testing.T) {
	t.Run(
		"raises each player's key cost per power-5 creature the other side threatens",
		func(t *testing.T) {
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Brobnar,
					InPlay: ct.Cards(
						GrumpBuggy,
						ct.Creature(ct.Power(5)),
						ct.Creature(ct.Power(4)),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(ct.Creature(ct.Power(6)), ct.Creature(ct.Power(7))),
				},
			})

			base := 6
			if got := h.Game().CurrentKeyCost(1); got != base+1 {
				t.Errorf("opponent key cost = %d, want %d", got, base+1)
			}
			if got := h.Game().CurrentKeyCost(0); got != base+2 {
				t.Errorf("controller key cost = %d, want %d", got, base+2)
			}
		},
	)
}
