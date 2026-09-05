package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ronnie Wristclocks
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Play: Steal 1 Æmber, or 2 if your opponent has 7 Æmber or more.
func TestRonnieWristclocks(t *testing.T) {
	t.Run("steals 1 Æmber when the opponent has fewer than 7", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(RonnieWristclocks)},
			P2: ct.Side{House: card.House.Brobnar, Amber: 6},
		})

		h.P1.Play(RonnieWristclocks)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(5)
	})

	t.Run("steals 2 Æmber when the opponent has 7 or more", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(RonnieWristclocks)},
			P2: ct.Side{House: card.House.Brobnar, Amber: 7},
		})

		h.P1.Play(RonnieWristclocks)

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(5)
	})
}
