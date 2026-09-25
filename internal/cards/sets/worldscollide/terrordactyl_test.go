package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Terrordactyl
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  12
//	Traits: Beast
//
//	Terrordactyl deals 4 damage when fighting.
//	Terrordactyl enters play stunned.
//	Before Fight: Deal 4 damage to each neighbor of the creature Terrordactyl fights.
func TestTerrordactyl(t *testing.T) {
	t.Run("enters play stunned", func(t *testing.T) {
		var terror ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(ct.Bind(&terror, Terrordactyl)),
			},
		})

		h.P1.Play(Terrordactyl)

		h.Expect(terror).Stunned(true)
	})

	t.Run(
		"before fight, deals 4 damage to each neighbor of the fought creature",
		func(t *testing.T) {
			var terror, left, target, right ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Saurian,
					InPlay: ct.Cards(ct.Bind(&terror, Terrordactyl)),
				},
				P2: ct.Side{InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(6))),
					ct.Bind(&target, ct.Creature(ct.Power(6))),
					ct.Bind(&right, ct.Creature(ct.Power(6))),
				)},
			})

			h.P1.Fight(terror, target)

			h.Expect(left).At(ct.PlayArea).Damage(4)
			h.Expect(right).At(ct.PlayArea).Damage(4)
		},
	)
}
