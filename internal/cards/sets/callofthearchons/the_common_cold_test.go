package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Common Cold
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Deal 1 damage to each creature. You may destroy each Mars creature.
func TestTheCommonCold(t *testing.T) {
	t.Run("damages each creature and may destroy each Mars creature", func(t *testing.T) {
		var marsFoe, brobFoe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(TheCommonCold)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&marsFoe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&brobFoe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Play(TheCommonCold)
		h.P1.ClickOption("Yes")

		h.Expect(marsFoe).At(ct.Discard)
		h.Expect(brobFoe).At(ct.PlayArea).Damage(1)
	})
}
