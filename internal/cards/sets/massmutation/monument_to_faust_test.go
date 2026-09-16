package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Monument to Faust
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: If Faust the Great is in your discard pile, keys cost +2 Æmber during your opponent's next turn. Otherwise, keys cost +1 Æmber during your opponent's next turn.
func TestMonumentToFaust(t *testing.T) {
	t.Run("raises the opponent's key cost by 1", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(MonumentToFaust),
			},
			P2: ct.Side{House: card.House.Saurian},
		})

		h.P1.UseAction(MonumentToFaust)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Saurian)

		if got := h.Game().CurrentKeyCost(1); got != 7 {
			t.Errorf("key cost = %d, want 7 (6 + 1)", got)
		}
	})

	t.Run("raises by 2 when Faust the Great is in your discard pile", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Saurian,
				InPlay:  ct.Cards(MonumentToFaust),
				Discard: ct.Cards(FaustTheGreat),
			},
			P2: ct.Side{House: card.House.Saurian},
		})

		h.P1.UseAction(MonumentToFaust)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Saurian)

		if got := h.Game().CurrentKeyCost(1); got != 8 {
			t.Errorf("key cost = %d, want 8 (6 + 2)", got)
		}
	})
}
