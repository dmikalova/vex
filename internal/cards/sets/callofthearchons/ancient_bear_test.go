package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ancient Bear
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Beast
//
//	Assault 2.
func TestAncientBear(t *testing.T) {
	t.Run("deals 2 assault damage before fight damage and takes retaliation", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(AncientBear)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(10))),
				),
			},
		})

		h.P1.Fight(AncientBear, foe)

		h.Expect(foe).Damage(7)              // 2 assault + 5 fight
		h.Expect(AncientBear).At(ct.Discard) // destroyed by the 10-power foe's retaliation
	})
}
