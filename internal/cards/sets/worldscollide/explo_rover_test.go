package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Explo-rover
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Robot
//
//	Skirmish.
//	Explo-rover may be played as an Upgrade instead of a Creature, with the text: "This Creature gains skirmish."
func TestExploRover(t *testing.T) {
	t.Run("played as a creature deals no retaliation damage when it fights", func(t *testing.T) {
		var rover, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&rover, ExploRover)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(5)))),
			},
		})

		h.P1.Fight(rover, foe)

		h.Expect(rover).At(ct.PlayArea).Damage(0) // Skirmish: no return damage
		h.Expect(foe).Damage(3)
	})

	t.Run("played as an upgrade grants its host skirmish", func(t *testing.T) {
		var host, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(ExploRover),
				InPlay: ct.Cards(
					ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(5)))),
			},
		})

		h.P1.Play(ExploRover)
		h.P1.ClickOption("Upgrade")
		h.P1.ClickCard(host) // attach to the friendly host, not the enemy

		h.P1.Fight(host, foe)

		h.Expect(host).At(ct.PlayArea).Damage(0) // host now has Skirmish
		h.Expect(foe).Damage(4)
	})
}
