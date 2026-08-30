package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lights Out
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Put up to 2 creatures into their owners' hands.
func TestLightsOut(t *testing.T) {
	t.Run("returns up to 2 enemy creatures to hand", func(t *testing.T) {
		var a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(LightsOut)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Mars))),
					ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Mars))),
				),
			},
		})

		h.P1.Play(LightsOut)
		h.P1.ClickOption(a.Name())
		h.P1.ClickOption(b.Name())

		h.Expect(a).At(ct.Hand)
		h.Expect(b).At(ct.Hand)
	})
}
