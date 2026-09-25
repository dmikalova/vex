package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Molina's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Molina's Blaster to Armsmaster Molina -> deal 3 damage to a creature."
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
					ct.Creature(
						ct.OfHouse(card.House.StarAlliance),
						ct.Named(ArmsmasterMolina.Name),
					),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&victim, ct.Creature(ct.Power(6)))),
			},
		})

		h.P1.Reap(carrier)
		h.P1.ExpectPrompt("Choose one")
		h.P1.ClickOption("attach")
		h.P1.ClickCard(victim)

		h.Expect(victim).Damage(3)
	})
}
