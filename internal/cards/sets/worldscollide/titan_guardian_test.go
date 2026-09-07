package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Titan Guardian
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Armor:  1
//	Traits: Beast • Cyborg
//
//	Taunt.
//	Destroyed: If Titan Guardian is not on a flank, draw 2 cards.
func TestTitanGuardian(t *testing.T) {
	t.Run("draws 2 cards when destroyed off a flank", func(t *testing.T) {
		var guardian, left, right, enemy, d1, d2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2))),
					ct.Bind(&guardian, TitanGuardian),
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2))),
				),
				Deck: ct.Cards(
					ct.Bind(&d1, ct.Creature(ct.Power(1))),
					ct.Bind(&d2, ct.Creature(ct.Power(1))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(10))))},
		})

		h.P1.Fight(guardian, enemy)

		h.Expect(guardian).At(ct.Discard)
		h.Expect(left).At(ct.PlayArea)
		h.Expect(right).At(ct.PlayArea)
		h.Expect(d1).At(ct.Hand)
		h.Expect(d2).At(ct.Hand)
	})

	t.Run("draws no cards when destroyed on a flank", func(t *testing.T) {
		var guardian, enemy, d1 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&guardian, TitanGuardian)),
				Deck:   ct.Cards(ct.Bind(&d1, ct.Creature(ct.Power(1)))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(10))))},
		})

		h.P1.Fight(guardian, enemy)

		h.Expect(guardian).At(ct.Discard)
		h.Expect(d1).At(ct.Deck) // still in the deck, no draw
	})
}
