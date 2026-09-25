package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Keyfrog
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Beast
//
//	Destroyed: Forge a key at current cost -> purge Keyfrog.
func TestKeyfrog(t *testing.T) {
	t.Run("forges a key at current cost when destroyed", func(t *testing.T) {
		var frog ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&frog, Keyfrog)),
				Amber:  engine.KeyCost,
			},
		})

		h.Game().DestroyEach(0, []engine.LocalID{frog.ID()})

		h.P1.ExpectKeys(1)
		h.P1.ExpectAmber(0)
	})
}
