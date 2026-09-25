package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Humble
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Exhaust a creature -> move 3 Æmber from the chosen creature to the common supply.
func TestHumble(t *testing.T) {
	t.Run("exhausts a creature and moves 3 Æmber off it to the supply", func(t *testing.T) {
		var rich ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(Humble),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&rich, ct.Creature(ct.Power(4))))},
		})

		h.Game().AddAmberOn(rich.ID(), 3)
		h.P1.Play(Humble) // rich is the sole creature, so it is the forced target.

		h.Expect(rich).AmberOn(0)
		if !rich.Exhausted() {
			t.Error("the chosen creature should be exhausted")
		}
	})

	t.Run("empties a creature holding fewer than 3 Æmber", func(t *testing.T) {
		var poor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(Humble),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&poor, ct.Creature(ct.Power(4))))},
		})

		h.Game().AddAmberOn(poor.ID(), 1)
		h.P1.Play(Humble)

		h.Expect(poor).AmberOn(0)
		if !poor.Exhausted() {
			t.Error("the chosen creature should be exhausted")
		}
	})
}
