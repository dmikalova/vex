package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Smiling Ruth
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	Reap: If you forged a key this turn, take control of an enemy flank creature.
func TestSmilingRuth(t *testing.T) {
	t.Run("takes an enemy flank creature after forging", func(t *testing.T) {
		var ruth, flank, middle, otherFlank ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				Amber:  6,
				InPlay: ct.Cards(ct.Bind(&ruth, SmilingRuth)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&flank, ct.Creature(ct.Power(3))),
					ct.Bind(&middle, ct.Creature(ct.Power(4))),
					ct.Bind(&otherFlank, ct.Creature(ct.Power(5))),
				),
			},
		})

		h.P1.EndTurn()
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Shadows)
		h.P1.ExpectKeys(1)
		h.P1.Reap(ruth)
		h.P1.ClickCard(otherFlank)

		h.Expect(otherFlank).At(ct.PlayArea)
		inMine := false
		for _, id := range h.Game().Battleline(0) {
			if id == otherFlank.ID() {
				inMine = true
			}
		}
		if !inMine {
			t.Error("the seized flank creature should be in the controller's battleline")
		}
		if h.Game().Owner(middle.ID()) != 1 {
			t.Error("the middle creature should be untouched")
		}
	})

	t.Run("stays quiet without a key forged this turn", func(t *testing.T) {
		var ruth, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&ruth, SmilingRuth)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Reap(ruth)
		h.P1.ExpectAmber(1)
	})
}
