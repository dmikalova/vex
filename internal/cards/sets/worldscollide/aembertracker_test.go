package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Aembertracker
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Beast
//
//	Play: Deal 2 damage to each enemy creature with Æmber on it, ignoring armor.
func TestAembertracker(t *testing.T) {
	t.Run("deals 2 unpreventable damage to each enemy creature with Æmber", func(t *testing.T) {
		var withAember, armored, bare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(Aembertracker),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&withAember, ct.Creature(ct.Power(5))),
					ct.Bind(&armored, ct.Creature(ct.Power(5), ct.Armor(3))),
					ct.Bind(&bare, ct.Creature(ct.Power(5))),
				),
			},
		})
		h.Game().State.Cards[withAember.ID()].Amber = 1
		h.Game().State.Cards[armored.ID()].Amber = 1

		h.P1.Play(Aembertracker)

		h.Expect(withAember).Damage(2)
		h.Expect(armored).Damage(2) // armor does not prevent it
		h.Expect(bare).Damage(0)    // no Æmber, untouched
	})
}
