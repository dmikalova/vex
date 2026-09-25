package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Z-Particle Tracker
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Connected
//	Bonus:  Æmber
//
//	This creature gains, "Fight: Search your deck for an upgrade, reveal it, and put it into your hand. Shuffle your deck."
func TestZParticleTracker(t *testing.T) {
	t.Run("host tutors an upgrade from the deck when it fights", func(t *testing.T) {
		var host, enemy, upgrade, filler ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
						ZParticleTracker,
					),
				),
				Deck: ct.Cards(
					ct.Bind(&upgrade, ct.Upgrade()),
					ct.Bind(&filler, ct.Creature()),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Fight(host, enemy)

		h.Expect(upgrade).At(ct.Hand)
		h.Expect(filler).At(ct.Deck)
	})
}
