package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Pale Star
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Power
//
//	Versatile.
//	Action: Destroy The Pale Star. For the remainder of the turn, each creature is considered to have 1 power and 0 armor.
func TestThePaleStar(t *testing.T) {
	t.Run(
		"masks each creature to 1 power and 0 armor, destroying newly lethal ones",
		func(t *testing.T) {
			var mine, doomed ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Dis,
					InPlay: ct.Cards(
						ThePaleStar,
						ct.Bind(&mine, ct.Creature(ct.Power(5), ct.Armor(2))),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(ct.Bind(&doomed, ct.Creature(ct.Power(6)))),
				},
			})
			// An enemy creature with 3 damage survives at 6 power, but not once masked to 1.
			h.Game().SetDamage(doomed.ID(), 3)

			h.P1.UseAction(ThePaleStar)

			h.Expect(ThePaleStar).At(ct.Discard)
			if got := h.Game().Power(mine.ID()); got != 1 {
				t.Errorf("masked power = %d, want 1", got)
			}
			if got := h.Game().Armor(mine.ID()); got != 0 {
				t.Errorf("masked armor = %d, want 0", got)
			}
			h.Expect(doomed).At(ct.Discard)
		},
	)
}
