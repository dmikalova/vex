package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lollop the Titanic
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  11
//	Traits: Giant • Location
//
//	Lollop the Titanic deals no damage when attacked.
func TestLollopTheTitanic(t *testing.T) {
	var attacker ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Brobnar,
			InPlay: ct.Cards(ct.Bind(&attacker, ct.Creature(ct.Power(6)))),
		},
		P2: ct.Side{InPlay: ct.Cards(LollopTheTitanic)},
	})

	h.P1.Fight(attacker, LollopTheTitanic)

	// Lollop deals no retaliation damage, so the attacker is unharmed; Lollop still
	// takes the attacker's fight damage.
	h.Expect(attacker).Damage(0)
	h.Expect(LollopTheTitanic).Damage(6)
}
