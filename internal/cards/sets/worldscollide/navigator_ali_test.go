package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Navigator Ali
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Human
//
//	Play/Fight/Reap: Look at the top 3 cards of your deck and put them back in any order.
func TestNavigatorAli(t *testing.T) {
	var ali, top1, top2, top3 ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.StarAlliance,
			InPlay: ct.Cards(ct.Bind(&ali, NavigatorAli)),
			Deck: ct.Cards(
				ct.Bind(&top1, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
				ct.Bind(&top2, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
				ct.Bind(&top3, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
			),
		},
	})

	h.P1.Reap(ali)
	// Picks send cards toward the bottom in order: top3 first (deepest), top1
	// next, leaving the unpicked top2 riding on top (drawn next).
	h.P1.ClickCard(top3)
	h.P1.ClickCard(top1)

	deck := h.Game().Deck(0)
	if len(deck) < 3 ||
		deck[0] != top2.ID() || deck[1] != top1.ID() || deck[2] != top3.ID() {
		t.Errorf("deck top = %v, want [%d %d %d]", deck, top2.ID(), top1.ID(), top3.ID())
	}
}
