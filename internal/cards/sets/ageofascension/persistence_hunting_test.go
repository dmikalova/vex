package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Persistence Hunting
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose a house - exhaust each enemy creature of the chosen house.
func TestPersistenceHunting(t *testing.T) {
	t.Run("exhausts each enemy creature of the chosen house", func(t *testing.T) {
		var foe1, foe2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(PersistenceHunting)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe1, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
					ct.Bind(&foe2, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
				),
			},
		})

		h.P1.Play(PersistenceHunting)
		h.P1.ClickOption("Brobnar")

		if !foe1.Exhausted() {
			t.Error("foe1 should be exhausted")
		}
		if !foe2.Exhausted() {
			t.Error("foe2 should be exhausted")
		}
	})
}
