package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ortannu's Binding
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Connected
//	Æmber:  1
//
//	Play: Deal 2 damage to a friendly creature.
func TestOrtannusBinding(t *testing.T) {
	t.Run("deals 2 damage to a friendly creature when played", func(t *testing.T) {
		var friend ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(OrtannusBinding),
				InPlay: ct.Cards(
					ct.Bind(&friend, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(20))),
				),
			},
			P2: ct.Side{},
		})

		h.P1.Play(OrtannusBinding)

		h.Expect(friend).Damage(2)
	})
}
