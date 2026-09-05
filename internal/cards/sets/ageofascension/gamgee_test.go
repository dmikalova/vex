package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gamgee
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Elusive.
//	Reap: If your opponent has more Æmber than you, steal 1 Æmber.
func TestGamgee(t *testing.T) {
	t.Run("steals 1 aember when reaping while the opponent has more aember", func(t *testing.T) {
		var gamgee ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&gamgee, Gamgee)),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Reap(gamgee)

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(4)
	})
}
