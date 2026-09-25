package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Gron Nine-Toes
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Gron Nine-Toes gains +4 power while it is damaged.
func TestGronNineToes(t *testing.T) {
	t.Run("undamaged, its power stays at base", func(t *testing.T) {
		var gron ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&gron, GronNineToes)),
			},
		})

		h.Expect(gron).Power(5)
	})

	t.Run("while damaged, it gains +4 power", func(t *testing.T) {
		var gron ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&gron, GronNineToes)),
			},
		})

		gron.Damaged(1)
		h.Expect(gron).Power(9)
	})
}
