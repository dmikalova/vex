package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Yzphyz Knowdrone
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Armor:  1
//	Traits: Martian • Scientist
//
//	Play: Archive a card from your hand. You may purge a card from your archives to stun a creature.
func TestYzphyzKnowdrone(t *testing.T) {
	var spare, victim ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Mars,
			Hand: ct.Cards(
				YzphyzKnowdrone,
				ct.Bind(&spare, ct.Tactic()),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(ct.Bind(&victim, ct.Creature(ct.Power(3)))),
		},
	})

	h.P1.Play(YzphyzKnowdrone)
	h.P1.ClickCard(spare)  // archive the spare card from hand
	h.P1.ClickCard(victim) // stun a creature

	h.Expect(spare).At(ct.Purge)
	h.Expect(victim).Stunned(true)
}
