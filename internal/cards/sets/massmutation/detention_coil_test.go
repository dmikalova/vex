package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Detention Coil
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature cannot fight.
func TestDetentionCoil(t *testing.T) {
	t.Run("bars its host from fighting but not reaping", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
						DetentionCoil,
					),
				),
			},
		})

		if !h.Game().CannotBeUsedTo(host.ID(), engine.FightUse) {
			t.Error("host with Detention Coil should not be able to fight")
		}
		if h.Game().CannotBeUsedTo(host.ID(), engine.ReapUse) {
			t.Error("host with Detention Coil should still be able to reap")
		}
	})
}
