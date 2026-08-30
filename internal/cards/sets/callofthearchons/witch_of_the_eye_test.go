package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Witch of the Eye
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Witch
//
//	Reap: Put a card from your discard pile into your hand.
func TestWitchOfTheEye(t *testing.T) {
	t.Run("puts a card from your discard pile into your hand when it reaps", func(t *testing.T) {
		var buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Untamed,
				InPlay:  ct.Cards(WitchOfTheEye),
				Discard: ct.Cards(ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Untamed)))),
			},
		})

		h.P1.Reap(WitchOfTheEye)

		h.Expect(buried).At(ct.Hand)
	})
}
