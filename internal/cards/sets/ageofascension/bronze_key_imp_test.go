package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Bronze Key Imp
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	Players cannot forge their first key.
func TestBronzeKeyImp(t *testing.T) {
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Dis,
			InPlay: ct.Cards(BronzeKeyImp),
		},
	})
	g := h.Game()

	g.State.ForgeCanonicalKeys(0, 0)
	g.State.Aember[0] = engine.KeyCost
	h.P1.EndTurn() // to P2
	h.P2.EndTurn() // back to P1: forge phase runs, first key is barred

	h.P1.ExpectKeys(0)
}
