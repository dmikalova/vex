package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Commander Remiel
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Knight
//
//	Reap: Use a friendly non-Sanctum creature.
func TestCommanderRemiel(t *testing.T) {
	t.Run("reaps, then uses a friendly non-Sanctum creature to reap again", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					CommanderRemiel,
					ct.Creature(
						ct.OfHouse(card.House.Mars),
						ct.Power(3),
					), // the only non-Sanctum friendly
				),
			},
		})

		h.P1.Reap(CommanderRemiel)

		h.P1.ExpectAmber(2) // Remiel's reap + the used ally's reap
	})
}
