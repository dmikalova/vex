package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Po's Pixies
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Faerie
//
//	Elusive.
//	Æmber stolen or captured from your pool is taken from the common supply instead.
func TestPosPixies(t *testing.T) {
	t.Run(
		"a theft from the Pixies controller's pool is drawn from the common supply",
		func(t *testing.T) {
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Untamed,
					InPlay: ct.Cards(PosPixies),
					Amber:  3,
				},
				P2: ct.Side{},
			})
			g := h.Game()
			// The opponent steals 2 from the Pixies controller; the Pixies keep the pool
			// intact and the 2 come from the common supply instead.
			engine.StealAember{
				Amount: 2,
				Player: engine.Opponent,
			}.Resolve(&engine.EffectContext{
				Resolver:   g,
				Controller: 0,
			})
			if got := g.State.Aember[0]; got != 3 {
				t.Errorf("Pixies controller pool = %d, want 3 (kept)", got)
			}
			if got := g.State.Aember[1]; got != 2 {
				t.Errorf("thief pool = %d, want 2 (from supply)", got)
			}
		},
	)
}
