package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Kelifi Dragon
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  12
//	Traits: Dragon
//
//	Kelifi Dragon cannot be played unless you have 5 Æmber or more.
//	Fight/Reap: Gain 1 Æmber, and deal 5 damage to a creature.
func TestKelifiDragon(t *testing.T) {
	t.Run("cannot be played on a pool under 5", func(t *testing.T) {
		var dragon ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(ct.Bind(&dragon, KelifiDragon)),
				Amber: 4,
			},
		})

		h.P1.ExpectCannotPlay(dragon)
	})

	t.Run("reaping gains 1 Æmber and deals 5 damage", func(t *testing.T) {
		var dragon, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&dragon, KelifiDragon)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(9)))),
			},
		})

		h.P1.Reap(dragon)
		h.P1.ClickCard(enemy)

		h.P1.ExpectAmber(2) // 1 for the reap, 1 from the ability
		h.Expect(enemy).Damage(5)
	})
}
