package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Standardized Testing
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy each creature with the lowest or highest power.
func TestStandardizedTesting(t *testing.T) {
	var low, mid, high ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Logos,
			Hand:  ct.Cards(StandardizedTesting),
			InPlay: ct.Cards(
				ct.Bind(&mid, ct.Creature(ct.Power(4))),
			),
		},
		P2: ct.Side{InPlay: ct.Cards(
			ct.Bind(&low, ct.Creature(ct.Power(2))),
			ct.Bind(&high, ct.Creature(ct.Power(6))),
		)},
	})

	h.P1.Play(StandardizedTesting)

	h.Expect(low).At(ct.Discard)
	h.Expect(high).At(ct.Discard)
	h.Expect(mid).At(ct.PlayArea)
}
