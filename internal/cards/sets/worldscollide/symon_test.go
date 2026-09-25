package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Symon
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Alien • Thief
//
//	Skirmish.
//	Fight: Put the creature Symon fought on top of its owner's deck.
func TestSymon(t *testing.T) {
	t.Run("puts the creature it fought on top of its owner's deck", func(t *testing.T) {
		var symon, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&symon, Symon)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
		})

		h.P1.Fight(symon, foe)

		h.Expect(foe).At(ct.Deck) // top of P2's deck, not the play area
		h.Expect(symon).At(ct.PlayArea)
	})
}
