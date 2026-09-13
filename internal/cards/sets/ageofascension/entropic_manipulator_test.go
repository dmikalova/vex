package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Entropic Manipulator
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Redistribute the damage among a player's Creatures.
func TestEntropicManipulator(t *testing.T) {
	t.Run("redistributes damage among the chosen player's creatures", func(t *testing.T) {
		var a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(EntropicManipulator),
				InPlay: ct.Cards(
					ct.Bind(&a, ct.Creature(ct.Power(5))),
					ct.Bind(&b, ct.Creature(ct.Power(5))),
				),
			},
		})
		a.Damaged(2)

		h.P1.Play(EntropicManipulator)
		h.P1.ClickOption("P1")
		h.P1.ClickOption("Yes")
		// Move the 2 damage from a onto b.
		h.P1.ClickCard(b)
		h.P1.ClickCard(b)

		h.Expect(a).Damage(0)
		h.Expect(b).Damage(2)
	})

	t.Run("may pile damage past a creature's power to destroy it", func(t *testing.T) {
		var a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(EntropicManipulator),
				InPlay: ct.Cards(
					ct.Bind(&a, ct.Creature(ct.Power(2))),
					ct.Bind(&b, ct.Creature(ct.Power(2))),
				),
			},
		})
		a.Damaged(1)
		b.Damaged(1)

		h.P1.Play(EntropicManipulator)
		h.P1.ClickOption("P1")
		h.P1.ClickOption("Yes")
		// Pile both damage onto a to lethal.
		h.P1.ClickCard(a)
		h.P1.ClickCard(a)

		if h.Game().InPlay(a.ID()) {
			t.Error("a should be destroyed")
		}
		if !h.Game().InPlay(b.ID()) {
			t.Error("b should survive")
		}
	})

	t.Run("may decline to redistribute", func(t *testing.T) {
		var a ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				Hand:   ct.Cards(EntropicManipulator),
				InPlay: ct.Cards(ct.Bind(&a, ct.Creature(ct.Power(5)))),
			},
		})
		a.Damaged(2)

		h.P1.Play(EntropicManipulator)
		h.P1.ClickOption("P1")
		h.P1.ClickOption("No")

		h.Expect(a).Damage(2)
	})
}
