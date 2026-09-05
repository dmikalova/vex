package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Way of the Porcupine
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains +3 hazardous.
func TestWayOfThePorcupine(t *testing.T) {
	t.Run("deals hazardous damage to an attacker", func(t *testing.T) {
		var attacker, host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Upgraded(ct.Bind(&host, ct.Creature(ct.Power(3))), WayOfThePorcupine),
				),
			},
		})

		h.P1.Fight(attacker, host)

		h.Expect(attacker).At(ct.Discard)
	})
}
