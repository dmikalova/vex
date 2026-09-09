package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mega Narp
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  10
//	Armor:  1
//	Traits: Giant
//
//	Each neighboring creature cannot reap.
func TestMegaNarp(t *testing.T) {
	t.Run("bars its neighbors from reaping but not distant friends", func(t *testing.T) {
		var left, right, far ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Brobnar))),
					MegaNarp,
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&far, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				),
			},
		})

		h.P1.ExpectCannotUseTo(left, card.UseKind.Reap)
		h.P1.ExpectCannotUseTo(right, card.UseKind.Reap)

		h.P1.Reap(far)
		h.P1.ExpectAmber(1)
	})
}
