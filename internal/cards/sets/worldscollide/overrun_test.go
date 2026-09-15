package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Overrun
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: If 3 or more enemy creatures have been destroyed this turn, your opponent loses 2 Æmber.
func TestOverrun(t *testing.T) {
	t.Run("opponent loses 2 Æmber at or above the threshold", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(Overrun)},
			P2: ct.Side{Amber: 3},
		})
		h.Game().State.TurnHistory[0][card.TurnStat.EnemyCreaturesDestroyed] = 3

		h.P1.Play(Overrun)

		h.P2.ExpectAmber(1)
	})

	t.Run("nothing happens below the threshold", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(Overrun)},
			P2: ct.Side{Amber: 3},
		})
		h.Game().State.TurnHistory[0][card.TurnStat.EnemyCreaturesDestroyed] = 2

		h.P1.Play(Overrun)

		h.P2.ExpectAmber(3)
	})
}
