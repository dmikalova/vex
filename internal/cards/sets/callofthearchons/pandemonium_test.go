package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Pandemonium
func TestPandemonium(t *testing.T) {
	var mine, hurt, theirs ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Dis,
			Hand:  ct.Cards(Pandemonium),
			InPlay: ct.Cards(
				ct.Bind(&mine, ct.Creature(ct.Power(4))),
				ct.Bind(&hurt, ct.Creature(ct.Power(4))),
			),
			Amber: 3,
		},
		P2: ct.Side{
			InPlay: ct.Cards(ct.Bind(&theirs, ct.Creature(ct.Power(5)))),
			Amber:  3,
		},
	})
	hurt.Damaged(1)

	h.P1.Play(Pandemonium)

	// Each side's creatures reach across the board, so both pools are drained.
	h.Expect(mine).AmberOn(1)
	h.Expect(theirs).AmberOn(1)
	h.Expect(hurt).AmberOn(0) // damaged, so it captures nothing
	h.P1.ExpectAmber(3)       // 3 + 1 bonus from Pandemonium - 1 to the enemy creature
	h.P2.ExpectAmber(2)
}
