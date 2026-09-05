package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Yurk
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Demon
//
//	Play: Discard a card from your hand.
func TestYurk(t *testing.T) {
	t.Run("discards a card from hand", func(t *testing.T) {
		var fodder ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand: ct.Cards(
					Yurk,
					ct.Bind(&fodder, ct.Creature(ct.Power(3))),
				),
			},
		})

		h.P1.Play(Yurk)

		h.Expect(fodder).At(ct.Discard)
	})
}
