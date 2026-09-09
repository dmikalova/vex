package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Qincan's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may choose one:
//	- Deal 2 damage to a creature
//	- Attach this creature to Sci. Officer Qincan, and you may archive a creature from play."
func TestQincansBlaster(t *testing.T) {
	t.Run("archives a creature on the attach payoff", func(t *testing.T) {
		var carrier, prey ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&carrier,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						QincansBlaster,
					),
					ct.Bind(&prey, ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(carrier)
		h.P1.ClickOption("Yes")
		h.P1.ExpectPrompt("Choose one")
		h.P1.ClickOption("attach")
		h.P1.ClickCard(prey)

		h.Expect(prey).At(ct.Archives)
	})
}
