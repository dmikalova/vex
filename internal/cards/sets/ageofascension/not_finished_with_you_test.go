package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Not Finished with You
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Shuffle any number of Creatures from your discard pile into your deck.
func TestNotFinishedWithYou(t *testing.T) {
	var a, b ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Dis,
			Hand:  ct.Cards(NotFinishedWithYou),
			Discard: ct.Cards(
				ct.Bind(&a, ct.Creature()),
				ct.Bind(&b, ct.Creature()),
			),
		},
	})

	h.P1.Play(NotFinishedWithYou)
	h.P1.ClickCard(a)
	h.P1.ClickCard(b)

	h.Expect(a).At(ct.Deck)
	h.Expect(b).At(ct.Deck)
}
