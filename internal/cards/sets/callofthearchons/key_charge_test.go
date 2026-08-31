package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Key Charge
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Lose 1 Æmber -> forge a key at current cost.
func TestKeyCharge(t *testing.T) {
	t.Run("loses 1 Æmber then forges a key at current cost", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(KeyCharge), Amber: 7},
		})

		h.P1.Play(KeyCharge)

		h.P1.ExpectKeys(1)
		h.P1.ExpectAmber(0)
	})

	t.Run("does not forge when there is no Æmber to lose", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(KeyCharge), Amber: 0},
		})

		h.P1.Play(KeyCharge)

		h.P1.ExpectKeys(0)
	})
}
