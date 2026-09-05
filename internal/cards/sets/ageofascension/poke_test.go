package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Poke
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 1 damage to an enemy creature. If this damage destroys that creature, draw a card.
func TestPoke(t *testing.T) {
	t.Run("draws a card when the damage destroys the creature", func(t *testing.T) {
		var enemy, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(Poke),
				Deck:  ct.Cards(ct.Bind(&top, ct.Creature(ct.Power(3)))),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&enemy, ct.Creature(ct.Power(1))),
			)},
		})

		h.P1.Play(Poke)

		h.Expect(enemy).At(ct.Discard)
		h.Expect(top).At(ct.Hand)
	})

	t.Run("draws no card when the creature survives", func(t *testing.T) {
		var enemy, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(Poke),
				Deck:  ct.Cards(ct.Bind(&top, ct.Creature(ct.Power(3)))),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&enemy, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Play(Poke)

		h.Expect(enemy).Damage(1)
		h.Expect(top).At(ct.Deck)
	})
}
