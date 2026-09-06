package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mindworm
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Beast
//
//	Elusive.
//	Before Fight: Deal damage equal to its power to each neighbor of the creature Mindworm fights.
func TestMindworm(t *testing.T) {
	var worm, left, target, right ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Mars,
			InPlay: ct.Cards(ct.Bind(&worm, Mindworm)),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.Power(5))),
				ct.Bind(&target, ct.Creature(ct.Power(3))),
				ct.Bind(&right, ct.Creature(ct.Power(5))),
			),
		},
	})

	h.P1.Fight(worm, target)

	h.Expect(left).Damage(3)
	h.Expect(right).Damage(3)
}
