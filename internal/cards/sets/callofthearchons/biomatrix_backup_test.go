package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Biomatrix Backup
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Common
//	Bonus:  Æmber
//
//	This creature gains, "Destroyed: Put this creature into its owner's archives."
func TestBiomatrixBackup(t *testing.T) {
	t.Run("relocates the destroyed host to its owner's archives", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
						BiomatrixBackup,
					),
				),
			},
		})

		h.Game().DestroyEach(0, []engine.LocalID{host.ID()})

		h.Expect(host).At(ct.Archives)           // host archived instead of discarded
		h.Expect(BiomatrixBackup).At(ct.Discard) // the upgrade itself is discarded
	})
}
