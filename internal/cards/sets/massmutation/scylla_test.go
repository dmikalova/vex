package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Scylla
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Beast
//
//	Each enemy creature gains, "Reap: Deal 4 damage to this creature."
func TestScylla(t *testing.T) {
	t.Run("an enemy creature takes 4 damage when it reaps", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(Scylla),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6))))},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.Reap(enemy)

		h.Expect(enemy).Damage(4)
	})

	t.Run("a friendly creature reaping is unaffected", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					Scylla,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(6))),
				),
			},
		})

		h.P1.Reap(ally)

		h.Expect(ally).Damage(0)
	})
}
