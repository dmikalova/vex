package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Uncharted Lands
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Each Star Alliance Creature gains, "Reap: Move 1 Æmber from Uncharted Lands to your pool."
//	Play: Place 6 Æmber from the common supply on Uncharted Lands.
func TestUnchartedLands(t *testing.T) {
	t.Run("play places 6 Æmber on it", func(t *testing.T) {
		var lands ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(ct.Bind(&lands, UnchartedLands)),
			},
		})

		h.P1.Play(lands)

		h.Expect(lands).AmberOn(6)
	})

	t.Run(
		"a Star Alliance creature reaps to move 1 Æmber from it to your pool",
		func(t *testing.T) {
			var lands, reaper ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.StarAlliance,
					InPlay: ct.Cards(
						ct.Bind(&lands, UnchartedLands),
						ct.Bind(
							&reaper,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(3)),
						),
					),
				},
			})
			h.Game().AddAmberOn(lands.ID(), 6)

			h.P1.Reap(reaper)

			// Base reap gains 1, the granted reap moves 1 more off the artifact.
			h.P1.ExpectAmber(2)
			h.Expect(lands).AmberOn(5)
		},
	)

	t.Run(
		"two copies each grant their own reap, drawing from each artifact",
		func(t *testing.T) {
			var landsA, landsB, reaper ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.StarAlliance,
					InPlay: ct.Cards(
						ct.Bind(&landsA, UnchartedLands),
						ct.Bind(&landsB, UnchartedLands),
						ct.Bind(
							&reaper,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(3)),
						),
					),
				},
			})
			h.Game().AddAmberOn(landsA.ID(), 6)
			h.Game().AddAmberOn(landsB.ID(), 6)

			h.P1.Reap(reaper)

			// Base reap gains 1; each copy's granted reap moves 1 off its own
			// artifact, so both drop to 5 rather than one dropping to 4.
			h.P1.ExpectAmber(3)
			h.Expect(landsA).AmberOn(5)
			h.Expect(landsB).AmberOn(5)
		},
	)
}
