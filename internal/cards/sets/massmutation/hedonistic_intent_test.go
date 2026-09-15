package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hedonistic Intent
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Exalt each flank creature.
func TestHedonisticIntent(t *testing.T) {
	t.Run("exalts flank creatures but not the middle", func(t *testing.T) {
		var left, middle, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(HedonisticIntent),
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature()),
					ct.Bind(&middle, ct.Creature()),
					ct.Bind(&right, ct.Creature()),
				),
			},
		})

		h.P1.Play(HedonisticIntent)

		h.Expect(left).AmberOn(1)
		h.Expect(middle).AmberOn(0)
		h.Expect(right).AmberOn(1)
	})
}
