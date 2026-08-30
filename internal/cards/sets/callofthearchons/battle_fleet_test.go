package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Battle Fleet
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Reveal any number of Mars cards from your hand, and for each card revealed this way, draw a card.
func TestBattleFleet(t *testing.T) {
	t.Run("draws a card for each Mars card left in hand", func(t *testing.T) {
		var d1, d2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand: ct.Cards(
					BattleFleet,
					ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2)),
					ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2)),
					ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(2)), // non-Mars, not counted
				),
				Deck: ct.Cards(
					ct.Bind(&d1, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
					ct.Bind(&d2, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
			},
		})

		h.P1.Play(BattleFleet) // two Mars cards remain in hand -> draw 2

		h.Expect(d1).At(ct.Hand)
		h.Expect(d2).At(ct.Hand)
	})
}
