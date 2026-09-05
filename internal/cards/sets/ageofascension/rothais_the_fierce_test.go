package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Rothais the Fierce
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Taunt, Hazardous 4.
func TestRothaisTheFierce(t *testing.T) {
	t.Run("deals hazardous damage to an attacker", func(t *testing.T) {
		var attacker, rothais ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&rothais, RothaisTheFierce))},
		})

		h.P1.Fight(attacker, rothais)

		h.Expect(attacker).At(ct.Discard)
	})
}
