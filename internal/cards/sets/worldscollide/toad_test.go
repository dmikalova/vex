package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Toad
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Connected
//	Power:  1
//	Traits: Beast
//
//	Toad cannot reap.
func TestToad(t *testing.T) {
	t.Run("cannot reap", func(t *testing.T) {
		var toad ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&toad, Toad)),
			},
		})

		h.P1.ExpectCannotUseTo(toad, card.UseKind.Reap)
	})
}
