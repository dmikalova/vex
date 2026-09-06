package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Eureka!
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Alpha.
//	Play: Gain 2 Æmber. Archive 2 random cards from your hand.
func TestEureka(t *testing.T) {
	var a, b ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Logos,
			Hand: ct.Cards(
				Eureka,
				ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Logos))),
				ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Logos))),
			),
		},
	})

	h.P1.Play(Eureka)

	// One pip from Eureka plus 2 gained.
	h.P1.ExpectAmber(3)
	// Both other hand cards were archived.
	h.Expect(a).At(ct.Archives)
	h.Expect(b).At(ct.Archives)
}
