package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hallowed Shield
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Common
//	Traits: Item
//
//	Action: For the remainder of the turn, a creature cannot be dealt damage.
func TestHallowedShield(t *testing.T) {
	t.Run("prevents damage to a creature this turn", func(t *testing.T) {
		var attacker, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(20))),
					HallowedShield,
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20))))},
		})

		h.P1.UseAction(HallowedShield)
		h.P1.ClickCard(foe)
		h.P1.Fight(attacker, foe)

		h.Expect(foe).Damage(0)
	})
}
