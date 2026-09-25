package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Mastermindy
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Elf • Thief
//
//	Elusive.
//	At the end of your turn, put a scheme counter on Mastermindy.
//	Action: For each scheme counter on Mastermindy, steal 1 Æmber. Remove each scheme counter from Mastermindy.
func TestMastermindy(t *testing.T) {
	t.Run(
		"accrues a scheme counter each turn and steals one Æmber per counter",
		func(t *testing.T) {
			var mindy ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Shadows,
					InPlay: ct.Cards(ct.Bind(&mindy, Mastermindy)),
				},
				P2: ct.Side{
					House: card.House.Shadows,
					Amber: 5,
				},
			})

			// Two of P1's turns end, each placing a scheme counter.
			h.P1.EndTurn()
			h.P2.ChooseHouse(card.House.Shadows)
			h.P2.EndTurn()
			h.P1.ChooseHouse(card.House.Shadows)
			h.P1.EndTurn()
			h.P2.ChooseHouse(card.House.Shadows)
			h.P2.EndTurn()

			// A third turn: cash the two scheme counters in for two stolen Æmber.
			h.P1.ChooseHouse(card.House.Shadows)
			h.P1.UseAction(Mastermindy)

			h.P1.ExpectAmber(2)
			h.P2.ExpectAmber(3)
		},
	)

	t.Run("with no scheme counters the Action steals nothing", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(Mastermindy),
			},
			P2: ct.Side{
				House: card.House.Shadows,
				Amber: 4,
			},
		})

		h.P1.UseAction(Mastermindy)

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(4)
	})
}
