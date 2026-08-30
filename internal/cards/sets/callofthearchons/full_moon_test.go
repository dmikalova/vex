package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Full Moon
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, each time you play a creature, gain 1 Æmber.
func TestFullMoon(t *testing.T) {
	t.Run("gains 1 Æmber for each creature played after it this turn", func(t *testing.T) {
		var c1, c2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand: ct.Cards(
					FullMoon,
					ct.Bind(&c1, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
					ct.Bind(&c2, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
		})

		h.P1.Play(FullMoon)
		h.P1.ExpectAmber(0) // playing Full Moon itself gains nothing

		h.P1.Play(c1)
		h.P1.ExpectAmber(1)
		h.P1.Play(c2)
		h.P1.ExpectAmber(2)
	})
}
