package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Yantzee Gang
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Elf • Thief
//
//	Action: Steal 1 Æmber.
func TestYantzeeGang(t *testing.T) {
	t.Run("steals 1 aember when used", func(t *testing.T) {
		var yantzee ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&yantzee, YantzeeGang)),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.UseAction(yantzee)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(4)
	})
}
