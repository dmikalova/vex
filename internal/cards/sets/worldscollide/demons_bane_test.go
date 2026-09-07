package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Demons' Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy a Demon creature.
func TestDemonsBane(t *testing.T) {
	t.Run("destroys a Demon creature and gains a pip", func(t *testing.T) {
		var demon ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(DemonsBane)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&demon, ct.Creature(ct.Traits(card.Traits.Demon))),
			)},
		})

		h.P1.Play(DemonsBane)

		h.Expect(demon).At(ct.Discard)
		h.P1.ExpectAmber(1)
	})
}
