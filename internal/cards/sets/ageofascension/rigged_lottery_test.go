package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Rigged Lottery
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Discard the top 5 cards of each player's deck. For each Shadows card discarded this way, its owner gains 1 Æmber.
func TestRiggedLottery(t *testing.T) {
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Shadows,
			Hand:  ct.Cards(RiggedLottery),
			Deck:  ct.DeckOf(card.House.Shadows, 5),
		},
		P2: ct.Side{
			Deck: ct.DeckOf(card.House.Mars, 5),
		},
	})

	h.P1.Play(RiggedLottery)

	// P1's five Shadows cards each pay their owner 1 Æmber, on top of the card's
	// own 1 Æmber bonus. P2's Mars cards pay nothing.
	if got := h.P1.Amber(); got != 6 {
		t.Errorf("P1 Æmber = %d, want 6", got)
	}
	if got := h.P2.Amber(); got != 0 {
		t.Errorf("P2 Æmber = %d, want 0", got)
	}
	if got := len(h.Game().Discard(0)); got != 6 {
		t.Errorf("P1 discard = %d, want 6", got)
	}
	if got := len(h.Game().Discard(1)); got != 5 {
		t.Errorf("P2 discard = %d, want 5", got)
	}
}
