package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Take Hostages
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, each time a friendly creature fights, it captures 1 Æmber from your opponent.
func TestTakeHostages(t *testing.T) {
	t.Run("a friendly creature captures 1 Æmber when it fights this turn", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(TakeHostages),
				InPlay: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(5))),
				),
			},
			P2: ct.Side{
				Amber: 3,
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
		})

		h.P1.Play(TakeHostages)
		h.P1.Fight(ally, foe)

		h.Expect(ally).AmberOn(1)
		h.P2.ExpectAmber(2)
	})
}
