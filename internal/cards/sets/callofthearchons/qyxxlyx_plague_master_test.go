package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Qyxxlyx Plague Master
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Martian • Scientist
//
//	Fight/Reap: Deal 3 damage to each Human trait creature, ignoring armor.
func TestQyxxlyxPlagueMaster(t *testing.T) {
	t.Run("deals 3 to each Human creature, bypassing armor", func(t *testing.T) {
		var human, beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, InPlay: ct.Cards(QyxxlyxPlagueMaster)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(
					&human,
					ct.Creature(
						ct.OfHouse(card.House.Sanctum),
						ct.Power(5),
						ct.Armor(2),
						ct.Traits(card.Traits.Human),
					),
				),
				ct.Bind(
					&beast,
					ct.Creature(
						ct.OfHouse(card.House.Sanctum),
						ct.Power(5),
						ct.Traits(card.Traits.Beast),
					),
				),
			)},
		})

		h.P1.Reap(QyxxlyxPlagueMaster)

		h.Expect(human).Damage(3) // armor does not absorb any of it
		h.Expect(beast).Damage(0)
	})
}
