package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Creed of Nature
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Power
//
//	Versatile.
//	Action: Destroy Creed of Nature. Choose a Creature - for the remainder of the turn, it gains skirmish and assault equal to its power.
func TestCreedOfNature(t *testing.T) {
	t.Run(
		"sacrifices itself, then a chosen creature gains skirmish and assault equal to its power",
		func(t *testing.T) {
			var chosen, defender ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Untamed,
					InPlay: ct.Cards(
						CreedOfNature,
						ct.Bind(&chosen, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
					),
				},
				P2: ct.Side{InPlay: ct.Cards(
					ct.Bind(&defender, ct.Creature(ct.Power(5))),
				)},
			})

			h.P1.UseAction(CreedOfNature)
			h.P1.ClickCard(chosen)

			// The Omni sacrifices Creed of Nature.
			h.Expect(CreedOfNature).At(ct.Discard)

			// The chosen creature (power 3) now has assault 3: 3 assault + 3 fight damage
			// destroys the 5-power defender, and skirmish spares the attacker its return
			// damage, so it survives undamaged.
			h.P1.Fight(chosen, defender)
			h.Expect(defender).At(ct.Discard)
			h.Expect(chosen).At(ct.PlayArea).Damage(0)
		},
	)

	t.Run("assault scales with the chosen creature's own power", func(t *testing.T) {
		var chosen, defender ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					CreedOfNature,
					ct.Bind(&chosen, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(6))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&defender, ct.Creature(ct.Power(9))),
			)},
		})

		h.P1.UseAction(CreedOfNature)
		h.P1.ClickCard(chosen)

		// Assault 6 (its power) + 6 fight damage = 12 destroys the 9-power defender.
		h.P1.Fight(chosen, defender)
		h.Expect(defender).At(ct.Discard)
		h.Expect(chosen).At(ct.PlayArea).Damage(0)
	})
}
