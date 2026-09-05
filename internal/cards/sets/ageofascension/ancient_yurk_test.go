package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ancient Yurk
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Demon
//
//	Play: Discard 3 cards from your hand.
func TestAncientYurk(t *testing.T) {
	t.Run("discards 3 cards from hand", func(t *testing.T) {
		var a, b, c ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand: ct.Cards(
					AncientYurk,
					ct.Bind(&a, ct.Creature(ct.Power(3))),
					ct.Bind(&b, ct.Creature(ct.Power(4))),
					ct.Bind(&c, ct.Creature(ct.Power(5))),
				),
			},
		})

		h.P1.Play(AncientYurk)
		h.P1.ClickCard(a)
		h.P1.ClickCard(b)

		h.Expect(a).At(ct.Discard)
		h.Expect(b).At(ct.Discard)
		h.Expect(c).At(ct.Discard)
	})
}
