package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Random Access Archives
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Archive the top card of your deck.
func TestRandomAccessArchives(t *testing.T) {
	t.Run("archives the top card of the deck", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(RandomAccessArchives),
				Deck: ct.Cards(
					ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2))),
				),
			},
		})

		h.P1.Play(RandomAccessArchives)

		h.Expect(top).At(ct.Archives)
	})
}
