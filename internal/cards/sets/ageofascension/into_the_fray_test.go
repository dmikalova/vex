package ageofascension_test

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/cards/sets/ageofascension"
)

// Into the Fray
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: A friendly Brobnar creature gains, "Fight: Ready this creature."
func TestIntoTheFray(t *testing.T) {
	var fray, brute, foe ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Brobnar,
			Hand:  ct.Cards(ct.Bind(&fray, ageofascension.IntoTheFray)),
			InPlay: ct.Cards(
				ct.Bind(&brute, ct.Creature(ct.Power(5), ct.OfHouse(card.House.Brobnar))),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3)))),
		},
	})

	// Grant the Brobnar creature "Fight: Ready this creature", then have it fight;
	// it readies instead of staying exhausted.
	h.P1.Play(fray)
	h.P1.Fight(brute, foe)

	if h.Game().Exhausted(brute.ID()) {
		t.Error("the Brobnar creature should be readied after it fights")
	}
}
