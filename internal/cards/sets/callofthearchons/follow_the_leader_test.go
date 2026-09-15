package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Follow the Leader
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For the remainder of the turn, each friendly creature may fight.
func TestFollowTheLeader(t *testing.T) {
	t.Run("lets an out-of-house friendly creature fight", func(t *testing.T) {
		var follow, outsider, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(ct.Bind(&follow, FollowTheLeader)),
				InPlay: ct.Cards(
					ct.Bind(&outsider, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(5))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(2)))),
			},
		})

		h.P1.ExpectCannotUse(outsider)
		h.P1.Play(follow)
		h.P1.Fight(outsider, enemy)

		h.Expect(enemy).At(ct.Discard)
	})
}
