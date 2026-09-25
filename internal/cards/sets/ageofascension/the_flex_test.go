package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// The Flex
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Choose a friendly ready Brobnar creature. Exhaust it. Gain Æmber equal to half its power, rounded down.
func TestTheFlex(t *testing.T) {
	var beefy ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Brobnar,
			Hand:   ct.Cards(TheFlex),
			InPlay: ct.Cards(ct.Bind(&beefy, ct.Creature(ct.Power(5)))),
		},
	})

	h.P1.Play(TheFlex)

	h.Expect(beefy).Exhausted()
	h.P1.ExpectAmber(2)
}
