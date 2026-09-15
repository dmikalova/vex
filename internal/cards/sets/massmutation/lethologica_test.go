package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lethologica
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Discard cards from the top of your deck until you discard a Logos card or run out of cards -> put the discarded card into your hand.
func TestLethologica(t *testing.T) {
	t.Run("digs to the first Logos card and takes it", func(t *testing.T) {
		var skipped, found, buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(Lethologica),
				Deck: ct.Cards(
					ct.Bind(&skipped, ct.Creature(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&found, ct.Creature(ct.OfHouse(card.House.Logos))),
					ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		h.P1.Play(Lethologica)

		h.Expect(found).At(ct.Hand)
		h.Expect(skipped).At(ct.Discard)
		h.Expect(buried).At(ct.Deck)
	})

	t.Run("empties the deck when it finds no Logos card", func(t *testing.T) {
		var offHouse ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(Lethologica),
				Deck: ct.Cards(
					ct.Bind(&offHouse, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				),
			},
		})

		h.P1.Play(Lethologica)

		h.Expect(offHouse).At(ct.Discard)
	})
}
