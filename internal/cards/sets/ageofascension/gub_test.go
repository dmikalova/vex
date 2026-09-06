package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gub
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Demon
//
//	Gub gains +5 power and taunt while it is not on a flank.
func TestGub(t *testing.T) {
	t.Run("on a flank Gub has its base power", func(t *testing.T) {
		var gub ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(ct.Bind(&gub, Gub))},
		})
		h.Expect(gub).Power(1)
	})

	t.Run("off a flank Gub gets +5 power", func(t *testing.T) {
		var gub ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ct.Creature(ct.Power(3)),
					ct.Bind(&gub, Gub),
					ct.Creature(ct.Power(3)),
				),
			},
		})
		h.Expect(gub).Power(6)
	})
}
