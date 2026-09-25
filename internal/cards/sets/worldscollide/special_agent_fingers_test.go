package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Special Agent "Fingers"
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	Action: Steal 1 Æmber.
func TestSpecialAgentFingers(t *testing.T) {
	t.Run("steals 1 Æmber with its Action", func(t *testing.T) {
		var fingers ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&fingers, SpecialAgentFingers)),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(fingers)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
	})
}
