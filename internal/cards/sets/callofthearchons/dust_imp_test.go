package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Dust Imp
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Destroyed: Gain 2 Æmber.
func TestDustImp(t *testing.T) {
	t.Run("gains 2 Æmber for its controller when it is destroyed", func(t *testing.T) {
		var imp ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&imp, DustImp)),
			},
		})

		h.Game().DestroyEach(0, []engine.LocalID{imp.ID()})

		h.P1.ExpectAmber(2)
	})
}
