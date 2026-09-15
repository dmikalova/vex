package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Daemo-Knight
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Mutant • Knight
//
//	Destroyed: Steal 1 Æmber.
func TestDaemoKnight(t *testing.T) {
	t.Run("steals 1 Æmber when it is destroyed", func(t *testing.T) {
		var knight ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&knight, DaemoKnight)),
			},
			P2: ct.Side{Amber: 2},
		})

		h.Game().DestroyEach(0, []engine.LocalID{knight.ID()})

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})
}
