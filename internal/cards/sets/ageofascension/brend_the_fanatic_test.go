package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Brend the Fanatic
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	Skirmish.
//	Play: Your opponent gains 1 Æmber.
//	Destroyed: Steal 3 Æmber.
func TestBrendTheFanatic(t *testing.T) {
	t.Run("gives the opponent 1 aember when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(BrendTheFanatic)},
			P2: ct.Side{},
		})

		h.P1.Play(BrendTheFanatic)

		h.P2.ExpectAmber(1)
	})
}
