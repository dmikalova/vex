package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dust Pixie
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Æmber:  2
//	Traits: Faerie
func TestDustPixie(t *testing.T) {
	t.Run("gains 2 Æmber when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(DustPixie)},
		})

		h.P1.Play(DustPixie)

		h.Expect(DustPixie).At(ct.PlayArea).Power(1)
		h.P1.ExpectAmber(2)
	})
}
