package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Helmsman Spears
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Human
//
//	Fight/Reap: Discard any number of cards from your hand -> for each card discarded this way, draw a card.
func TestHelmsmanSpears(t *testing.T) {
	var spears, toss1, toss2, drawn1, drawn2 ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.StarAlliance,
			InPlay: ct.Cards(ct.Bind(&spears, HelmsmanSpears)),
			Hand: ct.Cards(
				ct.Bind(&toss1, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
				ct.Bind(&toss2, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
			),
			Deck: ct.Cards(
				ct.Bind(&drawn1, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
				ct.Bind(&drawn2, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
			),
		},
	})

	h.P1.Reap(spears)
	h.P1.ClickCard(toss1)
	h.P1.ClickCard(toss2)

	// The two discarded cards leave the hand; a card is drawn for each.
	h.Expect(toss1).At(ct.Discard)
	h.Expect(toss2).At(ct.Discard)
	h.Expect(drawn1).At(ct.Hand)
	h.Expect(drawn2).At(ct.Hand)
}
