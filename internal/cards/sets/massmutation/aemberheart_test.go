package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Aemberheart
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: Choose a friendly creature - exalt and ward the chosen creature, and fully heal the chosen creature.
func TestAemberheart(t *testing.T) {
	t.Run("exalts, wards, and fully heals a friendly creature", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					Aemberheart,
					ct.Bind(&ally, ct.Creature(ct.Power(5))),
				),
			},
		})
		h.Game().State.Cards[ally.ID()].Damage = 3

		h.P1.UseAction(Aemberheart)

		h.Expect(ally).AmberOn(1).Damage(0)
		if !h.Game().Warded(ally.ID()) {
			t.Errorf("%s should be warded", ally.Name())
		}
	})
}
