package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Beasts' Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy a Beast creature.
func TestBeastsBane(t *testing.T) {
	t.Run("destroys a Beast creature and gains a pip", func(t *testing.T) {
		var beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(BeastsBane)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&beast, ct.Creature(ct.Traits(card.Traits.Beast))),
			)},
		})

		h.P1.Play(BeastsBane)

		h.Expect(beast).At(ct.Discard)
		h.P1.ExpectAmber(1)
	})
}
