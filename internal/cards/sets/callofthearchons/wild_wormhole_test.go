package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Wild Wormhole
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Play the top card of your deck.
func TestWildWormhole(t *testing.T) {
	t.Run("plays the top card of the deck", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(WildWormhole),
				Deck: ct.Cards(
					ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
		})

		h.P1.Play(WildWormhole)

		h.Expect(top).At(ct.PlayArea)
	})
}
