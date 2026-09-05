package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Reverse Time
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Swap your deck and your discard pile, then shuffle your deck.
func TestReverseTime(t *testing.T) {
	t.Run("swaps the deck and the discard pile", func(t *testing.T) {
		var inDeck, inDiscard ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Logos,
				Hand:    ct.Cards(ReverseTime),
				Deck:    ct.Cards(ct.Bind(&inDeck, ct.Creature(ct.OfHouse(card.House.Logos)))),
				Discard: ct.Cards(ct.Bind(&inDiscard, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
		})

		h.P1.Play(ReverseTime)

		h.Expect(inDiscard).At(ct.Deck)
		// Reverse Time itself is discarded as it resolves, so the old deck card
		// joins it in the new discard pile.
		h.Expect(inDeck).At(ct.Discard)
	})
}
