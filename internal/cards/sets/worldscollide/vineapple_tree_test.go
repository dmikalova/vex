package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Vineapple Tree
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Location
//
//	Each player's keys cost +1 Æmber for each growth counter on Vineapple Tree.
//	After a player forges a key, remove each growth counter from Vineapple Tree.
//	Action: Put a growth counter on Vineapple Tree.
func TestVineappleTree(t *testing.T) {
	t.Run("Action places a growth counter that raises each key's cost", func(t *testing.T) {
		var tree ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&tree, VineappleTree)),
			},
		})

		if got := h.Game().CurrentKeyCost(0); got != 6 {
			t.Errorf("key cost with no counters = %d, want 6", got)
		}

		h.P1.UseAction(tree)
		if got := h.Game().CountersOn(tree.ID(), card.Counter.Growth); got != 1 {
			t.Fatalf("growth counters after action = %d, want 1", got)
		}
		if got := h.Game().CurrentKeyCost(0); got != 7 {
			t.Errorf("owner key cost with one counter = %d, want 7", got)
		}
		if got := h.Game().CurrentKeyCost(1); got != 7 {
			t.Errorf("opponent key cost with one counter = %d, want 7", got)
		}
	})

	t.Run("forging a key removes each growth counter", func(t *testing.T) {
		var tree ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&tree, VineappleTree)),
			},
		})
		h.Game().PlaceCounter(tree.ID(), card.Counter.Growth, 3)
		if got := h.Game().CurrentKeyCost(0); got != 9 {
			t.Errorf("key cost with three counters = %d, want 9", got)
		}

		h.Game().ForgeKeyFree(0)

		if got := h.Game().CountersOn(tree.ID(), card.Counter.Growth); got != 0 {
			t.Errorf("growth counters after forge = %d, want 0", got)
		}
		if got := h.Game().CurrentKeyCost(0); got != 6 {
			t.Errorf("key cost after forge = %d, want 6", got)
		}
	})
}
