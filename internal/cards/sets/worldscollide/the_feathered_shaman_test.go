package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Feathered Shaman
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Human • Witch
//
//	Elusive.
//	Fight/Reap: Ward each neighboring Creature.
func TestTheFeatheredShaman(t *testing.T) {
	t.Run("reaping wards each neighbor", func(t *testing.T) {
		var left, shaman, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4))),
					ct.Bind(&shaman, TheFeatheredShaman),
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4))),
				),
			},
		})

		h.P1.Reap(shaman)

		if !h.Game().Warded(left.ID()) {
			t.Errorf("%s should be warded", left.Name())
		}
		if !h.Game().Warded(right.ID()) {
			t.Errorf("%s should be warded", right.Name())
		}
	})
}
