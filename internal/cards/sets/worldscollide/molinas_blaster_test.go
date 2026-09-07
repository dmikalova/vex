package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Molina's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may choose one:
//	- Deal 2 damage to a creature
//	- Attach this creature to Armsmaster Molina, and you may deal 3 damage to a creature."
func TestMolinasBlaster(t *testing.T) {
	t.Run("deals 3 damage on the attach payoff", func(t *testing.T) {
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
						MolinasBlaster,
					),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&victim, ct.Creature(ct.Power(6)))),
			},
		})

		h.P1.Reap(carrier)
		h.P1.ClickOption("Yes")
		h.P1.ExpectPrompt("Choose one")
		h.P1.ClickOption("attach")
		h.P1.ClickCard(victim)

		h.Expect(victim).Damage(3)
	})
}
