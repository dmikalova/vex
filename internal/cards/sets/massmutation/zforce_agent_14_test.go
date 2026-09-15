package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Z-Force Agent 14
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Cyborg
//
//	Fight: For each upgrade on Z-Force Agent 14, gain 1 Æmber.
func TestZForceAgent14(t *testing.T) {
	t.Run("gains no Æmber fighting with no upgrades", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ZForceAgent14),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Fight(ZForceAgent14, enemy)

		h.P1.ExpectAmber(0)
	})

	t.Run("gains 1 Æmber for each upgrade on it", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(ZForceAgent14, ct.Upgrade(), ct.Upgrade()),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Fight(ZForceAgent14, enemy)

		h.P1.ExpectAmber(2)
	})
}
