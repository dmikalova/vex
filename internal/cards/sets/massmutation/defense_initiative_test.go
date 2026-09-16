package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Defense Initiative
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Choose a creature - ward the chosen creature, and you may exalt the chosen creature -> ward the chosen creature.
func TestDefenseInitiative(t *testing.T) {
	t.Run("wards a creature, then exalts it and wards its neighbors", func(t *testing.T) {
		var left, chosen, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(DefenseInitiative),
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(3))),
					ct.Bind(&chosen, ct.Creature(ct.Power(3))),
					ct.Bind(&right, ct.Creature(ct.Power(3))),
				),
			},
		})

		h.P1.Play(DefenseInitiative)
		h.P1.ClickCard(chosen)  // choose the creature to ward
		h.P1.ClickOption("Yes") // accept the may: exalt it and ward its neighbors

		h.Expect(chosen).AmberOn(1)
		if !h.Game().Warded(chosen.ID()) {
			t.Error("chosen creature should be warded")
		}
		if !h.Game().Warded(left.ID()) || !h.Game().Warded(right.ID()) {
			t.Error("both neighbors should be warded")
		}
	})

	t.Run("declining the exalt wards only the chosen creature", func(t *testing.T) {
		var left, chosen, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(DefenseInitiative),
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(3))),
					ct.Bind(&chosen, ct.Creature(ct.Power(3))),
					ct.Bind(&right, ct.Creature(ct.Power(3))),
				),
			},
		})

		h.P1.Play(DefenseInitiative)
		h.P1.ClickCard(chosen) // choose the creature to ward
		h.P1.ClickOption("No") // decline the may

		h.Expect(chosen).AmberOn(0)
		if !h.Game().Warded(chosen.ID()) {
			t.Error("chosen creature should still be warded")
		}
		if h.Game().Warded(left.ID()) || h.Game().Warded(right.ID()) {
			t.Error("neighbors should not be warded when the exalt is declined")
		}
	})
}
