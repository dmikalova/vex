package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Beware the Ides
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 23 damage to a creature in the center of its controller's battleline.
func TestBewareTheIdes(t *testing.T) {
	t.Run("deals 23 damage to the center creature", func(t *testing.T) {
		var left, center, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(BewareTheIdes),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				ct.Bind(&center, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
			)},
		})

		h.P1.Play(BewareTheIdes)

		h.Expect(center).At(ct.Discard)
		h.Expect(left).At(ct.PlayArea)
		h.Expect(right).At(ct.PlayArea)
	})
}
