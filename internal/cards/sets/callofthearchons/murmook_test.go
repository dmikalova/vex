package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Murmook
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Your opponent's keys cost +1 Æmber.
func TestMurmook(t *testing.T) {
	t.Run("raises the opponent's key cost by 1 while in play", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(Murmook)},
		})

		g := h.Game()
		g.State.Aember[1] = engine.KeyCost
		g.BeginTurn(1)
		h.P2.ExpectKeys(0) // one Æmber short of the raised cost

		g.State.Aember[1] = engine.KeyCost + 1
		g.BeginTurn(1)
		h.P2.ExpectKeys(1)
	})
}
