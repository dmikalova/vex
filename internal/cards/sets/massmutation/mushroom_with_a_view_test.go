package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mushroom with a View
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Location
//
//	Versatile.
//	Action: Heal 1 damage from each friendly creature.
func TestMushroomWithAView(t *testing.T) {
	t.Run("heals 1 damage from each friendly creature", func(t *testing.T) {
		var friend ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					MushroomWithAView,
					ct.Bind(&friend, ct.Creature(ct.Power(4))),
				),
			},
		})

		friend.Damaged(2)
		h.P1.UseAction(MushroomWithAView)

		h.Expect(friend).Damage(1)
	})
}
