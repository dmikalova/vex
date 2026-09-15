package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mab the Mad
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Bonus:  Æmber
//	Traits: Faerie
//
//	Reap: Shuffle Mab the Mad into its owner's deck.
func TestMabTheMad(t *testing.T) {
	t.Run("shuffles itself into your deck when it reaps", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(MabTheMad)},
		})

		h.P1.Reap(MabTheMad)

		h.Expect(MabTheMad).At(ct.Deck)
	})
}
