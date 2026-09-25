package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// No Safety in Numbers
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Deal 3 damage to each creature that belongs to a house that has 3 or more creatures in play.
func TestNoSafetyInNumbers(t *testing.T) {
	t.Run(
		"damages every creature of a house with 3+ across both players, sparing smaller houses",
		func(t *testing.T) {
			var mars1, mars2, mars3, untamed ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Shadows,
					Hand:  ct.Cards(NoSafetyInNumbers),
					InPlay: ct.Cards(
						ct.Bind(&mars1, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
						ct.Bind(&mars2, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						// The third Mars creature sits on the opponent's side, so Mars
						// only reaches three by counting both players together.
						ct.Bind(&mars3, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
						// Untamed has a single creature in play, below the threshold.
						ct.Bind(&untamed, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
					),
				},
			})

			h.P1.Play(NoSafetyInNumbers)

			// Mars has three creatures across both players, so all take 3 damage: the
			// two power-3 Mars creatures die, the power-5 one survives damaged.
			h.Expect(mars1).At(ct.Discard)
			h.Expect(mars2).At(ct.Discard)
			h.Expect(mars3).At(ct.PlayArea).Damage(3)

			// Untamed has only one creature in play, so it is untouched.
			h.Expect(untamed).At(ct.PlayArea).Damage(0)
		},
	)
}
