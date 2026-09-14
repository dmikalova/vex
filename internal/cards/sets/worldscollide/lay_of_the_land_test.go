package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lay of the Land
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Look at the top 3 cards of your deck and put them back in any order, and draw a card.
func TestLayOfTheLand(t *testing.T) {
	var top1, top2, top3 ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.StarAlliance,
			Hand:  ct.Cards(LayOfTheLand),
			Deck: ct.Cards(
				ct.Bind(&top1, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
				ct.Bind(&top2, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
				ct.Bind(&top3, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
			),
		},
	})

	h.P1.Play(LayOfTheLand)
	// Picks send top3 then top1 toward the bottom, leaving the unpicked top2 on
	// top; the draw that follows takes top2.
	h.P1.ClickCard(top3)
	h.P1.ClickCard(top1)

	h.Expect(top2).At(ct.Hand)
	h.P1.ExpectAmber(1)
}
