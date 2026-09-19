package worldscollide

import (
	"errors"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Snag
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Demon
//
//	Fight: Your opponent must choose the house of the creature Snag fights as their active house during their next turn.
func TestSnag(t *testing.T) {
	var snag, foe ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Dis,
			InPlay: ct.Cards(ct.Bind(&snag, Snag)),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3), ct.Armor(6))),
			),
		},
	})

	h.P1.Fight(snag, foe)

	if h.Game().State.HouseConstraintCountNext[1] != 1 {
		t.Fatalf("armed constraints = %d, want 1", h.Game().State.HouseConstraintCountNext[1])
	}

	h.P1.EndTurn() // the opponent's turn begins, promoting the must

	// The must reads the fought creature's house live at choice time (Logos), so
	// any other house is rejected and Logos is required.
	if err := h.Game().ChooseHouse(1, card.House.Mars); !errors.Is(err, engine.ErrHouseNotAllowed) {
		t.Errorf("a house other than the fought creature's = %v, want ErrHouseNotAllowed", err)
	}
	h.P2.ChooseHouse(card.House.Logos)
}
