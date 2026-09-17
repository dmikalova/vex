package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Blast Shielding
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Common
//	Bonus:  Æmber
//
//	This creature gains +2 armor.
//	This creature gains, "After this creature is used, you may attach Blast Shielding to a friendly neighboring creature."
func TestBlastShielding(t *testing.T) {
	t.Run("moves onto a neighbor after the host is used", func(t *testing.T) {
		var host, neighbor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&host,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						BlastShielding,
					),
					ct.Bind(&neighbor, ct.Creature(ct.Power(3))),
				),
			},
		})

		h.Expect(host).Armor(2)
		h.Expect(neighbor).Armor(0)

		h.P1.Reap(host)
		h.P1.ClickCard(neighbor)

		h.Expect(BlastShielding).At(ct.Attached)
		h.Expect(host).Armor(0)
		h.Expect(neighbor).Armor(2)
	})

	t.Run("stays on the host when its controller declines", func(t *testing.T) {
		var host, neighbor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&host,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						BlastShielding,
					),
					ct.Bind(&neighbor, ct.Creature(ct.Power(3))),
				),
			},
		})

		h.P1.Reap(host)
		h.P1.ClickDone()

		h.Expect(host).Armor(2)
		h.Expect(neighbor).Armor(0)
	})
}
