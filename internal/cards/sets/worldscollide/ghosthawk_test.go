package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Ghosthawk
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Beast
//
//	Deploy.
//	Play: Reap with each of Ghosthawk's neighbors, one at a time.
func TestGhosthawk(t *testing.T) {
	untamed := ct.OfHouse(card.House.Untamed)

	var left, right ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(untamed)),
				ct.Bind(&right, ct.Creature(untamed)),
			),
			Hand: ct.Cards(Ghosthawk),
		},
	})

	h.P1.Play(Ghosthawk)
	h.P1.ClickOption("Between") // deploy between the two neighbors

	// Reap with each neighbor, one at a time: pick the first, the last is automatic.
	h.P1.ClickCard(left)

	h.Expect(left).At(ct.PlayArea).Exhausted()
	h.Expect(right).At(ct.PlayArea).Exhausted()
	if got := h.Game().Aember(0); got != 2 {
		t.Errorf("aember = %d, want 2 (each neighbor reaped for 1)", got)
	}
}
