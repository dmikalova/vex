package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Cera Severa
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	Fight/Reap: Cera Severa captures 1 Æmber from your opponent.
//	Destroyed: For each Æmber on Cera Severa, deal 1 damage to an enemy creature.
func TestCeraSevera(t *testing.T) {
	t.Run("captures 1 Æmber after it reaps", func(t *testing.T) {
		var cera ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&cera, CeraSevera)),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Reap(cera)

		h.Expect(cera).AmberOn(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("destroyed deals 1 damage per Æmber on it to a chosen enemy", func(t *testing.T) {
		var cera, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&cera, CeraSevera)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(6))))},
		})
		h.Game().AddAmberOn(cera.ID(), 2)

		h.Game().DestroyEach(0, []engine.LocalID{cera.ID()})

		h.Expect(foe).Damage(2)
	})

	t.Run("destroyed with no Æmber deals no damage", func(t *testing.T) {
		var cera, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&cera, CeraSevera)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(6))))},
		})

		h.Game().DestroyEach(0, []engine.LocalID{cera.ID()})

		h.Expect(foe).Damage(0)
	})
}
