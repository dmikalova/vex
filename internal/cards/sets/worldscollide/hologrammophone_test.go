package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hologrammophone
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Common
//	Æmber:  1
//	Traits: Item
//
//	Action: Ward a Creature.
func TestHologrammophone(t *testing.T) {
	t.Run("wards a chosen creature", func(t *testing.T) {
		var friend ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(Hologrammophone, ct.Bind(&friend, ct.Creature())),
			},
		})

		h.P1.UseAction(Hologrammophone)

		if !h.Game().Warded(friend.ID()) {
			t.Errorf("%s should be warded", friend.Name())
		}
	})
}
