package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Charge!
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: For the remainder of the turn, each time you play a creature, deal 2 damage to an enemy creature.
func TestCharge(t *testing.T) {
	t.Run("each creature played after it deals 2 damage to a chosen enemy", func(t *testing.T) {
		var minion, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand: ct.Cards(
					Charge,
					ct.Bind(&minion, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(5)))),
			},
		})

		h.P1.Play(Charge)
		h.P1.Play(minion) // enters and, per Charge!, deals 2 to the sole enemy

		h.Expect(foe).Damage(2)
	})
}
