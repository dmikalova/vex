package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Velum
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Connected
//	Power:  2
//	Traits: Human • Scientist
//
//	Reap: Archive a card from your hand. If you control Hyde, archive a card from your hand.
//	Destroyed: Archive Hyde from your discard pile -> archive Velum from play.
func TestVelum(t *testing.T) {
	t.Run("reaps to archive 1 without Hyde", func(t *testing.T) {
		var velum, h1, h2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&velum, Velum)),
				Hand: ct.Cards(
					ct.Bind(&h1, ct.Creature()),
					ct.Bind(&h2, ct.Creature()),
				),
			},
		})
		velum.Ready()

		h.P1.Reap(velum)
		h.P1.ClickCard(h1)

		h.Expect(h1).At(ct.Archives)
		h.Expect(h2).At(ct.Hand)
	})

	t.Run("reaps to archive 2 while controlling Hyde", func(t *testing.T) {
		var velum, h1, h2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&velum, Velum), Hyde),
				Hand: ct.Cards(
					ct.Bind(&h1, ct.Creature()),
					ct.Bind(&h2, ct.Creature()),
				),
			},
		})
		velum.Ready()

		h.P1.Reap(velum)
		h.P1.ClickCard(h1)

		h.Expect(h1).At(ct.Archives)
		h.Expect(h2).At(ct.Archives)
	})

	t.Run("destroyed archives Hyde from discard and archives Velum", func(t *testing.T) {
		var velum, hyde, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Logos,
				InPlay:  ct.Cards(ct.Bind(&velum, Velum)),
				Discard: ct.Cards(ct.Bind(&hyde, Hyde)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6))))},
		})

		h.P1.Fight(velum, enemy)

		h.Expect(hyde).At(ct.Archives)
		h.Expect(velum).At(ct.Archives)
	})
}
