package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Old Yurk
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Demon
//
//	Play: Discard 2 cards from your hand.
func TestOldYurk(t *testing.T) {
	t.Run("discards 2 cards from hand", func(t *testing.T) {
		var a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand: ct.Cards(
					OldYurk,
					ct.Bind(&a, ct.Creature(ct.Power(3))),
					ct.Bind(&b, ct.Creature(ct.Power(4))),
				),
			},
		})

		h.P1.Play(OldYurk)
		h.P1.ClickCard(a)

		h.Expect(a).At(ct.Discard)
		h.Expect(b).At(ct.Discard)
	})
}
