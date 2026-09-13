package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Collector Worm
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Armor:  5
//	Traits: Beast
//
//	Fight: Put the Creature Collector Worm fought into your archives.
func TestCollectorWorm(t *testing.T) {
	var worm, prey ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Mars,
			InPlay: ct.Cards(ct.Bind(&worm, CollectorWorm)),
		},
		P2: ct.Side{
			InPlay: ct.Cards(ct.Bind(&prey, ct.Creature(ct.Power(4)))),
		},
	})

	h.P1.Fight(worm, prey)

	h.Expect(prey).At(ct.Archives)
}
