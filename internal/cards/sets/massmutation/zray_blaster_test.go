package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Z-Ray Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Connected
//	Bonus:  Æmber
//
//	This creature gains +3 power and +3 splash-attack.
func TestZRayBlaster(t *testing.T) {
	t.Run("adds 3 power to its host", func(t *testing.T) {
		var host ct.Card
		ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&host,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						ZRayBlaster,
					),
				),
			},
		})

		if got := host.Power(); got != 7 {
			t.Errorf("host power = %d, want 7", got)
		}
	})

	t.Run("damages the neighbors of the creature its host fights", func(t *testing.T) {
		var host, left, target, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&host,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						ZRayBlaster,
					),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(6))),
					ct.Bind(&target, ct.Creature(ct.Power(3))),
					ct.Bind(&right, ct.Creature(ct.Power(6))),
				),
			},
		})

		h.P1.Fight(host, target)

		h.Expect(left).Damage(3)
		h.Expect(right).Damage(3)
	})
}
