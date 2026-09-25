package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Odoac the Patrician
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Dinosaur • Politician
//
//	While Odoac the Patrician has Æmber on it, your Æmber cannot be stolen.
//	Play: Odoac the Patrician captures 1 Æmber from your opponent.
func TestOdoacThePatrician(t *testing.T) {
	var odoac ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Saurian,
			Hand:  ct.Cards(ct.Bind(&odoac, OdoacThePatrician)),
		},
		P2: ct.Side{Amber: 3},
	})

	h.P1.Play(odoac)

	// Play: captures 1 Æmber from the opponent onto itself.
	h.Expect(odoac).AmberOn(1)
	h.P2.ExpectAmber(2)
}
