package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Xenos Bloodshadow
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Human • Witch
//
//	Elusive, Poison, Skirmish, Hazardous 6.
func TestXenosBloodshadow(t *testing.T) {
	t.Run("elusive: the first time it is attacked it takes no damage", func(t *testing.T) {
		var xenos, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&xenos, XenosBloodshadow)),
			},
			P2: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(10)))),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.Fight(foe, xenos)

		h.Expect(xenos).Damage(0)
	})
}
