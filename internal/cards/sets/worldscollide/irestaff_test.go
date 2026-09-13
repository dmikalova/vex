package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Irestaff
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Common
//	Æmber:  1
//	Traits: Weapon
//
//	Action: Choose a Creature - enrage it, and give it a +1 power counter.
func TestIrestaff(t *testing.T) {
	t.Run("enrages a creature and gives it a +1 power counter", func(t *testing.T) {
		var troll ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(Irestaff, ct.Bind(&troll, ct.Creature(ct.Power(4)))),
			},
		})

		h.P1.UseAction(Irestaff)

		h.Expect(troll).Power(5)
		if !h.Game().Enraged(troll.ID()) {
			t.Errorf("%s should be enraged", troll.Name())
		}
	})
}
