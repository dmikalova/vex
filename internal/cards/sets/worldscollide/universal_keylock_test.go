package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Universal Keylock
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	Each player's keys cost +3 Æmber.
//	After a player forges a key, destroy Universal Keylock.
func TestUniversalKeylock(t *testing.T) {
	t.Run("raises each player's key cost by 3", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, InPlay: ct.Cards(UniversalKeylock)},
		})

		g := h.Game()
		g.State.Aember[0] = engine.KeyCost + 2
		g.StartTurn(0)

		h.P1.ExpectKeys(0) // needs KeyCost+3, only has KeyCost+2
	})

	t.Run("is destroyed after a player forges a key", func(t *testing.T) {
		var keylock ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&keylock, UniversalKeylock)),
			},
		})

		g := h.Game()
		g.State.Aember[0] = engine.KeyCost + 3
		g.StartTurn(0)

		h.P1.ExpectKeys(1)
		h.Expect(keylock).At(ct.Discard)
	})
}
