package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Tezmal
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	Reap: Choose a house - your opponent cannot choose that house as their active house on their next turn.
func TestTezmal(t *testing.T) {
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Tezmal)},
		P2: ct.Side{House: card.House.Mars},
	})

	h.P1.Reap(Tezmal)
	h.P1.ExpectPrompt("Choose a house").Source("Tezmal")
	h.P1.ClickOption("Mars")

	if got := h.Game().State.ForbiddenHouseNext[1].Value; got != card.House.Mars {
		t.Fatalf("armed forbidden house = %v, want Mars", got)
	}

	h.P1.EndTurn() // the opponent's turn begins, promoting the forbidden house

	if err := h.Game().ChooseHouse(1, card.House.Mars); err != engine.ErrHouseForbidden {
		t.Errorf("forbidden house = %v, want ErrHouseForbidden", err)
	}
	h.P2.ChooseHouse(card.House.Sanctum) // any other house is allowed
}
