package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// World Tree
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Put a creature from your discard pile on top of your deck.
func TestWorldTree(t *testing.T) {
	t.Run("returns a creature from the discard pile to the top of the deck", func(t *testing.T) {
		var ghost ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(WorldTree),
				Discard: ct.Cards(
					ct.Bind(&ghost, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4))),
				),
			},
		})

		h.P1.UseAction(WorldTree)

		deck := h.Game().State.Deck[0]
		if deck.Count == 0 || deck.IDs[0] != ghost.ID() {
			t.Errorf("deck top = %v, want Ghost (%d)", deck.IDs[:deck.Count], ghost.ID())
		}
	})
}
