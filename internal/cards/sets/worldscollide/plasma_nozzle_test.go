package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Plasma Nozzle
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains +2 assault and +2 splash-attack.
func TestPlasmaNozzle(t *testing.T) {
	t.Run(
		"host deals 2 assault to the fought creature and 2 splash to its neighbors",
		func(t *testing.T) {
			var host, left, target, right ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.StarAlliance,
					InPlay: ct.Cards(
						ct.Upgraded(
							ct.Bind(
								&host,
								ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
							),
							PlasmaNozzle,
						),
					),
				},
				P2: ct.Side{InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(10))),
					ct.Bind(&target, ct.Creature(ct.Power(10))),
					ct.Bind(&right, ct.Creature(ct.Power(10))),
				)},
			})

			h.P1.Fight(host, target)

			// Assault 2 deals 2 to the fought creature, then the 4-power host
			// deals its fight damage, for 6 total; splash-attack 2 hits each neighbor.
			h.Expect(target).At(ct.PlayArea).Damage(6)
			h.Expect(left).At(ct.PlayArea).Damage(2)
			h.Expect(right).At(ct.PlayArea).Damage(2)
		},
	)
}
