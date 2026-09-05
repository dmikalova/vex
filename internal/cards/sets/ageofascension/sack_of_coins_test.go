package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sack of Coins
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Choose a creature - for each Æmber in your pool, deal 1 damage to the chosen creature.
func TestSackOfCoins(t *testing.T) {
	t.Run("deals 1 damage per aember in your pool", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Amber: 2,
				Hand:  ct.Cards(SackOfCoins),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20))))},
		})

		h.P1.Play(SackOfCoins)

		h.Expect(foe).Damage(3) // pool of 2 + 1 Æmber bonus on play
	})
}
