package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lethal Distraction
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: For the remainder of the turn, whenever a creature takes damage, it takes an additional 2 damage.
func TestLethalDistraction(t *testing.T) {
	t.Run("the chosen creature takes 2 extra damage from each instance", func(t *testing.T) {
		var attacker, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(LethalDistraction),
				InPlay: ct.Cards(
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(10)))),
			},
		})

		h.P1.Play(LethalDistraction)
		h.P1.ClickCard(foe)
		h.P1.Fight(attacker, foe)

		// The 3 combat damage lands as 3 + 2 = 5 on the chosen enemy.
		h.Expect(foe).Damage(5)
	})
}
