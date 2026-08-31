package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Clear Mind
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Unstun each friendly creature.
func TestClearMind(t *testing.T) {
	t.Run("unstuns each friendly creature but not the opponent's", func(t *testing.T) {
		var ally1, ally2, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Bind(&ally1, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
					ct.Bind(&ally2, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
				),
				Hand: ct.Cards(ClearMind),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
				),
			},
		})
		ally1.Stun()
		ally2.Stun()
		foe.Stun()

		h.P1.Play(ClearMind)

		h.Expect(ally1).Stunned(false)
		h.Expect(ally2).Stunned(false)
		h.Expect(foe).Stunned(true) // enemy creatures are unaffected
	})
}
