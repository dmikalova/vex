package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Access Denied
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Common
//	Bonus:  Æmber
//
//	This creature cannot reap.
func TestAccessDenied(t *testing.T) {
	t.Run("bars its host from reaping but not fighting", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
						AccessDenied,
					),
				),
			},
		})

		if !h.Game().CannotBeUsedTo(host.ID(), engine.ReapUse) {
			t.Error("host with Access Denied should not be able to reap")
		}
		if h.Game().CannotBeUsedTo(host.ID(), engine.FightUse) {
			t.Error("host with Access Denied should still be able to fight")
		}
	})
}
