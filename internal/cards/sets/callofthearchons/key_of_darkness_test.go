package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Key of Darkness
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: If your opponent has no Æmber, forge a key at +2 Æmber current cost. Otherwise, forge a key at +6 Æmber current cost.
func TestKeyOfDarkness(t *testing.T) {
	t.Run("forges at +6 while the opponent holds Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(KeyOfDarkness),
				Amber: 12,
			},
			P2: ct.Side{Amber: 1},
		})

		h.P1.Play(KeyOfDarkness)
		h.P1.ExpectKeys(1)
		h.P1.ExpectAmber(0)
	})

	t.Run("forges at +2 while the opponent has no Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(KeyOfDarkness),
				Amber: 12,
			},
		})

		h.P1.Play(KeyOfDarkness)
		h.P1.ExpectKeys(1)
		h.P1.ExpectAmber(4)
	})

	t.Run("forges nothing it cannot pay for", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(KeyOfDarkness),
				Amber: 7,
			},
			P2: ct.Side{Amber: 1},
		})

		h.P1.Play(KeyOfDarkness)
		h.P1.ExpectKeys(0)
		h.P1.ExpectAmber(7)
	})
}
