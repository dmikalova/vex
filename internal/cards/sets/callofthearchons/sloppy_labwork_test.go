package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sloppy Labwork
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Archive a card from your hand. Discard a card from your hand.
func TestSloppyLabwork(t *testing.T) {
	t.Run("archives one chosen card and discards another", func(t *testing.T) {
		var keep, toss ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					SloppyLabwork,
					ct.Bind(&keep, ct.Creature(ct.OfHouse(card.House.Logos))),
					ct.Bind(&toss, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		h.P1.Play(SloppyLabwork)
		h.P1.ClickCard(keep) // archive; toss is the sole remaining card to discard

		h.Expect(keep).At(ct.Archives)
		h.Expect(toss).At(ct.Discard)
	})
}
