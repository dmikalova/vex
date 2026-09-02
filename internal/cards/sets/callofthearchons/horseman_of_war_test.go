package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Horseman of War
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Connected
//	Power:  5
//	Traits: Horseman • Spirit
//
//	Play: For the remainder of the turn, each friendly creature may fight.
func TestHorsemanOfWar(t *testing.T) {
	t.Run("lets an out-of-house friendly creature fight", func(t *testing.T) {
		var horseman, outsider, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(ct.Bind(&horseman, HorsemanOfWar)),
				InPlay: ct.Cards(
					ct.Bind(&outsider, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(5))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(2)))),
			},
		})

		h.P1.ExpectCannotUse(outsider)
		h.P1.Play(horseman)
		h.P1.Fight(outsider, enemy)

		h.Expect(enemy).At(ct.Discard)
	})
}
