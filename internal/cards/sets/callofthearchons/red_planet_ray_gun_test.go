package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Red Planet Ray Gun
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This Creature gains, "Reap: For each Mars Creature in play, deal 1 damage to a Creature."
func TestRedPlanetRayGun(t *testing.T) {
	t.Run("deals 1 damage per Mars creature in play, counting both players", func(t *testing.T) {
		var host, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
						RedPlanetRayGun,
					),
					ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5)),
				ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(6))),
			)},
		})

		h.P1.Reap(host)
		// Three Mars creatures in play, so three instances of 1 damage; all on enemy.
		h.P1.ClickCard(enemy)
		h.P1.ClickCard(enemy)
		h.P1.ClickCard(enemy)

		h.Expect(enemy).Damage(3) // three Mars creatures in play
	})
}
