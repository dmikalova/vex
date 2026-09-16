package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sloth
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Traits: Demon • Sin
//
//	At the end of your turn, if you did not use any creatures this turn, for each friendly Sin creature, gain 1 Æmber.
func TestSloth(t *testing.T) {
	t.Run("gains 1 Æmber per friendly Sin creature when no creature was used", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					Sloth,
					ct.Creature(
						ct.OfHouse(card.House.Dis),
						ct.Traits(card.Traits.Sin),
						ct.Power(3),
					),
				),
			},
		})

		h.P1.EndTurn()

		h.P1.ExpectAmber(2) // Sloth and one other Sin creature
	})

	t.Run("gains nothing when a creature was used", func(t *testing.T) {
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					Sloth,
					ct.Bind(&reaper, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(reaper)
		h.P1.EndTurn()

		h.P1.ExpectAmber(1) // only the reap's Æmber, no Sloth gain
	})
}
