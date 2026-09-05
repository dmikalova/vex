package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// 1-2 Punch
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose an enemy creature - if that creature was already stunned, destroy it. Otherwise, stun it.
func TestCard12Punch(t *testing.T) {
	t.Run("stuns an unstunned enemy creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(Card12Punch)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(4))))},
		})

		h.P1.Play(Card12Punch)

		h.Expect(foe).At(ct.PlayArea)
		h.Expect(foe).Stunned(true)
	})

	t.Run("destroys an enemy creature that was already stunned", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(Card12Punch)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(4))))},
		})
		foe.Stun()

		h.P1.Play(Card12Punch)

		h.Expect(foe).At(ct.Discard)
	})
}
