package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Song of Spring
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Shuffle any number of friendly Untamed creatures from your hand, discard pile, or battleline into your deck.
func TestSongOfSpring(t *testing.T) {
	var inHand, inDiscard, onBoard ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			Hand: ct.Cards(
				SongOfSpring,
				ct.Bind(&inHand, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
			),
			Discard: ct.Cards(
				ct.Bind(&inDiscard, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
			),
			InPlay: ct.Cards(
				ct.Bind(&onBoard, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
			),
		},
	})

	h.P1.Play(SongOfSpring)
	h.P1.ClickCard(inHand)
	h.P1.ClickCard(inDiscard)
	h.P1.ClickCard(onBoard)

	h.Expect(inHand).At(ct.Deck)
	h.Expect(inDiscard).At(ct.Deck)
	h.Expect(onBoard).At(ct.Deck)
}
