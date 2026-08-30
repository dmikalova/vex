package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Wardrummer
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Goblin
//
//	Play: Put each other friendly Brobnar creature into its owner's hand.
func TestWardrummer(t *testing.T) {
	t.Run("returns each other friendly Brobnar creature to hand", func(t *testing.T) {
		var brobAlly, marsAlly ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(Wardrummer),
				InPlay: ct.Cards(
					ct.Bind(&brobAlly, ct.Creature(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&marsAlly, ct.Creature(ct.OfHouse(card.House.Mars))),
				),
			},
		})

		h.P1.Play(Wardrummer)

		h.Expect(brobAlly).At(ct.Hand)
		h.Expect(marsAlly).At(ct.PlayArea)
		h.Expect(Wardrummer).At(ct.PlayArea)
	})
}
