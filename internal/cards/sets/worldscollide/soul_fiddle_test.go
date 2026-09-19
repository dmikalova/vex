package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Soul Fiddle
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: Enrage a creature.
func TestSoulFiddle(t *testing.T) {
	t.Run("enrages a chosen creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(SoulFiddle),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature()))},
		})

		h.P1.UseAction(SoulFiddle)

		if !h.Game().Enraged(foe.ID()) {
			t.Errorf("%s should be enraged", foe.Name())
		}
	})
}
