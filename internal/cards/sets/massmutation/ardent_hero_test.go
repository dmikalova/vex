package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Ardent Hero
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human • Knight
//
//	Taunt.
//	Ardent Hero cannot be dealt damage by Mutant creatures or creatures with power 5 or higher.
func TestArdentHero(t *testing.T) {
	t.Run("a Mutant creature deals it no damage", func(t *testing.T) {
		var hero, mutant ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&hero, ArdentHero)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(
					&mutant,
					ct.Creature(ct.Power(3), ct.Traits(card.Traits.Mutant), ct.Armor(4)),
				),
			)},
		})

		h.P1.Fight(hero, mutant)

		h.Expect(hero).Damage(0)
	})

	t.Run("a power 5 or higher creature deals it no damage", func(t *testing.T) {
		var hero, brute ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&hero, ArdentHero)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&brute, ct.Creature(ct.Power(6))))},
		})

		h.P1.Fight(hero, brute)

		h.Expect(hero).Damage(0)
	})

	t.Run("a weaker non-Mutant creature deals it damage", func(t *testing.T) {
		var hero, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&hero, ArdentHero)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3), ct.Armor(4))))},
		})

		h.P1.Fight(hero, foe)

		h.Expect(hero).Damage(3)
	})
}
