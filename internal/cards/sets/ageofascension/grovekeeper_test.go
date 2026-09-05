package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Grovekeeper
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Witch
//
//	At the end of your turn, give each neighboring creature a +1 power counter.
func TestGrovekeeper(t *testing.T) {
	t.Run("adds a power counter to each neighbor at the end of the turn", func(t *testing.T) {
		var left, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4))),
					Grovekeeper,
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4))),
				),
			},
			P2: ct.Side{},
		})

		h.P1.EndTurn()

		h.Expect(left).Power(5)
		h.Expect(right).Power(5)
	})
}
