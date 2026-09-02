package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Routine Job
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Steal 1 Æmber, and for each copy of Routine Job in your discard pile, steal 1 Æmber.
func TestRoutineJob(t *testing.T) {
	t.Run("steals 1 with no copies in the discard pile", func(t *testing.T) {
		var job ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(ct.Bind(&job, RoutineJob)),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(job)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(4)
	})

	t.Run("steals 1 more for each copy already discarded", func(t *testing.T) {
		var job ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Shadows,
				Hand:    ct.Cards(ct.Bind(&job, RoutineJob)),
				Discard: ct.Cards(RoutineJob, RoutineJob),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(job)

		h.P1.ExpectAmber(3)
		h.P2.ExpectAmber(2)
	})
}
