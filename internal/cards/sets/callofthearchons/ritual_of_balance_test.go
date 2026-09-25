package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Ritual of Balance
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Power
//
//	Action: If your opponent has 6 Æmber or more, steal 1 Æmber.
func TestRitualOfBalance(t *testing.T) {
	t.Run("does nothing while the opponent is below 6 Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(RitualOfBalance),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.UseAction(RitualOfBalance)

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(5)
	})

	t.Run("steals 1 Æmber once the opponent has 6 or more", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(RitualOfBalance),
			},
			P2: ct.Side{Amber: 6},
		})

		h.P1.UseAction(RitualOfBalance)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(5)
	})
}
