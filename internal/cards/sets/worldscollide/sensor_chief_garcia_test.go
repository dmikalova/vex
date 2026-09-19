package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sensor Chief Garcia
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human
//
//	Play/Fight/Reap: Keys cost +2 Æmber during your opponent's next turn.
func TestSensorChiefGarcia(t *testing.T) {
	t.Run("raises the opponent's key cost during their next turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(SensorChiefGarcia),
			},
			P2: ct.Side{},
		})

		h.P1.Play(SensorChiefGarcia)
		h.P1.EndTurn()

		if got := h.Game().CurrentKeyCost(1); got != 8 {
			t.Errorf("opponent key cost = %d, want 8", got)
		}
	})
}
