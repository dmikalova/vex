package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Font of the Eye
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Common
//	Traits: Location
//
//	Versatile.
//	Action: If an enemy creature has been destroyed this turn, a friendly creature captures 1 Æmber from your opponent.
func TestFontOfTheEye(t *testing.T) {
	t.Run(
		"a friendly creature captures 1 after an enemy creature was destroyed",
		func(t *testing.T) {
			var ally, enemy ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Sanctum,
					InPlay: ct.Cards(
						FontOfTheEye,
						ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(6))),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(1)))),
					Amber:  3,
				},
			})

			h.P1.Fight(ally, enemy) // destroys the enemy creature this turn
			h.P1.UseAction(FontOfTheEye)

			h.Expect(ally).AmberOn(1)
			h.P2.ExpectAmber(2)
		},
	)

	t.Run("captures nothing when no enemy creature was destroyed", func(t *testing.T) {
		var ally, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					FontOfTheEye,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(6))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(1)))),
				Amber:  3,
			},
		})

		h.P1.UseAction(FontOfTheEye)

		h.Expect(ally).AmberOn(0)
		h.P2.ExpectAmber(3)
	})
}
