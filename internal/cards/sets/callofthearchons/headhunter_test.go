package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Headhunter
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Fight: Gain 1 Æmber.
func TestHeadhunter(t *testing.T) {
	t.Run("gains 1 Æmber when it fights", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(Headhunter)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
			},
		})

		h.P1.Fight(Headhunter, foe)

		h.P1.ExpectAmber(1)
		h.Expect(foe).At(ct.Discard)
	})
}
