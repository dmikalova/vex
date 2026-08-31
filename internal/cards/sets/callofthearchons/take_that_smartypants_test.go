package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Take that, Smartypants
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If there are 3 enemy Logos cards in play, steal 2 Æmber.
func TestTakeThatSmartypants(t *testing.T) {
	t.Run("steals 2 when opponent has 3 or more Logos cards in play", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(TakeThatSmartypants)},
			P2: ct.Side{
				Amber: 4,
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3)),
					ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3)),
					ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3)),
				),
			},
		})

		h.P1.Play(TakeThatSmartypants)

		h.P1.ExpectAmber(3) // 1 bonus + 2 stolen
		h.P2.ExpectAmber(2)
	})

	t.Run("steals nothing when opponent has fewer than 3 Logos cards", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(TakeThatSmartypants)},
			P2: ct.Side{
				Amber: 4,
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3)),
					ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3)),
				),
			},
		})

		h.P1.Play(TakeThatSmartypants)

		h.P1.ExpectAmber(1) // only the bonus
		h.P2.ExpectAmber(4)
	})
}
