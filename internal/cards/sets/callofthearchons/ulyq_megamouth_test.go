package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ulyq Megamouth
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Martian • Scientist
//
//	Fight/Reap: Use a friendly non-Mars creature.
func TestUlyqMegamouth(t *testing.T) {
	t.Run("uses a friendly non-Mars creature when it reaps", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(UlyqMegamouth, ct.Creature(ct.OfHouse(card.House.Brobnar))),
			},
		})

		h.P1.Reap(UlyqMegamouth)

		h.P1.ExpectAmber(2) // Ulyq's reap plus the used creature's reap
	})
}
