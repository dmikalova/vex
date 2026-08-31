package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Francus
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Armor:  1
//	Traits: Knight • Spirit
//
//	After a creature is destroyed fighting Francus, Francus captures 1 Æmber from your opponent.
func TestFrancus(t *testing.T) {
	t.Run("captures 1 Æmber when a creature is destroyed fighting it", func(t *testing.T) {
		var prey ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(Francus),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&prey, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
				Amber: 3,
			},
		})

		h.P1.Fight(Francus, prey)

		h.Expect(prey).At(ct.Discard) // the 1-power prey is destroyed
		h.Expect(Francus).AmberOn(1)  // Francus survives (armor 1) and captures 1
		h.P2.ExpectAmber(2)           // 3 - 1 captured
	})
}
