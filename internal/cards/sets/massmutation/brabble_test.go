package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Brabble
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Destroyed: If it is your turn, your opponent loses 1 Æmber. Otherwise, your opponent loses 3 Æmber.
func TestBrabble(t *testing.T) {
	t.Run("destroyed on your turn drains 1 Æmber", func(t *testing.T) {
		var brabble ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&brabble, Brabble)),
			},
			P2: ct.Side{Amber: 5},
		})

		h.Game().DestroyEach(0, []engine.LocalID{brabble.ID()})

		h.P2.ExpectAmber(4)
	})

	t.Run("destroyed on the opponent's turn drains 3 Æmber", func(t *testing.T) {
		var brabble ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&brabble, Brabble)),
			},
			P2: ct.Side{Amber: 5},
		})
		h.Game().State.ActivePlayer = 1

		h.Game().DestroyEach(1, []engine.LocalID{brabble.ID()})

		h.P2.ExpectAmber(2)
	})
}
