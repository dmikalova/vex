package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lateral Shift
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Special
//
//	Play: Play a card from your opponent's hand.
func TestLateralShift(t *testing.T) {
	t.Run("plays a card out of the opponent's hand as your own", func(t *testing.T) {
		var borrowed ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(LateralShift),
			},
			P2: ct.Side{
				Hand: ct.Cards(
					ct.Bind(&borrowed, ct.Creature(ct.Power(4))),
				),
			},
		})

		h.P1.Play(LateralShift)

		h.Expect(borrowed).At(ct.PlayArea)
	})
}
