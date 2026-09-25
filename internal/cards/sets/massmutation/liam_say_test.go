package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Liam Say
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Elf
//
//	Elusive.
//	At the start of your turn, you may deal 1 damage to a creature.
func TestLiamSay(t *testing.T) {
	t.Run("may deal 1 damage to a creature at the start of your turn", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(LiamSay),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(5)))),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ClickCard(enemy)

		h.Expect(enemy).Damage(1)
	})

	t.Run("deals no damage when declined", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(LiamSay),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(5)))),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ClickDone()

		h.Expect(enemy).Damage(0)
	})
}
