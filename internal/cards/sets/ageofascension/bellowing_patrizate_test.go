package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bellowing Patrizate
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Giant
//
//	After a Creature enters play, if Bellowing Patrizate is ready, deal 1 damage to it.
func TestBellowingPatrizate(t *testing.T) {
	t.Run("zaps each creature that enters play while ready", func(t *testing.T) {
		var newcomer ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(BellowingPatrizate),
				Hand: ct.Cards(
					ct.Bind(&newcomer, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Play(newcomer)

		h.Expect(newcomer).Damage(1)
	})

	t.Run("does nothing while exhausted", func(t *testing.T) {
		var patrizate, newcomer ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&patrizate, BellowingPatrizate)),
				Hand: ct.Cards(
					ct.Bind(&newcomer, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})
		patrizate.Exhaust()

		h.P1.Play(newcomer)

		h.Expect(newcomer).Damage(0)
	})
}
