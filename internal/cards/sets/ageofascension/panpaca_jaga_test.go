package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Panpaca, Jaga
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Beast
//
//	Skirmish.
//	Each creature to the left of Panpaca, Jaga gains skirmish.
func TestPanpacaJaga(t *testing.T) {
	var left, right ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.Power(4))),
				PanpacaJaga,
				ct.Bind(&right, ct.Creature(ct.Power(4))),
			),
		},
	})
	// Only the creature to Jaga's left gains skirmish.
	if !h.Game().HasKeyword(left.ID(), engine.Skirmish) {
		t.Error("creature to Jaga's left should gain skirmish")
	}
	if h.Game().HasKeyword(right.ID(), engine.Skirmish) {
		t.Error("creature to Jaga's right should not gain skirmish")
	}
}
