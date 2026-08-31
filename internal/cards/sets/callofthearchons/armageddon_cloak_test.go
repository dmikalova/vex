package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Armageddon Cloak
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains +2 hazardous and, "If this creature would be destroyed, instead fully heal it, and destroy Armageddon Cloak."
func TestArmageddonCloak(t *testing.T) {
	t.Run("fully heals its host and destroys itself instead of the host once", func(t *testing.T) {
		var host, cloak ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
						ct.Bind(&cloak, ArmageddonCloak),
					),
				),
			},
		})
		host.Damaged(2)

		h.Game().DealDamage(0, []engine.DamageTarget{{ID: host.ID(), Amount: 1}})

		h.Expect(host).At(ct.PlayArea).Damage(0)
		h.Expect(cloak).At(ct.Discard)

		h.Game().DealDamage(0, []engine.DamageTarget{{ID: host.ID(), Amount: 3}})

		h.Expect(host).At(ct.Discard)
	})

	t.Run("grants the host +2 hazardous", func(t *testing.T) {
		var attacker, host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(2))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
						ArmageddonCloak,
					),
				),
			},
		})

		h.P1.Fight(attacker, host)

		h.Expect(attacker).At(ct.Discard)
		h.Expect(host).At(ct.PlayArea).Damage(0)
	})
}
