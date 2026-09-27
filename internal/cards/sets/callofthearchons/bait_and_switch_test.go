package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Bait and Switch
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Steal 1 Æmber -> if your opponent has more Æmber than you, repeat this effect.
func TestBaitAndSwitch(t *testing.T) {
	t.Run("steals 1 Æmber at a time while the opponent still leads", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(BaitAndSwitch),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(BaitAndSwitch)

		// 5/0 -> 4/1 -> 3/2 -> 2/3 (opponent no longer leads).
		h.P1.ExpectAmber(3)
		h.P2.ExpectAmber(2)
	})

	t.Run("steals at most 6 Æmber, the Rule of 6 bound on repeats", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(BaitAndSwitch),
			},
			P2: ct.Side{Amber: 20},
		})

		h.P1.Play(BaitAndSwitch)

		// The play steals once and repeats five more times, not until the lead is gone.
		h.P1.ExpectAmber(6)
		h.P2.ExpectAmber(14)
	})

	t.Run("a protected pool stops the loop without spending the Rule of 6", func(t *testing.T) {
		// The arrow is a result gate: a steal that moves nothing ends the loop, so
		// the play spends only the one usage it records for itself and a second copy
		// of the name is still playable this turn.
		var first, second ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand: ct.Cards(
					ct.Bind(&first, BaitAndSwitch),
					ct.Bind(&second, BaitAndSwitch),
				),
			},
			P2: ct.Side{
				Amber:  5,
				InPlay: ct.Cards(TheVaultkeeper),
			},
		})

		h.P1.Play(first)

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(5)

		h.P1.Play(second)

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(5)
	})
}
