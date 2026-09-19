package ageofascension

import (
	"errors"
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
//	Reap: Choose a house. Your opponent cannot choose that house as their active house during their next turn.
func TestTezmal(t *testing.T) {
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Tezmal)},
		P2: ct.Side{House: card.House.Mars},
	})

	h.P1.Reap(Tezmal)
	h.P1.ExpectPrompt("Choose a house").Source("Tezmal")
	h.P1.ClickOption("Mars")

	if got := h.Game().State.HouseConstraintsNext[1]; h.Game().State.HouseConstraintCountNext[1] != 1 ||
		got[0].House != card.House.Mars {
		t.Fatalf("armed constraint = %+v (count %d), want one on Mars",
			got[0], h.Game().State.HouseConstraintCountNext[1])
	}

	h.P1.EndTurn() // the opponent's turn begins, promoting the forbidden house

	if err := h.Game().ChooseHouse(1, card.House.Mars); !errors.Is(err, engine.ErrHouseNotAllowed) {
		t.Errorf("forbidden house = %v, want ErrHouseNotAllowed", err)
	}
	h.P2.ChooseHouse(card.House.Sanctum) // any other house is allowed
}
