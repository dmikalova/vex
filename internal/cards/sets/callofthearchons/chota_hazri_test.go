package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Chota Hazri
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Human • Witch
//
//	Play: Lose 1 Æmber. Forge a key at current cost -> purge Chota Hazri.
func TestChotaHazri(t *testing.T) {
	t.Run("loses 1 Æmber to forge a key at current cost", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(ChotaHazri),
				Amber: engine.KeyCost + 1, // enough to pay the 1 and the key cost
			},
		})

		h.P1.Play(ChotaHazri)

		h.P1.ExpectKeys(1)  // forged
		h.P1.ExpectAmber(0) // 1 lost + key cost paid
	})
}
