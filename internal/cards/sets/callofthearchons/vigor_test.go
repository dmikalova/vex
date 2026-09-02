package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Vigor
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Heal 3 damage from a creature. If you healed exactly 3 damage, gain 1 Æmber.
func TestVigor(t *testing.T) {
	t.Run("gains 1 Æmber when it heals the full 3 damage", func(t *testing.T) {
		var hurt ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Hand:   ct.Cards(Vigor),
				InPlay: ct.Cards(ct.Bind(&hurt, ct.Creature(ct.Power(6)))),
			},
		})
		hurt.Damaged(4)

		h.P1.Play(Vigor)

		h.Expect(hurt).Damage(1)
		h.P1.ExpectAmber(2) // the Æmber bonus plus the healing bonus
	})

	t.Run("gains nothing when there is less than 3 damage to heal", func(t *testing.T) {
		var hurt ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Hand:   ct.Cards(Vigor),
				InPlay: ct.Cards(ct.Bind(&hurt, ct.Creature(ct.Power(6)))),
			},
		})
		hurt.Damaged(2)

		h.P1.Play(Vigor)

		h.Expect(hurt).Damage(0)
		h.P1.ExpectAmber(1)
	})
}
