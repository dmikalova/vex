package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Frane's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may choose one:
//	- Deal 2 damage to a creature
//	- Attach this creature to First Officer Frane, and move all Æmber from First Officer Frane to your pool."
func TestFranesBlaster(t *testing.T) {
	t.Run("attaches to Frane and moves its Æmber to the pool", func(t *testing.T) {
		var carrier, frane ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&carrier,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						FranesBlaster,
					),
					ct.Bind(&frane, FirstOfficerFrane),
				),
			},
		})
		h.Game().AddAmberOn(frane.ID(), 2)

		h.P1.Reap(carrier) // +1 Æmber
		h.P1.ClickOption("Yes")
		h.P1.ExpectPrompt("Choose one")
		h.P1.ClickOption("attach")

		h.Expect(frane).AmberOn(0)
		h.P1.ExpectAmber(3) // 1 reap + 2 moved off Frane
	})
}
