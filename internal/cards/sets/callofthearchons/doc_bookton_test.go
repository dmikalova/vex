package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Doc Bookton
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Human • Scientist
//
//	Reap: Draw a card.
func TestDocBookton(t *testing.T) {
	t.Run("draws a card when it reaps", func(t *testing.T) {
		var doc, drawn ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&doc, DocBookton)),
				Deck: ct.Cards(
					ct.Bind(&drawn, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(1))),
				),
			},
		})

		h.Expect(doc).Power(5)
		h.P1.Reap(doc)

		h.Expect(drawn).At(ct.Hand) // drew the top card
	})
}
