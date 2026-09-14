package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Silver Key Imp
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	Players cannot forge their second key.
func TestSilverKeyImp(t *testing.T) {
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(SilverKeyImp)},
	})
	g := h.Game()

	g.State.Keys[0] = 1 // has one key; the second is barred
	g.State.Aember[0] = engine.KeyCost
	h.P1.EndTurn() // to P2
	h.P2.EndTurn() // back to P1: forge phase runs, second key is barred

	h.P1.ExpectKeys(1)
}
