package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Moor Wolf
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Beast • Wolf
//
//	Skirmish.
//	Play: Ready each other friendly Wolf creature.
func TestMoorWolf(t *testing.T) {
	t.Run("readies each other friendly Wolf creature when played", func(t *testing.T) {
		var wolf ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(MoorWolf),
				InPlay: ct.Cards(
					ct.Bind(&wolf, ct.Creature(
						ct.OfHouse(card.House.Untamed),
						ct.Traits(card.Traits.Wolf),
						ct.Power(3),
					)),
				),
			},
		})
		wolf.Exhaust()

		h.P1.Play(MoorWolf)

		h.Expect(wolf).Ready()
	})
}
