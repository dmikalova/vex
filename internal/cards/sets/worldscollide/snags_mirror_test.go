package worldscollide

import (
	"errors"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Snag's Mirror
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	After a player chooses an active house, their opponent cannot choose the same house as their active house during their next turn.
func TestSnagsMirror(t *testing.T) {
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(SnagsMirror)},
		P2: ct.Side{House: card.House.Mars},
	})

	h.P1.EndTurn() // the opponent's turn begins; they now pick a house

	h.P2.ChooseHouse(card.House.Sanctum)
	// The mirror reacts to the opponent's choice: P1 is barred from Sanctum next.
	if got := h.Game().State.HouseConstraintsNext[0]; h.Game().State.HouseConstraintCountNext[0] != 1 ||
		got[0].House != card.House.Sanctum {
		t.Fatalf(
			"after P2 chose Sanctum, armed constraint[P1] = %+v, want a cannot on Sanctum",
			got[0],
		)
	}

	h.P2.EndTurn() // back to P1, promoting the forbidden house

	if err := h.Game().
		ChooseHouse(0, card.House.Sanctum); !errors.Is(
		err,
		engine.ErrHouseNotAllowed,
	) {
		t.Errorf("P1 choosing Sanctum = %v, want ErrHouseNotAllowed", err)
	}
	h.P1.ChooseHouse(card.House.Dis) // any other house is allowed
	// The mirror reacts to P1's choice too: P2 is barred from Dis next turn.
	if got := h.Game().State.HouseConstraintsNext[1]; h.Game().State.HouseConstraintCountNext[1] != 1 ||
		got[0].House != card.House.Dis {
		t.Fatalf("after P1 chose Dis, armed constraint[P2] = %+v, want a cannot on Dis", got[0])
	}
}
