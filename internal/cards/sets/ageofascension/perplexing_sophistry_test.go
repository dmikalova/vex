package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Perplexing Sophistry
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If you have more Æmber than your opponent, your opponent discards a random card from their hand, and you draw a card.
func TestPerplexingSophistry(t *testing.T) {
	t.Run("opponent discards and you draw when you have more aember", func(t *testing.T) {
		var mine, theirs ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Amber: 5, // 6 after the play bonus
				Hand:  ct.Cards(PerplexingSophistry),
				Deck:  ct.Cards(ct.Bind(&mine, ct.Creature(ct.Power(1)))),
			},
			P2: ct.Side{
				Amber: 1,
				Hand:  ct.Cards(ct.Bind(&theirs, ct.Creature(ct.Power(1)))),
			},
		})

		h.P1.Play(PerplexingSophistry)

		h.Expect(theirs).At(ct.Discard)
		h.Expect(mine).At(ct.Hand)
	})

	t.Run("nothing happens when you do not have more aember", func(t *testing.T) {
		var theirs ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Amber: 0, // 1 after the play bonus
				Hand:  ct.Cards(PerplexingSophistry),
			},
			P2: ct.Side{
				Amber: 5,
				Hand:  ct.Cards(ct.Bind(&theirs, ct.Creature(ct.Power(1)))),
			},
		})

		h.P1.Play(PerplexingSophistry)

		h.Expect(theirs).At(ct.Hand)
	})
}
