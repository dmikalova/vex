package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Humans' Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy a Human creature.
func TestHumansBane(t *testing.T) {
	t.Run("destroys a Human creature and gains a pip", func(t *testing.T) {
		var human ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(HumansBane)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&human, ct.Creature(ct.Traits(card.Traits.Human))),
			)},
		})

		h.P1.Play(HumansBane)

		h.Expect(human).At(ct.Discard)
		h.P1.ExpectAmber(1)
	})
}
