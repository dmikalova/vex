package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Golden Spiral
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Common
//	Traits: Location
//
//	Action: Choose a friendly creature. Exalt the chosen creature. Ready and use the chosen creature.
func TestTheGoldenSpiral(t *testing.T) {
	t.Run("exalts a friendly creature, then readies and uses it", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					TheGoldenSpiral,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(4))),
				),
			},
			P2: ct.Side{},
		})

		h.P1.UseAction(TheGoldenSpiral)

		// The chosen creature is exalted, then reaps (no enemy to fight).
		h.Expect(ally).AmberOn(1)
		h.Expect(ally).Exhausted()
		h.P1.ExpectAmber(1)
	})
}
