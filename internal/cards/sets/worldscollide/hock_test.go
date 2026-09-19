package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hock
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Destroy an artifact -> gain 1 Æmber.
func TestHock(t *testing.T) {
	t.Run("destroys an artifact and gains 1 Æmber", func(t *testing.T) {
		var relic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(Hock),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&relic, ct.Artifact()))},
		})

		h.P1.Play(Hock)

		h.Expect(relic).At(ct.Discard)
		h.P1.ExpectAmber(2) // 1 from Hock's Æmber bonus, 1 from the destroy
	})
}
