package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Quant
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Scientist
//
//	Reap: Play a non-Logos tactic.
func TestQuant(t *testing.T) {
	t.Run("reap plays a non-Logos tactic from hand", func(t *testing.T) {
		var brobnar, logos ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(Quant),
				Hand: ct.Cards(
					ct.Bind(&brobnar, ct.Tactic(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&logos, ct.Tactic(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		h.P1.Reap(Quant)

		h.Expect(brobnar).At(ct.Discard)
		h.Expect(logos).At(ct.Hand)
	})

	t.Run("reap does nothing with only a Logos tactic in hand", func(t *testing.T) {
		var logos ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(Quant),
				Hand: ct.Cards(
					ct.Bind(&logos, ct.Tactic(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		h.P1.Reap(Quant)

		h.Expect(logos).At(ct.Hand)
	})
}
