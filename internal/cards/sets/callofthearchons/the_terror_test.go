package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Terror
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon • Knight
//
//	Play: If your opponent has exactly 0 Æmber, gain 2 Æmber.
func TestTheTerror(t *testing.T) {
	t.Run("gains 2 Æmber when the opponent has none", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(TheTerror)},
			P2: ct.Side{Amber: 0},
		})

		h.P1.Play(TheTerror)

		h.P1.ExpectAmber(2)
	})

	t.Run("does nothing when the opponent has Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(TheTerror)},
			P2: ct.Side{Amber: 1},
		})

		h.P1.Play(TheTerror)

		h.P1.ExpectAmber(0)
	})
}
