package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Heart of the Forest
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Location
//
//	Each player cannot forge keys while they have more forged keys than their opponent.
func TestHeartOfTheForest(t *testing.T) {
	t.Run("bars a player who leads on keys from forging", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(HeartOfTheForest)},
		})
		g := h.Game()
		g.State.Keys[0] = 1 // P1 leads
		g.State.Aember[0] = 3 * engine.KeyCost
		h.P1.EndTurn() // to P2
		h.P2.EndTurn() // back to P1: forge phase runs, barred while ahead

		h.P1.ExpectKeys(1)
	})

	t.Run("lets a tied player forge", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(HeartOfTheForest)},
		})
		g := h.Game()
		g.State.Aember[0] = 3 * engine.KeyCost
		h.P1.EndTurn()
		h.P2.EndTurn()

		h.P1.ExpectKeys(1)
	})
}
