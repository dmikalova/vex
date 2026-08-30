package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mindwarper
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Martian • Scientist
//
//	Elusive.
//	Action: An enemy creature captures 1 Æmber from their own side.
func TestMindwarper(t *testing.T) {
	t.Run("an enemy creature captures 1 Æmber from its own side", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, InPlay: ct.Cards(Mindwarper)},
			P2: ct.Side{
				Amber:  3,
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar)))),
			},
		})

		h.P1.UseAction(Mindwarper)

		h.Expect(foe).AmberOn(1)
		h.P2.ExpectAmber(2)
	})
}
