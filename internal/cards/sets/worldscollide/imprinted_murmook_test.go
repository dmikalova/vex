package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Imprinted Murmook
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Elusive.
//	Your keys cost -1 Æmber.
func TestImprintedMurmook(t *testing.T) {
	t.Run("lowers your key cost by 1 while in play", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(ImprintedMurmook)},
		})

		g := h.Game()
		g.State.Aember[0] = engine.KeyCost - 1
		g.StartTurn(0)

		h.P1.ExpectKeys(1)
	})
}
