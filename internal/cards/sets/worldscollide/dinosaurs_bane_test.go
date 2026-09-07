package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dinosaurs' Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Destroy a Dinosaur creature.
func TestDinosaursBane(t *testing.T) {
	t.Run("destroys a Dinosaur creature and gains a pip", func(t *testing.T) {
		var dino, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(DinosaursBane)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&dino, ct.Creature(ct.Traits(card.Traits.Dinosaur))),
				ct.Bind(&other, ct.Creature()),
			)},
		})

		h.P1.Play(DinosaursBane)

		h.Expect(dino).At(ct.Discard)
		h.Expect(other).At(ct.PlayArea)
		h.P1.ExpectAmber(1)
	})
}
