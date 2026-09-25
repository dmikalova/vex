package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Baldric the Bold
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Before Fight: If the fights creature is the most powerful enemy creature, gain 2 Æmber.
func TestBaldricTheBold(t *testing.T) {
	t.Run("gains 2 Æmber fighting the most powerful enemy creature", func(t *testing.T) {
		var baldric, big, small ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&baldric, BaldricTheBold)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&big, ct.Creature(ct.Power(6))),
				ct.Bind(&small, ct.Creature(ct.Power(3))),
			)},
		})

		h.P1.Fight(baldric, big)

		h.P1.ExpectAmber(2)
	})

	t.Run("a tie for most powerful still qualifies", func(t *testing.T) {
		var baldric, a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&baldric, BaldricTheBold)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&a, ct.Creature(ct.Power(5))),
				ct.Bind(&b, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Fight(baldric, a)

		h.P1.ExpectAmber(2)
	})

	t.Run("no Æmber fighting a lesser enemy creature", func(t *testing.T) {
		var baldric, big, small ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&baldric, BaldricTheBold)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&big, ct.Creature(ct.Power(6))),
				ct.Bind(&small, ct.Creature(ct.Power(3))),
			)},
		})

		h.P1.Fight(baldric, small)

		h.P1.ExpectAmber(0)
	})
}
