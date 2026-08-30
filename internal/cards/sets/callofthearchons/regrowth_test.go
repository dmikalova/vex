package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Regrowth
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Put a creature from your discard pile into your hand.
func TestRegrowth(t *testing.T) {
	t.Run("puts a creature from your discard pile into your hand", func(t *testing.T) {
		var buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Untamed,
				Hand:    ct.Cards(Regrowth),
				Discard: ct.Cards(ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Untamed)))),
			},
		})

		h.P1.Play(Regrowth)

		h.Expect(buried).At(ct.Hand)
	})
}
