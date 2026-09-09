package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Philophosaurus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Dinosaur • Philosopher
//
//	Reap: You may look at the top 3 cards of your deck, archive 1, put 1 into your hand, and discard 1.
func TestPhilophosaurus(t *testing.T) {
	var top, middle, third, bottom ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Saurian,
			InPlay: ct.Cards(Philophosaurus),
			Deck: ct.Cards(
				ct.Bind(&top, ct.Creature(ct.Power(3))),
				ct.Bind(&middle, ct.Creature(ct.Power(3))),
				ct.Bind(&third, ct.Creature(ct.Power(3))),
				ct.Bind(&bottom, ct.Creature(ct.Power(3))),
			),
		},
	})

	h.P1.Reap(Philophosaurus)
	h.P1.ClickOption("Yes") // take the optional look
	h.P1.ClickCard(top)     // archive one
	h.P1.ClickCard(middle)  // put one into hand; the last is discarded

	h.Expect(top).At(ct.Archives)
	h.Expect(middle).At(ct.Hand)
	h.Expect(third).At(ct.Discard)
	h.Expect(bottom).At(ct.Deck)
}
