package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Garcia's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This Creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a Creature
//	- Attach Garcia's Blaster to Sensor Chief Garcia -> steal 1 Æmber."
func TestGarciasBlaster(t *testing.T) {
	t.Run("attaches to Garcia and steals 1 Æmber", func(t *testing.T) {
		var carrier ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&carrier,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						GarciasBlaster,
					),
					SensorChiefGarcia,
				),
			},
			P2: ct.Side{Amber: 2},
		})

		h.P1.Reap(carrier) // +1 Æmber
		h.P1.ExpectPrompt("Choose one")
		h.P1.ClickOption("attach")

		h.P1.ExpectAmber(2) // 1 reap + 1 stolen
		h.P2.ExpectAmber(1) // 2 - 1 stolen
	})
}
