package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shield of Justice
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: For the remainder of the turn, each friendly creature cannot be dealt damage.
func TestShieldOfJustice(t *testing.T) {
	t.Run("friendly creatures take no damage for the turn", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				Hand:   ct.Cards(ShieldOfJustice),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4)))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))))},
		})

		h.P1.Play(ShieldOfJustice)
		h.P1.Fight(ally, foe)

		h.Expect(ally).Damage(0).At(ct.PlayArea)
	})
}
