package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Disruption Field
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Your opponent's keys cost +1 Æmber for each disruption counter on Disruption Field.
//	This creature gains, "Fight/Reap: Put a disruption counter on Disruption Field."
func TestDisruptionField(t *testing.T) {
	t.Run(
		"host reaping places a disruption counter that raises the opponent's key cost",
		func(t *testing.T) {
			var field, host ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.StarAlliance,
					InPlay: ct.Cards(
						ct.Bind(
							&host,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
					),
					Hand: ct.Cards(ct.Bind(&field, DisruptionField)),
				},
			})

			h.P1.Play(field) // the lone host auto-attaches

			if got := h.Game().CurrentKeyCost(1); got != 6 {
				t.Errorf("opponent key cost with no counters = %d, want 6", got)
			}

			h.P1.Reap(host)

			if got := h.Game().CountersOn(field.ID(), card.Counter.Disruption); got != 1 {
				t.Fatalf("disruption counters after reap = %d, want 1", got)
			}
			if got := h.Game().CurrentKeyCost(1); got != 7 {
				t.Errorf("opponent key cost with one counter = %d, want 7", got)
			}
		},
	)
}
