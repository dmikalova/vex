package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Ember Imp
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Your opponent cannot play more than 2 cards each turn.
func TestEmberImp(t *testing.T) {
	t.Run("bars the opponent from playing more than two cards a turn", func(t *testing.T) {
		var a, b, c ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand: ct.Cards(
					ct.Bind(&a, ct.Tactic(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&b, ct.Tactic(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&c, ct.Tactic(ct.OfHouse(card.House.Brobnar))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(EmberImp)},
		})

		h.P1.Play(a)
		h.P1.Play(b)

		// The third play of the turn is barred.
		if err := h.Game().PlayAction(0, 0); err != engine.ErrCardPlayLimit {
			t.Errorf("third play = %v, want ErrCardPlayLimit", err)
		}
		h.Expect(c).At(ct.Hand)
	})
}
