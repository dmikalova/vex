package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lomir Flamefist
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Giant
//
//	Play: If your opponent has 7 Æmber or more, your opponent loses 2 Æmber.
func TestLomirFlamefist(t *testing.T) {
	t.Run("drains 2 Æmber when the opponent has 7 or more", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(LomirFlamefist)},
			P2: ct.Side{Amber: 7},
		})

		h.P1.Play(LomirFlamefist)

		h.P2.ExpectAmber(5)
	})

	t.Run("does nothing when the opponent has fewer than 7", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(LomirFlamefist)},
			P2: ct.Side{Amber: 6},
		})

		h.P1.Play(LomirFlamefist)

		h.P2.ExpectAmber(6)
	})
}
