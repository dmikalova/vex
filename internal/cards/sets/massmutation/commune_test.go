package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Commune
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Omega.
//	Play: Lose all your Æmber. Gain 4 Æmber.
func TestCommune(t *testing.T) {
	t.Run("loses all Æmber, then gains 4", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(Commune),
				Amber: 7,
			},
		})

		h.P1.Play(Commune)

		h.P1.ExpectAmber(4)
	})

	t.Run("gains 4 from an empty pool", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(Commune),
			},
		})

		h.P1.Play(Commune)

		h.P1.ExpectAmber(4)
	})
}
