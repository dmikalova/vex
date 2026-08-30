package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sequis
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human • Knight
//
//	Reap: Sequis captures 1 Æmber from your opponent.
func TestSequis(t *testing.T) {
	t.Run("captures 1 Æmber when it reaps", func(t *testing.T) {
		var sequis ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(ct.Bind(&sequis, Sequis))},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Reap(sequis)

		h.Expect(sequis).AmberOn(1)
		h.P2.ExpectAmber(2)
	})
}
