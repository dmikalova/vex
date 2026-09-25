package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Mega Gron Nine-Toes
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  7
//	Traits: Giant
//
//	Mega Gron Nine-Toes gains +4 power while it is damaged.
func TestMegaGronNineToes(t *testing.T) {
	t.Run("undamaged, its power stays at base", func(t *testing.T) {
		var gron ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&gron, MegaGronNineToes)),
			},
		})

		h.Expect(gron).Power(7)
	})

	t.Run("while damaged, it gains +4 power", func(t *testing.T) {
		var gron ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&gron, MegaGronNineToes)),
			},
		})

		gron.Damaged(1)
		h.Expect(gron).Power(11)
	})
}
