package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Guardian Demon
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Demon
//
//	Play/Fight/Reap: Heal 2 damage from a creature. Deal that amount of damage to another creature.
func TestGuardianDemon(t *testing.T) {
	t.Run("heals a creature and deals the healed amount to another", func(t *testing.T) {
		var wounded, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					GuardianDemon,
					ct.Bind(&wounded, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(5))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(5))))},
		})
		wounded.Damaged(2)

		h.P1.Reap(GuardianDemon)
		// The wounded creature is the only one that can be healed, so the heal
		// resolves without a prompt; only the damage target is asked for.
		h.P1.ClickCard(foe) // deal the 2 healed to the enemy

		h.Expect(wounded).Damage(0)
		h.Expect(foe).Damage(2)
	})

	t.Run("only heals and deals as much damage as the creature has", func(t *testing.T) {
		var wounded, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					GuardianDemon,
					ct.Bind(&wounded, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(5))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(5))))},
		})
		wounded.Damaged(1)

		h.P1.Reap(GuardianDemon)
		h.P1.ClickCard(foe)

		h.Expect(wounded).Damage(0)
		h.Expect(foe).Damage(1)
	})

	t.Run("asks for no target when there is nothing to heal", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					GuardianDemon,
					ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(5)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(5))))},
		})

		// No creature has damage, so neither the heal nor the follow-up damage has
		// anything to do and the reap resolves without asking anything.
		h.P1.Reap(GuardianDemon)

		h.Expect(foe).Damage(0)
		h.Expect(GuardianDemon).Exhausted()
	})
}
