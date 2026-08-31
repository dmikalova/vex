package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Relentless Whispers
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 2 damage to a creature. If this damage destroys that creature, steal 1 Æmber.
func TestRelentlessWhispers(t *testing.T) {
	t.Run("steals 1 Æmber when its damage destroys the creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(RelentlessWhispers)},
			P2: ct.Side{
				Amber: 3,
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
		})

		h.P1.Play(RelentlessWhispers)

		h.Expect(foe).At(ct.Discard)
		h.P2.ExpectAmber(2) // 1 stolen
	})
}
