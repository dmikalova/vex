package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Sir Marrows
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	After an enemy creature reaps, Sir Marrows captures 1 Æmber from your opponent.
func TestSirMarrows(t *testing.T) {
	var marrows, reaper ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Brobnar,
			InPlay: ct.Cards(ct.Bind(&reaper, ct.Creature(ct.Power(3)))),
		},
		P2: ct.Side{
			InPlay: ct.Cards(ct.Bind(&marrows, SirMarrows)),
		},
	})

	// P1 (active) reaps: the Æmber it gains is captured by the opponent's Sir
	// Marrows instead of staying in P1's pool.
	h.P1.Reap(reaper)

	h.Expect(marrows).AmberOn(1)
	if got := h.P1.Amber(); got != 0 {
		t.Errorf("P1 pool = %d, want 0 (the reaped Æmber was captured)", got)
	}
}
