package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Silent Dagger
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "Reap: Deal 4 damage to a flank creature."
func TestSilentDagger(t *testing.T) {
	t.Run("grants the host Reap: deal 4 to a flank creature", func(t *testing.T) {
		var host, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3))),
						SilentDagger,
					),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.P1.Reap(host)
		h.P1.ClickCard(foe) // the host is also on a flank, so pick the enemy

		h.Expect(foe).Damage(4)
	})
}
