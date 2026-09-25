package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Panpaca, Anga
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Beast
//
//	Each creature to the right of Panpaca, Anga gains +2 power.
func TestPanpacaAnga(t *testing.T) {
	var left, right ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.Power(4))),
				PanpacaAnga,
				ct.Bind(&right, ct.Creature(ct.Power(4))),
			),
		},
	})
	// Only the creature to Anga's right gets +2 power.
	h.Expect(right).Power(6)
	h.Expect(left).Power(4)
}
