package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Blood Money
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Exalt an enemy creature 2 times.
func TestBloodMoney(t *testing.T) {
	t.Run("exalts a chosen enemy creature twice", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(BloodMoney)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))))},
		})

		h.P1.Play(BloodMoney)

		h.Expect(foe).AmberOn(2)
	})
}
