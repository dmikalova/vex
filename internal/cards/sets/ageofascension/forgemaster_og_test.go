package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Forgemaster Og
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Giant
//
//	After a player forges a key, that player loses all their Æmber.
func TestForgemasterOg(t *testing.T) {
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Brobnar,
			InPlay: ct.Cards(ForgemasterOg),
		},
	})
	g := h.Game()

	// Player 1 starts with enough Æmber to forge plus a surplus.
	g.State.ForgeCanonicalKeys(0, 0)
	g.State.Aember[0] = engine.KeyCost + 4
	h.P1.EndTurn() // to P2
	h.P2.EndTurn() // back to P1: forge phase forges a key, then Og drains the rest

	h.P1.ExpectKeys(1)
	h.P1.ExpectAmber(0)
}
