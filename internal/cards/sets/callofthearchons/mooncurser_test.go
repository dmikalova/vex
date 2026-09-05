package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mooncurser
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Elf • Thief
//
//	Skirmish, Poison.
//	Fight: Steal 1 Æmber.
func TestMooncurser(t *testing.T) {
	t.Run("steals 1 Æmber when it fights", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(Mooncurser)},
			P2: ct.Side{
				Amber: 3,
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
		})

		h.P1.Fight(Mooncurser, foe)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
		h.Expect(Mooncurser).At(ct.PlayArea).Damage(0) // skirmish
	})
}
