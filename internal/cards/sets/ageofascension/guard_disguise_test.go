package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Guard Disguise
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Destroy Guard Disguise, and if your opponent has 3 Æmber or fewer, steal 3 Æmber.
func TestGuardDisguise(t *testing.T) {
	t.Run("destroys itself and steals 3 when opponent has 3 or fewer", func(t *testing.T) {
		var disguise ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&disguise, GuardDisguise)),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(GuardDisguise)

		h.Expect(disguise).At(ct.Discard)
		h.P1.ExpectAmber(3)
		h.P2.ExpectAmber(0)
	})

	t.Run("steals nothing when opponent has 4 or more", func(t *testing.T) {
		var disguise ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&disguise, GuardDisguise)),
			},
			P2: ct.Side{Amber: 4},
		})

		h.P1.UseAction(GuardDisguise)

		h.Expect(disguise).At(ct.Discard)
		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(4)
	})
}
