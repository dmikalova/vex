package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Signal Fire
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Signal Fire. For the remainder of the turn, each friendly Brobnar creature may fight.
func TestSignalFire(t *testing.T) {
	t.Run("sacrifices itself and lets Brobnar creatures fight out of house", func(t *testing.T) {
		var brobnar, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					SignalFire,
					ct.Bind(&brobnar, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6)))),
			},
		})

		// Versatile lets the Brobnar artifact fire its Action out of the active house.
		h.P1.UseAction(SignalFire)
		h.Expect(SignalFire).At(ct.Discard)

		// The off-house Brobnar creature may now fight though Untamed is active.
		h.P1.Fight(brobnar, enemy)
		h.Expect(enemy).Damage(4)
	})
}
