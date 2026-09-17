package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lady Loreena
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Armor:  3
//	Traits: Spirit • Knight
//
//	Taunt.
//	Lady Loreena's taunt also applies to its neighbors' neighbors.
func TestLadyLoreena(t *testing.T) {
	t.Run("shields her neighbors and her neighbors' neighbors", func(t *testing.T) {
		var loreena, near, far ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Creature()),
			},
			P2: ct.Side{
				House: card.House.Sanctum,
				// Line: [far, near, loreena] — far is the neighbor's neighbor.
				InPlay: ct.Cards(
					ct.Bind(&far, ct.Creature()),
					ct.Bind(&near, ct.Creature()),
					ct.Bind(&loreena, LadyLoreena),
				),
			},
		})

		// Both the neighbor and the neighbor's neighbor are shielded.
		if !h.Game().TauntShielded(near.ID()) {
			t.Error("a neighbor of Lady Loreena should be shielded")
		}
		if !h.Game().TauntShielded(far.ID()) {
			t.Error("a neighbor's neighbor of Lady Loreena should be shielded")
		}
	})
}
