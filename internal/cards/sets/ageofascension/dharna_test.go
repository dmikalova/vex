package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dharna
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Witch
//
//	Play: For each friendly damaged creature in play, gain 1 Æmber.
//	Reap: Heal 2 damage from a friendly creature.
func TestDharna(t *testing.T) {
	t.Run("gains 1 aember for each damaged friendly creature when played", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(Dharna),
				InPlay: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(6))),
				),
			},
		})
		ally.Damaged(2)

		h.P1.Play(Dharna)

		h.P1.ExpectAmber(1)
	})
}
