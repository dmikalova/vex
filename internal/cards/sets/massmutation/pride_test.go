package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Pride
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Traits: Demon • Sin
//
//	Reap: Ward each friendly Sin creature.
func TestPride(t *testing.T) {
	t.Run("wards each friendly Sin creature and no others", func(t *testing.T) {
		var pride, otherSin, nonSin ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ct.Bind(&pride, Pride),
					ct.Bind(
						&otherSin,
						ct.Creature(
							ct.OfHouse(card.House.Dis),
							ct.Traits(card.Traits.Sin),
							ct.Power(3),
						),
					),
					ct.Bind(&nonSin, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(pride)

		if !h.Game().Warded(pride.ID()) {
			t.Error("Pride should ward itself")
		}
		if !h.Game().Warded(otherSin.ID()) {
			t.Error("friendly Sin creature should be warded")
		}
		if h.Game().Warded(nonSin.ID()) {
			t.Error("non-Sin creature should not be warded")
		}
	})
}
