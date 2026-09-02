package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Invasion Portal
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Discard cards from the top of your deck until you discard a Mars creature or run out of cards -> put it into your hand.
func TestInvasionPortal(t *testing.T) {
	t.Run("digs to the first Mars creature and takes it", func(t *testing.T) {
		var portal, skipped, martian ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ct.Bind(&portal, InvasionPortal)),
				Deck: ct.Cards(
					ct.Bind(&skipped, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
					ct.Bind(&martian, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.UseAction(portal)

		h.Expect(martian).At(ct.Hand)
		h.Expect(skipped).At(ct.Discard)
	})
}
