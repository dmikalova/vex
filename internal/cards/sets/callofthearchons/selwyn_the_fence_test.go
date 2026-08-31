package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Selwyn the Fence
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Elf • Thief
//
//	Fight/Reap: Move 1 Æmber from one of your cards to your pool.
func TestSelwynTheFence(t *testing.T) {
	t.Run("moves 1 Æmber from a friendly card to the pool on reap", func(t *testing.T) {
		var vault ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					SelwynTheFence,
					ct.Bind(&vault, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3))),
				),
			},
		})
		h.Game().State.Cards[vault.ID()].Amber = 2

		h.P1.Reap(SelwynTheFence)

		h.Expect(vault).AmberOn(1)
		h.P1.ExpectAmber(2) // 1 from reaping + 1 moved from the card
	})
}
