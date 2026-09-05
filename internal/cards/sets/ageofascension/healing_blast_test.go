package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Healing Blast
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Fully heal a creature. If you healed 4 or more damage, gain 2 Æmber.
func TestHealingBlast(t *testing.T) {
	t.Run(
		"fully heals a creature and gains aember when at least 4 damage is healed",
		func(t *testing.T) {
			var ally ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Sanctum,
					Hand:  ct.Cards(HealingBlast),
					InPlay: ct.Cards(
						ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(6))),
					),
				},
			})
			ally.Damaged(4)

			h.P1.Play(HealingBlast)

			h.Expect(ally).Damage(0)
			h.P1.ExpectAmber(3)
		},
	)
}
