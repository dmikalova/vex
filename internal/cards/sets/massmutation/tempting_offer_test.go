package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tempting Offer
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Put an enemy creature into its owner's hand -> your opponent gains 1 Æmber.
//	Enhance Capture.
func TestTemptingOffer(t *testing.T) {
	t.Run("returns an enemy creature to hand, its owner gains 1 Æmber", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(TemptingOffer)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
		})

		h.P1.Play(TemptingOffer)

		h.Expect(foe).At(ct.Hand)
		h.P2.ExpectAmber(1)
	})
}
