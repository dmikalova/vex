package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Regrettable Meteor
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Destroy each Dinosaur creature and each creature with power 6 or higher.
func TestRegrettableMeteor(t *testing.T) {
	t.Run(
		"destroys every Dinosaur and every power-6+ creature, sparing the rest",
		func(t *testing.T) {
			var dino, big, small, bigDino ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{House: card.House.Saurian, Hand: ct.Cards(RegrettableMeteor)},
				P2: ct.Side{InPlay: ct.Cards(
					ct.Bind(&dino, ct.Creature(ct.Traits(card.Traits.Dinosaur), ct.Power(2))),
					ct.Bind(&big, ct.Creature(ct.Power(6))),
					ct.Bind(&small, ct.Creature(ct.Power(3))),
					ct.Bind(&bigDino, ct.Creature(ct.Traits(card.Traits.Dinosaur), ct.Power(7))),
				)},
			})

			h.P1.Play(RegrettableMeteor)

			h.Expect(dino).At(ct.Discard)
			h.Expect(big).At(ct.Discard)
			h.Expect(bigDino).At(ct.Discard)
			h.Expect(small).At(ct.PlayArea)
		},
	)
}
