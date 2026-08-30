package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Cooperative Hunting
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For each friendly creature in play, deal 1 damage to a creature.
func TestCooperativeHunting(t *testing.T) {
	t.Run("deals damage equal to the number of friendly creatures to a chosen creature", func(t *testing.T) {
		var f1 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(CooperativeHunting),
				InPlay: ct.Cards(
					ct.Bind(&f1, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(5))),
					ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(5)),
				),
			},
		})

		h.P1.Play(CooperativeHunting)
		h.P1.ClickCard(f1)

		h.Expect(f1).Damage(2) // 1 × 2 friendly creatures
	})
}
