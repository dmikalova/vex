package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Aember Conduction Unit
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	After an enemy creature reaps, if this is the first time a creature has reaped this turn, stun it.
func TestAemberConductionUnit(t *testing.T) {
	var reaper, second ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Brobnar,
			InPlay: ct.Cards(
				ct.Bind(&reaper, ct.Creature(ct.Power(3))),
				ct.Bind(&second, ct.Creature(ct.Power(3))),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(AemberConductionUnit),
		},
	})

	h.P1.Reap(reaper)
	h.Expect(reaper).Stunned(true)

	// The second reap of the turn is not the first, so it is not stunned.
	h.P1.Reap(second)
	h.Expect(second).Stunned(false)
}
