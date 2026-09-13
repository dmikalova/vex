package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Redlock
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	Skirmish.
//	At the end of your turn, if you did not play any Creatures this turn, gain 1 Æmber.
func TestRedlock(t *testing.T) {
	var beef ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Shadows,
			InPlay: ct.Cards(ct.Bind(&beef, ct.Creature(ct.OfHouse(card.House.Shadows)))),
			Hand: ct.Cards(
				Redlock,
				ct.Bind(new(ct.Card), ct.Creature(ct.OfHouse(card.House.Shadows))),
			),
		},
	})

	// A turn on which a creature was played: no Æmber at end of turn.
	h.P1.Play(Redlock)
	h.P1.EndTurn()
	h.P1.ExpectAmber(0)
}

// Redlock pays out on a turn its controller played no creatures.
func TestRedlockNoCreaturesPlayed(t *testing.T) {
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Shadows,
			InPlay: ct.Cards(Redlock),
			Hand:   ct.Cards(ct.Tactic(ct.OfHouse(card.House.Shadows))),
		},
	})

	h.P1.EndTurn()
	h.P1.ExpectAmber(1)
}
