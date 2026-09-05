package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lost in the Woods
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Shuffle 2 friendly creatures into their owners' decks, and shuffle 2 enemy creatures into their owners' decks.
func TestLostInTheWoods(t *testing.T) {
	t.Run("shuffles 2 friendly and 2 enemy creatures into their owners' decks", func(t *testing.T) {
		var ally1, ally2, foe1, foe2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(LostInTheWoods),
				InPlay: ct.Cards(
					ct.Bind(&ally1, ct.Creature(ct.OfHouse(card.House.Untamed))),
					ct.Bind(&ally2, ct.Creature(ct.OfHouse(card.House.Untamed))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe1, ct.Creature(ct.OfHouse(card.House.Mars))),
					ct.Bind(&foe2, ct.Creature(ct.OfHouse(card.House.Mars))),
				),
			},
		})

		h.P1.Play(LostInTheWoods)
		// Each pick narrows the pool; the last candidate on each side is automatic.
		h.P1.ClickCard(ally1)
		h.P1.ClickCard(foe1)

		h.Expect(ally1).At(ct.Deck)
		h.Expect(ally2).At(ct.Deck)
		h.Expect(foe1).At(ct.Deck)
		h.Expect(foe2).At(ct.Deck)
	})
}
