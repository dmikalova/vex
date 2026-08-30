package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Guilty Hearts
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy each creature with Æmber on it.
func TestGuiltyHearts(t *testing.T) {
	t.Run("destroys each creature with Æmber on it", func(t *testing.T) {
		var rich, poor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(GuiltyHearts)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&rich, ct.Creature(ct.OfHouse(card.House.Mars))),
					ct.Bind(&poor, ct.Creature(ct.OfHouse(card.House.Mars))),
				),
			},
		})

		h.Game().State.Cards[rich.ID()].Amber = 2
		h.P1.Play(GuiltyHearts)

		h.Expect(rich).At(ct.Discard)
		h.Expect(poor).At(ct.PlayArea)
	})
}
