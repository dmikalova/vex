package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Cannon
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Weapon
//
//	Action: Deal 2 damage to a creature.
func TestCannon(t *testing.T) {
	t.Run("deals 2 damage to a chosen creature", func(t *testing.T) {
		var target ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(Cannon)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&target, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(5))),
				),
			},
		})

		h.P1.UseAction(Cannon)

		h.Expect(target).Damage(2)
	})
}
