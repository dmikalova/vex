package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Labwork
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Archive a card from your hand.
func TestLabwork(t *testing.T) {
	t.Run("archives a card from hand", func(t *testing.T) {
		var spare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					Labwork,
					ct.Bind(&spare, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2))),
				),
			},
		})

		h.P1.Play(Labwork)

		h.Expect(spare).At(ct.Archives)
	})
}
