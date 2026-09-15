package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Daemo-Beast
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant • Beast
//
//	Skirmish.
//	Destroyed: Steal 1 Æmber.
func TestDaemoBeast(t *testing.T) {
	t.Run("steals 1 aember when it is destroyed", func(t *testing.T) {
		var beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(ct.Bind(&beast, DaemoBeast))},
			P2: ct.Side{Amber: 2},
		})

		h.Game().DestroyEach(0, []engine.LocalID{beast.ID()})

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})
}
