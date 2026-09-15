package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Blypyp
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Martian • Scientist
//
//	Reap: The next Mars creature you play this turn enters play ready.
func TestBlypyp(t *testing.T) {
	var blypyp, next ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Mars,
			InPlay: ct.Cards(ct.Bind(&blypyp, Blypyp)),
			Hand:   ct.Cards(ct.Bind(&next, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3)))),
		},
	})

	h.P1.Reap(blypyp)
	h.P1.Play(next)

	h.Expect(next).At(ct.PlayArea).Ready()
}
