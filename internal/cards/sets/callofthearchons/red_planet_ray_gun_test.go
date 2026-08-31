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
//	This creature gains, "Reap: For each Mars creature in play, deal 1 damage to a creature."
func TestRedPlanetRayGun(t *testing.T) {
	t.Run("deals 1 damage per Mars creature in play, counting both players", func(t *testing.T) {
		var host, victim ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Upgraded(ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))), RedPlanetRayGun),
					ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5)),
				ct.Bind(&victim, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(6))),
			)},
		})

		h.P1.Reap(host)
		h.P1.ClickCard(victim)

		h.Expect(victim).Damage(3) // three Mars creatures in play
	})
}
