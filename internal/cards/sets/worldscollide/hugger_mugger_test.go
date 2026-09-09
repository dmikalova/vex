package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hugger-Mugger
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Elf • Thief
//
//	Elusive.
//	Play: Hugger-Mugger captures 1 Æmber from your opponent. If your opponent has more forged keys than you, steal 1 Æmber.
func TestHuggerMugger(t *testing.T) {
	t.Run("captures, then steals when behind on keys", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(HuggerMugger),
			},
			P2: ct.Side{Amber: 3},
		})
		h.Game().State.Keys[1] = 1

		h.P1.Play(HuggerMugger)

		h.Expect(HuggerMugger).AmberOn(1)
		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})

	t.Run("only captures when not behind on keys", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(HuggerMugger),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(HuggerMugger)

		h.Expect(HuggerMugger).AmberOn(1)
		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(2)
	})
}
