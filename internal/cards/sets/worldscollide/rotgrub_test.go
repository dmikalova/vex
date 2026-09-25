package worldscollide

import (
	"slices"
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Rotgrub
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Beast
//
//	Play: Your opponent loses 1 Æmber.
//	Reap: Archive Rotgrub.
func TestRotgrub(t *testing.T) {
	t.Run("opponent loses 1 Æmber when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Rotgrub),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(Rotgrub)

		h.P2.ExpectAmber(2)
	})

	t.Run("archives itself when it reaps", func(t *testing.T) {
		var grub ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&grub, Rotgrub)),
			},
		})

		h.P1.Reap(grub)

		h.Expect(grub).At(ct.Archives)
		if !slices.Contains(h.Game().Archives(0), grub.ID()) {
			t.Error("Rotgrub should be archived after it reaps")
		}
	})
}
