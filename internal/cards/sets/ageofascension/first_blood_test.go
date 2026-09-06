package ageofascension_test

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/cards/sets/ageofascension"
)

// First Blood
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Alpha.
//	Play: Deal 2 damage for each friendly Brobnar creature, divided among any number of creatures.
func TestFirstBlood(t *testing.T) {
	var blood, mook ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Brobnar,
			Hand:  ct.Cards(ct.Bind(&blood, ageofascension.FirstBlood)),
			InPlay: ct.Cards(
				ct.Bind(&mook, ct.Creature(ct.Power(9), ct.OfHouse(card.House.Brobnar))),
			),
		},
	})

	h.P1.Play(blood)

	// One friendly Brobnar creature means a pool of 2 damage, which the sole
	// candidate (that creature) takes in full.
	if got := h.Game().Damage(mook.ID()); got != 2 {
		t.Errorf("damage on the Brobnar creature = %d, want 2", got)
	}
}
