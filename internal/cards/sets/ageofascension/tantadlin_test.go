package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tantadlin
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  9
//	Traits: Tree
//
//	Tantadlin deals 2 Damage when fighting.
//	Fight: Your opponent discards a random card from their archives.
func TestTantadlin(t *testing.T) {
	t.Run("deals only 2 fight damage", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(Tantadlin),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20))))},
		})

		h.P1.Fight(Tantadlin, foe)

		h.Expect(foe).Damage(2)
	})

	t.Run("fighting discards a random card from the opponent's archives", func(t *testing.T) {
		var foe, archived ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(Tantadlin),
			},
			P2: ct.Side{
				InPlay:   ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(1)))),
				Archives: ct.Cards(ct.Bind(&archived, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Fight(Tantadlin, foe)

		h.Expect(archived).At(ct.Discard)
	})
}
