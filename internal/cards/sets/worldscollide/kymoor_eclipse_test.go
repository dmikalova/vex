package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Kymoor Eclipse
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Shuffle each flank Creature into its owner's deck.
func TestKymoorEclipse(t *testing.T) {
	t.Run("shuffles each flank creature into its owner's deck", func(t *testing.T) {
		var left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(KymoorEclipse),
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(3))),
					ct.Bind(&mid, ct.Creature(ct.Power(3))),
					ct.Bind(&right, ct.Creature(ct.Power(3))),
				),
			},
		})

		h.P1.Play(KymoorEclipse)

		h.Expect(left).At(ct.Deck)
		h.Expect(right).At(ct.Deck)
		h.Expect(mid).At(ct.PlayArea)
	})
}
