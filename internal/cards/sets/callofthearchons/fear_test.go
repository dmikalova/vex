package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Fear
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Put an enemy creature into its owner's hand.
func TestFear(t *testing.T) {
	t.Run("puts a chosen enemy creature into its owner's hand", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(Fear)},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3)))),
			},
		})

		h.P1.Play(Fear)

		h.Expect(enemy).At(ct.Hand)
	})
}
