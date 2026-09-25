package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Curia Saurus
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Location
//
//	Each creature with Æmber on it gains, "Destroyed: Move 1 Æmber from this creature to the most powerful enemy creature."
func TestCuriaSaurus(t *testing.T) {
	t.Run(
		"an Æmber-bearing creature moves 1 to the most powerful enemy when destroyed",
		func(t *testing.T) {
			var bearer, enemyStrong, enemyWeak ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Saurian,
					InPlay: ct.Cards(
						CuriaSaurus,
						ct.Bind(&bearer, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3))),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(&enemyStrong, ct.Creature(ct.Power(6))),
						ct.Bind(&enemyWeak, ct.Creature(ct.Power(2))),
					),
				},
			})
			h.Game().State.Cards[bearer.ID()].Amber = 1

			h.Game().DestroyEach(0, []engine.LocalID{bearer.ID()})

			h.Expect(enemyStrong).AmberOn(1)
			h.Expect(enemyWeak).AmberOn(0)
		},
	)

	t.Run("a creature with no Æmber grants no move when destroyed", func(t *testing.T) {
		var plain, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					CuriaSaurus,
					ct.Bind(&plain, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6)))),
			},
		})

		h.Game().DestroyEach(0, []engine.LocalID{plain.ID()})

		h.Expect(enemy).AmberOn(0)
	})
}
