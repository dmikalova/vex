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
	t.Run("draws a card for each Mars card revealed from hand", func(t *testing.T) {
		var d1, d2, m1, m2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand: ct.Cards(
					BattleFleet,
					ct.Bind(&m1, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
					ct.Bind(&m2, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
					ct.Creature(
						ct.OfHouse(card.House.Brobnar),
						ct.Power(2),
					), // non-Mars, not revealable
				),
				Deck: ct.Cards(
					ct.Bind(&d1, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
					ct.Bind(&d2, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
			},
		})

		h.P1.Play(BattleFleet)
		// Revealing is the player's choice; showing both Mars cards draws 2.
		h.P1.ClickCard(m1)
		h.P1.ClickCard(m2)

		h.Expect(d1).At(ct.Hand)
		h.Expect(d2).At(ct.Hand)
	})

	t.Run("reveals no card when the player declines", func(t *testing.T) {
		var d1 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand: ct.Cards(
					BattleFleet,
					ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2)),
				),
				Deck: ct.Cards(
					ct.Bind(&d1, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
			},
		})

		h.P1.Play(BattleFleet)
		h.P1.ClickDone()

		h.Expect(d1).At(ct.Deck)
	})
}
