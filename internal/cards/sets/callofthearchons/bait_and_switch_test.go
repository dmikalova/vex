package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bait and Switch
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: If your opponent has more Æmber than you, steal 1 Æmber -> repeat this effect.
func TestBaitAndSwitch(t *testing.T) {
	t.Run("steals 1 Æmber at a time while the opponent still leads", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(BaitAndSwitch)},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(BaitAndSwitch)

		// 5/0 -> 4/1 -> 3/2 -> 2/3 (opponent no longer leads).
		h.P1.ExpectAmber(3)
		h.P2.ExpectAmber(2)
	})
}
