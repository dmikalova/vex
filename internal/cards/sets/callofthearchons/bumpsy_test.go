package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bumpsy
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Play: Your opponent loses 1 Æmber.
func TestBumpsy(t *testing.T) {
	t.Run("makes the opponent lose 1 Æmber when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(Bumpsy),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(Bumpsy)

		h.Expect(Bumpsy).Power(5)
		h.P2.ExpectAmber(2)
	})
}
