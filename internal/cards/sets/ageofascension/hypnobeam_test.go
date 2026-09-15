package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hypnobeam
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Take control of an enemy creature.
func TestHypnobeam(t *testing.T) {
	t.Run("takes control of an enemy creature when played", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, Hand: ct.Cards(Hypnobeam)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
		})

		h.P1.Play(Hypnobeam)

		inMine := false
		for _, id := range h.Game().Battleline(0) {
			if id == foe.ID() {
				inMine = true
			}
		}
		if !inMine {
			t.Error("the seized creature should be in the controller's battleline")
		}
	})
}
