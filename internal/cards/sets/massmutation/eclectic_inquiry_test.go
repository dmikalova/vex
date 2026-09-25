package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Eclectic Inquiry
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Archive the top 2 cards of your deck.
func TestEclecticInquiry(t *testing.T) {
	t.Run("archives the top 2 cards of your deck", func(t *testing.T) {
		var first, second ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(EclecticInquiry),
				Deck: ct.Cards(
					ct.Bind(&first, ct.Creature()),
					ct.Bind(&second, ct.Creature()),
				),
			},
			P2: ct.Side{},
		})

		h.P1.Play(EclecticInquiry)

		h.Expect(first).At(ct.Archives)
		h.Expect(second).At(ct.Archives)
	})
}
