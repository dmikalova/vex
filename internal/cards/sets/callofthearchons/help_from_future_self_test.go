package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Help from Future Self
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: FIXED
//	Æmber:  1
//
//	Play: Search your deck and discard pile for a Timetraveller, reveal it, and put it into your hand, and shuffle your discard pile into your deck.
func TestHelpFromFutureSelf(t *testing.T) {
	var tt, buried ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:   card.House.Logos,
			Hand:    ct.Cards(HelpFromFutureSelf),
			Deck:    ct.Cards(ct.Bind(&tt, Timetraveller)),
			Discard: ct.Cards(ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Logos)))),
		},
	})

	h.P1.Play(HelpFromFutureSelf)

	h.Expect(tt).At(ct.Hand)     // the Timetraveller is tutored to hand
	h.Expect(buried).At(ct.Deck) // the discard pile is shuffled into the deck
}
