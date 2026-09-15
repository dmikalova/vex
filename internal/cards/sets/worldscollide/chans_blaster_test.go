package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Chan's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Chan's Blaster to Commander Chan -> use another creature."
func TestChansBlaster(t *testing.T) {
	t.Run("deals 2 damage with the blaster", func(t *testing.T) {
		var carrier, victim ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&carrier,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						ChansBlaster,
					),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&victim, ct.Creature(ct.Power(5)))),
			},
		})

		h.P1.Reap(carrier)
		h.P1.ExpectPrompt("Choose one")
		h.P1.ClickOption("deal 2")
		h.P1.ClickCard(victim)

		h.Expect(victim).Damage(2)
	})
}
