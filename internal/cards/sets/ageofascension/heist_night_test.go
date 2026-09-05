package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Heist Night
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//
//	Alpha.
//	Play: For each friendly Thief trait creature, steal 1 Æmber.
func TestHeistNight(t *testing.T) {
	t.Run("steals 1 aember for each friendly thief creature", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(HeistNight),
				InPlay: ct.Cards(
					ct.Creature(ct.Traits(card.Traits.Thief)),
					ct.Creature(ct.Traits(card.Traits.Thief)),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(HeistNight)

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(1)
	})
}
