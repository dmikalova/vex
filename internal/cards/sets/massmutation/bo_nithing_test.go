package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bo Nithing
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Play: For each forged key your opponent has, steal 1 Æmber.
func TestBoNithing(t *testing.T) {
	t.Run("steals 1 Æmber for each key the opponent has forged", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(BoNithing)},
			P2: ct.Side{Keys: 2, Amber: 3},
		})

		h.P1.Play(BoNithing)

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(1)
	})

	t.Run("steals nothing when the opponent has forged no keys", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(BoNithing)},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(BoNithing)

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(3)
	})
}
