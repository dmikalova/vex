package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Cowfyne's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains +2 splash-attack.
func TestCowfynesBrew(t *testing.T) {
	t.Run("host deals 2 damage to each neighbor of the creature it fights", func(t *testing.T) {
		var host, left, target, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
						CowfynesBrew,
					),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.Power(6))),
				ct.Bind(&target, ct.Creature(ct.Power(6))),
				ct.Bind(&right, ct.Creature(ct.Power(6))),
			)},
		})

		h.P1.Fight(host, target)

		h.Expect(left).At(ct.PlayArea).Damage(2)
		h.Expect(right).At(ct.PlayArea).Damage(2)
	})
}
