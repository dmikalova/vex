package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Exhume
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Play a creature from your discard pile.
func TestExhume(t *testing.T) {
	t.Run("plays a creature from your own discard pile", func(t *testing.T) {
		var buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Exhume),
				Discard: ct.Cards(
					ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(4))),
				),
			},
		})

		h.P1.Play(Exhume)

		h.Expect(buried).At(ct.PlayArea)
	})
}
