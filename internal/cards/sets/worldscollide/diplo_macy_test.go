package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Diplo-Macy
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Alpha.
//	Play: Until the start of your next turn, each creature gains, "Before Fight: Exalt this creature."
func TestDiploMacy(t *testing.T) {
	t.Run("a friendly creature exalts itself when it fights", func(t *testing.T) {
		var attacker, prey ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(DiploMacy),
				InPlay: ct.Cards(
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(6))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&prey, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(1))),
				),
			},
		})

		h.P1.Play(DiploMacy)
		h.P1.Fight(attacker, prey)

		// The attacker exalted itself before dealing fight damage.
		h.Expect(attacker).AmberOn(1)
	})

	t.Run(
		"an enemy creature exalts itself when it fights on the opponent's turn",
		func(t *testing.T) {
			var enemy, prey ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Saurian,
					Hand:  ct.Cards(DiploMacy),
					InPlay: ct.Cards(
						ct.Bind(&prey, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(1))),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6))),
					),
				},
			})

			h.P1.Play(DiploMacy)
			h.P1.EndTurn()
			h.P2.ChooseHouse(card.House.Brobnar)
			h.P2.Fight(enemy, prey)

			// The grant lasts until the caster's next turn, so an enemy creature
			// fighting on the opponent's turn exalts itself too.
			h.Expect(enemy).AmberOn(1)
		},
	)
}
