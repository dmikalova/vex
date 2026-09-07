package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Thieves' Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Destroy a Thief creature.
func TestThievesBane(t *testing.T) {
	t.Run("destroys a Thief creature and gains a pip", func(t *testing.T) {
		var thief ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(ThievesBane)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&thief, ct.Creature(ct.Traits(card.Traits.Thief))),
			)},
		})

		h.P1.Play(ThievesBane)

		h.Expect(thief).At(ct.Discard)
		h.P1.ExpectAmber(1)
	})
}
