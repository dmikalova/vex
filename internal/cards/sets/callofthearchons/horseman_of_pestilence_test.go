package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Horseman of Pestilence
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Horseman • Spirit
//
//	Play/Fight/Reap: Deal 1 damage to each non-Horseman trait creature.
func TestHorsemanOfPestilence(t *testing.T) {
	t.Run("deals 1 damage to each non-Horseman creature when played", func(t *testing.T) {
		var normal, horseman ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(HorsemanOfPestilence)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&normal, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(
						&horseman,
						ct.Creature(
							ct.OfHouse(card.House.Mars),
							ct.Power(3),
							ct.Traits(card.Traits.Horseman),
						),
					),
				),
			},
		})

		h.P1.Play(HorsemanOfPestilence)

		h.Expect(normal).Damage(1)
		h.Expect(horseman).Damage(0)
	})
}
