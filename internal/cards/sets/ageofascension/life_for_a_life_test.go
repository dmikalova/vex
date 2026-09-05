package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Life for a Life
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Destroy a friendly creature -> deal 6 damage to a creature.
func TestLifeForALife(t *testing.T) {
	t.Run("destroys a friendly creature to deal 6 damage to a creature", func(t *testing.T) {
		var sacrifice, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(LifeForALife),
				InPlay: ct.Cards(
					ct.Bind(&sacrifice, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(4))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20))))},
		})

		h.P1.Play(LifeForALife)

		h.Expect(sacrifice).At(ct.Discard)
		h.Expect(foe).Damage(6)
	})
}
