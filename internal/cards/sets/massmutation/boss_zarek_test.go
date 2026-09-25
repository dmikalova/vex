package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Boss Zarek
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant • Thief
//
//	Each friendly creature with Æmber on it gains elusive.
//	Enhance Capture Capture Capture.
func TestBossZarek(t *testing.T) {
	t.Run("a friendly creature with Æmber on it gains elusive", func(t *testing.T) {
		var attacker, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&attacker, ct.Creature(ct.Power(3)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(BossZarek, ct.Bind(&ally, ct.Creature(ct.Power(4)))),
			},
		})
		h.Game().State.Cards[ally.ID()].Amber = 1

		h.P1.Fight(attacker, ally)

		h.Expect(ally).Damage(0) // elusive: no damage the first time it is attacked
	})

	t.Run("a friendly creature without Æmber does not gain elusive", func(t *testing.T) {
		var attacker, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&attacker, ct.Creature(ct.Power(3)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(BossZarek, ct.Bind(&ally, ct.Creature(ct.Power(4)))),
			},
		})

		h.P1.Fight(attacker, ally)

		h.Expect(ally).Damage(3) // no elusive, so it takes the attacker's power
	})
}
