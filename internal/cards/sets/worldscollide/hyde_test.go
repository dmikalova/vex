package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hyde
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Human • Scientist
//
//	Reap: Draw a card. If you control Velum, draw a card.
//	Destroyed: Archive Velum from your discard pile -> archive Hyde from play.
func TestHyde(t *testing.T) {
	t.Run("reaps to draw 1 without Velum", func(t *testing.T) {
		var hyde ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&hyde, Hyde)),
				Deck:   ct.Cards(ct.Creature(), ct.Creature()),
			},
		})
		hyde.Ready()

		h.P1.Reap(hyde)

		if got := len(h.Game().Hand(0)); got != 1 {
			t.Errorf("hand = %d, want 1", got)
		}
	})

	t.Run("reaps to draw 2 while controlling Velum", func(t *testing.T) {
		var hyde ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&hyde, Hyde), Velum),
				Deck:   ct.Cards(ct.Creature(), ct.Creature()),
			},
		})
		hyde.Ready()

		h.P1.Reap(hyde)

		if got := len(h.Game().Hand(0)); got != 2 {
			t.Errorf("hand = %d, want 2", got)
		}
	})

	t.Run("destroyed archives Velum from discard and archives Hyde", func(t *testing.T) {
		var hyde, velum, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Logos,
				InPlay:  ct.Cards(ct.Bind(&hyde, Hyde)),
				Discard: ct.Cards(ct.Bind(&velum, Velum)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6))))},
		})

		h.P1.Fight(hyde, enemy)

		h.Expect(velum).At(ct.Archives)
		h.Expect(hyde).At(ct.Archives)
	})
}
