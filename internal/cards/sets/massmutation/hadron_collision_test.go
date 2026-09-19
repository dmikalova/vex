package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hadron Collision
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Choose a creature. Remove a ward from the chosen creature. Deal 3 damage to the chosen creature, ignoring armor.
func TestHadronCollision(t *testing.T) {
	t.Run("removes a ward and deals 3 damage that armor cannot prevent", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(HadronCollision),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.Power(6), ct.Armor(3))),
			)},
		})
		h.Game().SetWarded(foe.ID(), true)

		h.P1.Play(HadronCollision)

		if h.Game().Warded(foe.ID()) {
			t.Error("foe is still warded, want ward removed")
		}
		h.Expect(foe).Damage(3) // 3 armor absorbs none of it
	})
}
