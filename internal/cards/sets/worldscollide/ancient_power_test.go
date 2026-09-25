package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Ancient Power
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Ward each friendly creature with Æmber on it.
func TestAncientPower(t *testing.T) {
	t.Run("wards each friendly creature that has Æmber on it", func(t *testing.T) {
		var withAember, bare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(AncientPower),
				InPlay: ct.Cards(
					ct.Bind(&withAember, ct.Creature(ct.Power(3))),
					ct.Bind(&bare, ct.Creature(ct.Power(3))),
				),
			},
			P2: ct.Side{},
		})
		h.Game().State.Cards[withAember.ID()].Amber = 1

		h.P1.Play(AncientPower)

		if !h.Game().Warded(withAember.ID()) {
			t.Errorf("%s should be warded", withAember.Name())
		}
		if h.Game().Warded(bare.ID()) {
			t.Errorf("%s should not be warded", bare.Name())
		}
	})
}
