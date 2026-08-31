package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Flame-Wreathed
//
//	House:  Dis
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains +2 power and +2 hazardous.
func TestFlameWreathed(t *testing.T) {
	t.Run("grants the host +2 power and +2 hazardous", func(t *testing.T) {
		var enemy, attacker ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(1))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
						FlameWreathed,
					),
				),
			},
		})

		h.Expect(enemy).Power(6) // 4 + 2

		h.P1.Fight(attacker, enemy)

		h.Expect(attacker).At(ct.Discard) // destroyed by the granted Hazardous 2
	})
}
