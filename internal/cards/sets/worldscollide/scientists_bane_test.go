package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Scientists' Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy a Scientist creature.
func TestScientistsBane(t *testing.T) {
	t.Run("destroys a Scientist creature and gains a pip", func(t *testing.T) {
		var scientist ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(ScientistsBane)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&scientist, ct.Creature(ct.Traits(card.Traits.Scientist))),
			)},
		})

		h.P1.Play(ScientistsBane)

		h.Expect(scientist).At(ct.Discard)
		h.P1.ExpectAmber(1)
	})
}
