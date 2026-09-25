package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Virtuous Works
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber Æmber Æmber
func TestVirtuousWorks(t *testing.T) {
	t.Run("gains 3 Æmber when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(VirtuousWorks),
			},
		})

		h.P1.Play(VirtuousWorks)

		h.P1.ExpectAmber(3)
	})
}
