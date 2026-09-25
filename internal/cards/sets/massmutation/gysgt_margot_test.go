package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// GySgt. Margot
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	Fight/Reap: Deal 2 damage to an enemy creature. Ward a friendly creature.
func TestGySgtMargot(t *testing.T) {
	t.Run(
		"deals 2 damage to an enemy and wards a friendly creature after reaping",
		func(t *testing.T) {
			var margot, enemy ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.StarAlliance,
					InPlay: ct.Cards(ct.Bind(&margot, GySgtMargot)),
				},
				P2: ct.Side{InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				)},
			})

			h.P1.Reap(margot)

			h.Expect(enemy).Damage(2)
			if !h.Game().Warded(margot.ID()) {
				t.Error("Margot should ward the sole friendly creature (itself)")
			}
		},
	)
}
