package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Kartanoo
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Beast
//
//	Reap: Use an artifact.
func TestKartanoo(t *testing.T) {
	t.Run("reaps to use an artifact controlled by the opponent", func(t *testing.T) {
		var theirs ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(Kartanoo),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&theirs, MushroomWithAView))},
		})

		h.P1.Reap(Kartanoo)

		h.Expect(theirs).Exhausted() // used as if it were yours
		h.P1.ExpectAmber(1)          // the reap itself
	})

	t.Run("does nothing extra when no artifact is in play", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(Kartanoo),
			},
		})

		h.P1.Reap(Kartanoo)

		h.P1.ExpectAmber(1)
	})
}
