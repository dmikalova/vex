package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Nurse Soto
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Human
//
//	Deploy.
//	Play/Fight/Reap: Heal 3 damage from each neighboring creature.
func TestNurseSoto(t *testing.T) {
	t.Run("reaping heals 3 damage from each neighbor", func(t *testing.T) {
		var left, soto, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(6))),
					ct.Bind(&soto, NurseSoto),
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(6))),
				),
			},
		})
		left.Damaged(5)
		right.Damaged(5)

		h.P1.Reap(soto)

		h.Expect(left).Damage(2)
		h.Expect(right).Damage(2)
	})
}
