package anomalyexpansion

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ghostform
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Special
//	Æmber:  1
//
//	This Creature gains invulnerable.
//	This Creature gains, "Fight/Reap: Archive Ghostform."
func TestGhostform(t *testing.T) {
	t.Run("the host reaping archives Ghostform off it", func(t *testing.T) {
		var ghost, host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
				),
				Hand: ct.Cards(ct.Bind(&ghost, Ghostform)),
			},
		})

		h.P1.Play(ghost) // the lone host auto-attaches

		h.P1.Reap(host)

		h.Expect(ghost).At(ct.Archives)
		h.Expect(host).At(ct.PlayArea)
	})
}
