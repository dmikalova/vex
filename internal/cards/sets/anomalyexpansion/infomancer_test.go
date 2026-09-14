package anomalyexpansion

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Infomancer
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Traits: Human • Cyborg
//
//	Elusive.
//	Play: Put a Tactic card from your hand faceup under Infomancer.
//	Reap: Trigger the play effect of a Tactic grafted onto Infomancer.
func TestInfomancer(t *testing.T) {
	// Playing Infomancer grafts a Tactic from hand, and reaping triggers that
	// grafted Tactic's play effect while it stays grafted under Infomancer.
	var info, refrain, foe ct.Card
	// Warriors' Refrain (Play: stun each creature with power 3 or lower) stands in
	// for a real Tactic to graft; it is built here rather than imported from its
	// own set so this test stays self-contained.
	refrainDef := card.Build(
		"Warriors' Refrain",
		card.House.Brobnar,
		card.Type.Tactic,
		card.Rarity.Common,
		card.WithAemberBonus(1),
		card.WithAbility(
			card.Trigger.Play, card.Stun{
				Target: card.Target.EachCreature.PowerAtMost(3),
			}),
	)
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Brobnar,
			Hand: ct.Cards(
				ct.Bind(&info, Infomancer),
				ct.Bind(&refrain, refrainDef),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3)))),
		},
	})

	h.P1.Play(Infomancer)
	// The graft is the only Tactic in hand, so it is placed under Infomancer.
	h.Expect(refrain).At(ct.Under)

	// A creature enters play exhausted, so ready Infomancer before it reaps.
	info.Ready()
	h.P1.Reap(Infomancer)

	// Warriors' Refrain's play effect stunned the low-power enemy, and the Tactic
	// stayed grafted rather than moving to play or the discard pile.
	h.Expect(foe).Stunned(true)
	h.Expect(refrain).At(ct.Under)
	h.P1.ExpectAmber(1) // from the reap; a triggered Tactic grants no Æmber bonus
}
