package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Squire Alys
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Armor:  2
//	Traits: Human
//
//	Play: Squire Alys captures 2 Æmber from your opponent.
func TestSquireAlys(t *testing.T) {
	t.Run("captures 2 Æmber from the opponent when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(SquireAlys),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(SquireAlys)

		h.Expect(SquireAlys).AmberOn(2)
		h.P2.ExpectAmber(1)
	})
}
