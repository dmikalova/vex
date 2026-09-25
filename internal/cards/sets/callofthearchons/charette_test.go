package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Charette
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Demon
//
//	Play: Charette captures 3 Æmber from your opponent.
func TestCharette(t *testing.T) {
	t.Run("captures 3 Æmber when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Charette),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(Charette)

		h.Expect(Charette).Power(4).AmberOn(3)
		h.P2.ExpectAmber(2)
	})
}
