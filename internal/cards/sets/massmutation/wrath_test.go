package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Wrath
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Armor:  3
//	Traits: Demon • Sin
//
//	Taunt, Poison, Skirmish.
//	Fight: For each friendly Sin creature, enrage an enemy creature.
func TestWrath(t *testing.T) {
	t.Run("enrages an enemy creature for each friendly Sin creature", func(t *testing.T) {
		var wrath, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ct.Bind(&wrath, Wrath),
					ct.Creature(
						ct.OfHouse(card.House.Dis),
						ct.Traits(card.Traits.Sin),
						ct.Power(3),
					),
				),
			},
			P2: ct.Side{
				// Armor 3 absorbs Wrath's 3 power, so poison does not destroy the foe
				// and it survives to be enraged.
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(6), ct.Armor(3)))),
			},
		})

		h.P1.Fight(wrath, foe)

		if !h.Game().Enraged(foe.ID()) {
			t.Error("enemy creature should be enraged")
		}
	})
}
