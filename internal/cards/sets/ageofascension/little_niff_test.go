package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Little Niff
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Elf • Thief
//
//	Deploy, Elusive, Omega.
//	After a neighbor of Little Niff is used to fight, steal 1 Æmber.
func TestLittleNiff(t *testing.T) {
	var niff, attacker, defender ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Shadows,
			InPlay: ct.Cards(
				ct.Bind(&niff, LittleNiff),
				ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(4))),
			),
		},
		P2: ct.Side{
			Amber:  2,
			InPlay: ct.Cards(ct.Bind(&defender, ct.Creature(ct.Power(3)))),
		},
	})

	h.P1.Fight(attacker, defender)

	// The attacker fought beside Little Niff, so Niff steals 1 Æmber.
	if got := h.P1.Amber(); got != 1 {
		t.Errorf("P1 Æmber = %d, want 1", got)
	}
	if got := h.P2.Amber(); got != 1 {
		t.Errorf("P2 Æmber = %d, want 1", got)
	}
}
