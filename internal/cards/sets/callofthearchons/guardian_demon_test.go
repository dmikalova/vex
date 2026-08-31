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
		h.P1.ClickCard(wounded) // heal up to 2 from the wounded creature
		h.P1.ClickCard(foe)     // deal the 2 healed to the enemy

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
		h.P1.ClickCard(wounded)
		h.P1.ClickCard(foe)

		h.Expect(wounded).Damage(0)
		h.Expect(foe).Damage(1)
	})
}
