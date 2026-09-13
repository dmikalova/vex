package worldscollide

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
	// Playing Infomancer grafts an action from hand, and reaping triggers that
	// grafted action's play effect while it stays grafted under Infomancer.
	var info, brew ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Brobnar,
			Hand: ct.Cards(
				ct.Bind(&info, Infomancer),
				ct.Bind(&brew, ChieftainsBrew),
			),
		},
	})

	h.P1.Play(Infomancer)
	// The graft is the only action in hand, so it is placed under Infomancer.
	h.Expect(brew).At(ct.Under)

	// A creature enters play exhausted, so ready Infomancer before it reaps.
	info.Ready()
	h.P1.Reap(Infomancer)

	// Chieftain's Brew's play effect gave Infomancer two +1 power counters, and the
	// action stayed grafted rather than moving to play or the discard pile.
	h.Expect(Infomancer).Power(5)
	h.Expect(brew).At(ct.Under)
	h.P1.ExpectAmber(1) // from the reap; a triggered action grants no Æmber bonus
}
