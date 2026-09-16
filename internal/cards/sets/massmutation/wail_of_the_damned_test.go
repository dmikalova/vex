package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Wail of the Damned
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy a creature with no bonus icons.
//	Enhance Capture.
func TestWailOfTheDamned(t *testing.T) {
	var wail, bare, iconed ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Dis,
			Hand:  ct.Cards(ct.Bind(&wail, WailOfTheDamned)),
		},
		P2: ct.Side{InPlay: ct.Cards(
			ct.Bind(&bare, ct.Creature(ct.Power(3))),
			ct.Bind(&iconed, ct.Creature(ct.Power(3), ct.Bonus(card.Bonus.Aember))),
		)},
	})

	h.P1.Play(wail)

	h.Expect(bare).At(ct.Discard)
	h.Expect(iconed).At(ct.PlayArea)
}
