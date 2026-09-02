package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gatekeeper
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Armor:  1
//	Traits: Knight • Spirit
//
//	Play: If your opponent has 7 Æmber or more, Gatekeeper captures all but 5 Æmber from your opponent.
func TestGatekeeper(t *testing.T) {
	t.Run("captures the Æmber above five when the opponent has seven or more", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(Gatekeeper)},
			P2: ct.Side{Amber: 9},
		})

		h.P1.Play(Gatekeeper)

		h.P2.ExpectAmber(5)
		h.Expect(Gatekeeper).AmberOn(4)
	})

	t.Run("captures nothing when the opponent has fewer than seven", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(Gatekeeper)},
			P2: ct.Side{Amber: 6},
		})

		h.P1.Play(Gatekeeper)

		h.P2.ExpectAmber(6)
		h.Expect(Gatekeeper).AmberOn(0)
	})
}
