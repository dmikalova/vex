package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Imp-losion
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Destroy a friendly creature and an enemy creature.
func TestImplosion(t *testing.T) {
	t.Run("destroys a friendly creature and an enemy creature", func(t *testing.T) {
		var friend, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				Hand:   ct.Cards(Implosion),
				InPlay: ct.Cards(ct.Bind(&friend, ct.Creature(ct.OfHouse(card.House.Dis)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Mars)))),
			},
		})

		h.P1.Play(Implosion)

		h.Expect(friend).At(ct.Discard)
		h.Expect(enemy).At(ct.Discard)
	})
}
