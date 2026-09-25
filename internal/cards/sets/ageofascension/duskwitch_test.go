package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Duskwitch
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Human • Witch
//
//	Omega, Elusive.
//	Your creatures enter play ready.
func TestDuskwitch(t *testing.T) {
	var newbie ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Untamed,
			InPlay: ct.Cards(Duskwitch),
			Hand: ct.Cards(
				ct.Bind(&newbie, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
			),
		},
	})

	h.P1.Play(newbie)

	h.Expect(newbie).Ready()
}
