package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tunk
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Robot
//
//	After you play a Mars creature, fully heal Tunk.
func TestTunk(t *testing.T) {
	t.Run("fully heals itself after you play another Mars creature", func(t *testing.T) {
		var tunk, mars ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ct.Bind(&tunk, Tunk)),
				Hand:   ct.Cards(ct.Bind(&mars, ct.Creature(ct.OfHouse(card.House.Mars)))),
			},
		})

		tunk.Damaged(4)
		h.P1.Play(mars)

		h.Expect(tunk).Damage(0)
	})
}
