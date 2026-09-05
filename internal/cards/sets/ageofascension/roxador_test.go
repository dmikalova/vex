package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Roxador
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Beast
//
//	Skirmish.
//	Roxador deals 2 Damage when fighting.
//	Fight: Stun the creature Roxador fights.
func TestRoxador(t *testing.T) {
	t.Run("deals only 2 fight damage, takes none, and stuns the defender", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(Roxador),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(6))))},
		})

		h.P1.Fight(Roxador, foe)

		h.Expect(foe).Damage(2)
		h.Expect(foe).Stunned(true)
		h.Expect(Roxador).Damage(0)
	})
}
