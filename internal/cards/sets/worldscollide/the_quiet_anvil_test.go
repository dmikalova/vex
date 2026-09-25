package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// The Quiet Anvil
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	Each player's keys cost -2 Æmber.
//	After a player forges a key, destroy The Quiet Anvil.
func TestTheQuietAnvil(t *testing.T) {
	t.Run("lowers each player's key cost by 2", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(TheQuietAnvil),
			},
		})

		g := h.Game()
		g.State.Aember[0] = engine.KeyCost - 2
		g.StartTurn(0)

		h.P1.ExpectKeys(1)
	})

	t.Run("is destroyed after a player forges a key", func(t *testing.T) {
		var anvil ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&anvil, TheQuietAnvil)),
			},
		})

		g := h.Game()
		g.State.Aember[0] = engine.KeyCost - 2
		g.StartTurn(0)

		h.P1.ExpectKeys(1)
		h.Expect(anvil).At(ct.Discard)
	})
}
