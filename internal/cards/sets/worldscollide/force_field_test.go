package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Force Field
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "Reap: Ward this creature."
func TestForceField(t *testing.T) {
	t.Run("wards its host when the host reaps", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&host,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						ForceField,
					),
				),
			},
		})

		h.P1.Reap(host)

		if !h.Game().Warded(host.ID()) {
			t.Error("the host should be warded after reaping")
		}
	})
}
