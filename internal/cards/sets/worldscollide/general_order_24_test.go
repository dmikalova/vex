package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// General Order 24
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Law
//
//	At the start of each player's turn, if there are no creatures in play, destroy General Order 24. Otherwise, choose a friendly creature. Destroy each creature of that card's house.

// TestGeneralOrder24DestroysChosenHouse fires the artifact at the start of the
// opponent's turn: they choose one of their creatures and every creature of that
// creature's house is destroyed, while other houses survive.
func TestGeneralOrder24DestroysChosenHouse(t *testing.T) {
	var mars1, mars2, sanctum ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.StarAlliance,
			InPlay: ct.Cards(GeneralOrder24),
		},
		P2: ct.Side{House: card.House.Mars, InPlay: ct.Cards(
			ct.Bind(&mars1, ct.Creature(ct.OfHouse(card.House.Mars))),
			ct.Bind(&mars2, ct.Creature(ct.OfHouse(card.House.Mars))),
			ct.Bind(&sanctum, ct.Creature(ct.OfHouse(card.House.Sanctum))),
		)},
	})

	h.P1.EndTurn()        // P2's turn begins; P2 is the acting player who chooses.
	h.P2.ClickCard(mars1) // choose a Mars creature, so each Mars creature is destroyed.

	h.Expect(mars1).At(ct.Discard)
	h.Expect(mars2).At(ct.Discard)
	h.Expect(sanctum).At(ct.PlayArea) // a different house is untouched.
}

// TestGeneralOrder24SelfDestructsWithoutCreatures covers the "instead" clause: with
// no friendly creatures to choose, the artifact destroys itself.
func TestGeneralOrder24SelfDestructsWithoutCreatures(t *testing.T) {
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.StarAlliance,
			InPlay: ct.Cards(GeneralOrder24),
		},
		P2: ct.Side{House: card.House.Mars},
	})

	h.P1.EndTurn() // P2 controls no creatures, so General Order 24 destroys itself.

	h.Expect(GeneralOrder24).At(ct.Discard)
}
