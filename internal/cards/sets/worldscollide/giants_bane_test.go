package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Giants' Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Destroy a Giant creature.
func TestGiantsBane(t *testing.T) {
	t.Run("destroys a Giant creature", func(t *testing.T) {
		var giant ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(GiantsBane)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&giant, ct.Creature(ct.Power(6), ct.Traits(card.Traits.Giant))),
				),
			},
		})

		h.P1.Play(GiantsBane)

		h.Expect(giant).At(ct.Discard)
	})
}
