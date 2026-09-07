package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lilithal
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	Fight/Reap: Lilithal captures 1 Æmber from your opponent.
func TestLilithal(t *testing.T) {
	t.Run("captures 1 Æmber when it reaps", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Lilithal)},
			P2: ct.Side{Amber: 2},
		})

		h.P1.Reap(Lilithal)

		h.Expect(Lilithal).AmberOn(1)
		h.P2.ExpectAmber(1)
	})

	t.Run("captures 1 Æmber when it fights", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Lilithal)},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3)))),
				Amber:  2,
			},
		})

		h.P1.Fight(Lilithal, foe)

		h.Expect(Lilithal).AmberOn(1)
		h.P2.ExpectAmber(1)
	})
}
