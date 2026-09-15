package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

func camouflageCanFight(g *engine.Game, attacker, defender ct.Card) bool {
	for _, id := range g.FightTargets(0, attacker.ID()) {
		if id == defender.ID() {
			return true
		}
	}
	return false
}

// Camouflage
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Creatures not on a flank cannot fight this creature.
func TestCamouflage(t *testing.T) {
	var flank, interior, hidden ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			InPlay: ct.Cards(
				ct.Bind(&flank, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4))),
				ct.Bind(&interior, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4))),
				ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4)),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Upgraded(ct.Bind(&hidden, ct.Creature(ct.Power(6))), Camouflage),
			),
		},
	})

	// The interior creature is not on a flank, so it cannot fight the camouflaged
	// creature; a flank creature can.
	if camouflageCanFight(h.Game(), interior, hidden) {
		t.Error("interior creature should not be able to fight the camouflaged creature")
	}
	if !camouflageCanFight(h.Game(), flank, hidden) {
		t.Error("flank creature should be able to fight the camouflaged creature")
	}
	h.P1.Fight(flank, hidden)
	h.Expect(hidden).Damage(4)
}
