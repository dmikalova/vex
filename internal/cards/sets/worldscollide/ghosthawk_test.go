package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
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
//	Play: You may reap with up to 2 different neighboring Creatures, one at a time.
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
	h.P1.ClickOption("Yes")     // accept the May

	// Reap with each neighbor, one at a time.
	h.P1.ClickCard(left)
	h.P1.ClickCard(right)

	h.Expect(left).At(ct.PlayArea).Exhausted()
	h.Expect(right).At(ct.PlayArea).Exhausted()
	if got := h.Game().Aember(0); got != 2 {
		t.Errorf("aember = %d, want 2 (each neighbor reaped for 1)", got)
	}
}
