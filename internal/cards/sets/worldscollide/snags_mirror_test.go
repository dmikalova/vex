package worldscollide

import (
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
//	Æmber:  1
//	Traits: Item
//
//	After a player chooses an active house, their opponent cannot choose the same
//	house as their active house on their next turn.
func TestSnagsMirror(t *testing.T) {
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(SnagsMirror)},
		P2: ct.Side{House: card.House.Mars},
	})

	h.P1.EndTurn() // the opponent's turn begins; they now pick a house

	h.P2.ChooseHouse(card.House.Sanctum)
	// The mirror reacts to the opponent's choice: P1 is barred from Sanctum next.
	if got := h.Game().State.ForbiddenHouseNext[0].Value; got != card.House.Sanctum {
		t.Fatalf("after P2 chose Sanctum, forbidden-next[P1] = %v, want Sanctum", got)
	}

	h.P2.EndTurn() // back to P1, promoting the forbidden house

	if err := h.Game().ChooseHouse(0, card.House.Sanctum); err != engine.ErrHouseForbidden {
		t.Errorf("P1 choosing Sanctum = %v, want ErrHouseForbidden", err)
	}
	h.P1.ChooseHouse(card.House.Dis) // any other house is allowed
	// The mirror reacts to P1's choice too: P2 is barred from Dis next turn.
	if got := h.Game().State.ForbiddenHouseNext[1].Value; got != card.House.Dis {
		t.Fatalf("after P1 chose Dis, forbidden-next[P2] = %v, want Dis", got)
	}
}
