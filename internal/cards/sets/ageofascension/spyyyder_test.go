package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Spyyyder
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Demon
//
//	Skirmish.
//	Spyyyder gains poison while attacking an enemy flank Creature.
func TestSpyyyder(t *testing.T) {
	t.Run("poisons a flank creature it fights", func(t *testing.T) {
		var flank ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Spyyyder)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&flank, ct.Creature(ct.Power(10))),
			)},
		})

		h.P1.Fight(Spyyyder, flank)

		h.Expect(flank).At(ct.Discard) // 2 damage + poison destroys it
	})

	t.Run("does not poison a creature off a flank", func(t *testing.T) {
		var middle ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Spyyyder)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Creature(ct.Power(10)),
				ct.Bind(&middle, ct.Creature(ct.Power(10))),
				ct.Creature(ct.Power(10)),
			)},
		})

		h.P1.Fight(Spyyyder, middle)

		h.Expect(middle).Damage(2) // survives: no poison off a flank
	})
}
