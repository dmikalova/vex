package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Optio Gorkus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Armor:  3
//	Traits: Dinosaur • Soldier
//
//	Elusive.
//	Each of Optio Gorkus's neighbors gains, "Destroyed: Move all Æmber from this creature to Optio Gorkus."
func TestOptioGorkus(t *testing.T) {
	t.Run("moves a destroyed neighbor's aember to Optio Gorkus", func(t *testing.T) {
		var neighbor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&neighbor, ct.Creature()),
					OptioGorkus,
				),
			},
		})
		h.Game().AddAmberOn(neighbor.ID(), 3)

		h.Game().DestroyEach(0, []engine.LocalID{neighbor.ID()})

		h.Expect(OptioGorkus).AmberOn(3)
	})

	t.Run("moves nothing when the destroyed neighbor has no aember", func(t *testing.T) {
		var neighbor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&neighbor, ct.Creature()),
					OptioGorkus,
				),
			},
		})

		h.Game().DestroyEach(0, []engine.LocalID{neighbor.ID()})

		h.Expect(OptioGorkus).AmberOn(0)
	})
}
