package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Daemo-Thief
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Mutant • Thief
//
//	Elusive.
//	Destroyed: Steal 1 Æmber.
func TestDaemoThief(t *testing.T) {
	t.Run("steals 1 aember when it is destroyed", func(t *testing.T) {
		var thief ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(ct.Bind(&thief, DaemoThief))},
			P2: ct.Side{Amber: 2},
		})

		h.Game().DestroyEach(0, []engine.LocalID{thief.ID()})

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})
}
