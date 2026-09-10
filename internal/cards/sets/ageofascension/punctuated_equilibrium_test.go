package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Punctuated Equilibrium
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Each player discards their hand. Each player refills their hand as if it were the end of their turn.
func TestPunctuatedEquilibrium(t *testing.T) {
	var mine, theirs ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			Hand: ct.Cards(
				PunctuatedEquilibrium,
				ct.Bind(&mine, ct.Creature(ct.Power(3))),
			),
			Deck: ct.DeckOf(card.House.Untamed, 8),
		},
		P2: ct.Side{
			Hand: ct.Cards(ct.Bind(&theirs, ct.Creature(ct.Power(3)))),
			Deck: ct.DeckOf(card.House.Brobnar, 8),
		},
	})

	h.P1.Play(PunctuatedEquilibrium)

	h.Expect(mine).At(ct.Discard)
	h.Expect(theirs).At(ct.Discard)
	if got := len(h.Game().Hand(0)); got != 6 {
		t.Errorf("P1 hand = %d, want 6", got)
	}
	if got := len(h.Game().Hand(1)); got != 6 {
		t.Errorf("P2 hand = %d, want 6", got)
	}
}
