package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Finch Cloak
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Elf • Thief
//
//	Fight/Reap: If your opponent has more Æmber than you, steal 1 Æmber. Otherwise, each player gains 1 Æmber.
func TestFinchCloak(t *testing.T) {
	t.Run("steals 1 Æmber when behind the opponent", func(t *testing.T) {
		var finch ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				Amber:  1,
				InPlay: ct.Cards(ct.Bind(&finch, FinchCloak)),
			},
			P2: ct.Side{Amber: 4},
		})
		finch.Ready()

		h.P1.Reap(finch)

		h.P1.ExpectAmber(3) // 1 pool + 1 reap + 1 stolen
		h.P2.ExpectAmber(3)
	})

	t.Run("each player gains 1 Æmber when not behind", func(t *testing.T) {
		var finch ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				Amber:  4,
				InPlay: ct.Cards(ct.Bind(&finch, FinchCloak)),
			},
			P2: ct.Side{Amber: 1},
		})
		finch.Ready()

		h.P1.Reap(finch)

		h.P1.ExpectAmber(6) // 4 pool + 1 reap + 1 gained
		h.P2.ExpectAmber(2)
	})
}
