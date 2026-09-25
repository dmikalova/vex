package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Cooperative Hunting
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For each friendly creature in play, deal 1 damage to a creature.
func TestCooperativeHunting(t *testing.T) {
	t.Run(
		"concentrates every instance on one chosen creature",
		func(t *testing.T) {
			var f1, f2 ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Untamed,
					Hand:  ct.Cards(CooperativeHunting),
					InPlay: ct.Cards(
						ct.Bind(&f1, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(5))),
						ct.Bind(&f2, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(5))),
					),
				},
			})

			h.P1.Play(CooperativeHunting)
			// Two friendly creatures, so two instances of 1 damage; both aimed at f1.
			h.P1.ClickCard(f1)
			h.P1.ClickCard(f1)

			h.Expect(f1).Damage(2)
			h.Expect(f2).Damage(0)
		},
	)

	t.Run(
		"chooses a different creature for each instance",
		func(t *testing.T) {
			var f1, f2 ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Untamed,
					Hand:  ct.Cards(CooperativeHunting),
					InPlay: ct.Cards(
						ct.Bind(&f1, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(5))),
						ct.Bind(&f2, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(5))),
					),
				},
			})

			h.P1.Play(CooperativeHunting)
			h.P1.ClickCard(f1)
			h.P1.ClickCard(f2)

			h.Expect(f1).Damage(1)
			h.Expect(f2).Damage(1)
		},
	)
}
