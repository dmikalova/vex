package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Protectrix
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Knight • Spirit
//
//	Reap: You may fully heal a creature -> for the remainder of the turn, it cannot be dealt damage.
func TestProtectrix(t *testing.T) {
	t.Run("fully heals a creature and protects it from damage", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					Protectrix,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(6))),
				),
			},
		})

		ally.Damaged(2)

		h.P1.Reap(Protectrix)
		h.P1.ClickOption("Yes")

		h.Expect(ally).Damage(0)
		if !h.Game().State.Cards[ally.ID()].DamageImmune {
			t.Error("the healed creature should be protected from damage")
		}
	})
}
