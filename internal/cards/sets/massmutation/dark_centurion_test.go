package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dark Centurion
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Mutant • Soldier
//
//	Action: Move 1 Æmber from a creature to the common supply -> ward the chosen creature.
//	Enhance Capture Capture.
func TestDarkCenturion(t *testing.T) {
	t.Run("moves 1 Æmber to the supply and wards that creature", func(t *testing.T) {
		var centurion, rich ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&centurion, DarkCenturion)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&rich, ct.Creature(ct.Power(4))),
			)},
		})

		h.Game().AddAmberOn(rich.ID(), 2)
		h.P1.UseAction(centurion)
		h.P1.ClickCard(rich)

		h.Expect(rich).AmberOn(1)
		if !h.Game().Warded(rich.ID()) {
			t.Error("the creature Æmber was moved off should be warded")
		}
	})

	t.Run("does not ward a creature holding no Æmber", func(t *testing.T) {
		var centurion, bare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&centurion, DarkCenturion)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&bare, ct.Creature(ct.Power(4))),
			)},
		})

		h.P1.UseAction(centurion)
		h.P1.ClickCard(bare)

		h.Expect(bare).AmberOn(0)
		if h.Game().Warded(bare.ID()) {
			t.Error("a creature holding no Æmber should not be warded")
		}
	})
}
