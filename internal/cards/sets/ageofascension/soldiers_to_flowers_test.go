package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Soldiers to Flowers
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Purge each Untamed creature from each player's discard pile. For each card purged this way, its owner gains 1 Æmber.
func TestSoldiersToFlowers(t *testing.T) {
	var mine, spared, theirs ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			Hand:  ct.Cards(SoldiersToFlowers),
			Discard: ct.Cards(
				ct.Bind(&mine, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				ct.Bind(&spared, ct.Tactic()),
			),
		},
		P2: ct.Side{
			Discard: ct.Cards(ct.Bind(&theirs,
				ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4)))),
		},
	})

	h.P1.Play(SoldiersToFlowers)

	h.Expect(mine).At(ct.Purge)
	h.Expect(theirs).At(ct.Purge)
	h.Expect(spared).At(ct.Discard)
	// P1 gains 1 from the card's Æmber bonus on play, plus 1 for their purged
	// creature; P2 gains 1 for theirs.
	if got := h.P1.Amber(); got != 2 {
		t.Errorf("P1 Æmber = %d, want 2", got)
	}
	if got := h.P2.Amber(); got != 1 {
		t.Errorf("P2 Æmber = %d, want 1", got)
	}
}
