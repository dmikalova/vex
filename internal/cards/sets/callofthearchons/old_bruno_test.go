package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Old Bruno
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	Elusive.
//	Play: Old Bruno captures 3 Æmber from your opponent.
func TestOldBruno(t *testing.T) {
	t.Run("captures 3 Æmber when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(OldBruno)},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(OldBruno)

		h.Expect(OldBruno).AmberOn(3)
		h.P2.ExpectAmber(2)
	})
}
