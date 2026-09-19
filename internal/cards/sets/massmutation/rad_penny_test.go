package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Rad Penny
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Mutant • Thief
//
//	Play: Steal 1 Æmber.
//	Destroyed: Shuffle Rad Penny into its owner's deck.
func TestRadPenny(t *testing.T) {
	t.Run("steals 1 Æmber when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(RadPenny),
			},
			P2: ct.Side{Amber: 2},
		})

		h.P1.Play(RadPenny)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})

	t.Run("shuffles itself into your deck when destroyed", func(t *testing.T) {
		var penny ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&penny, RadPenny)),
			},
		})

		h.Game().DestroyEach(0, []engine.LocalID{penny.ID()})

		h.Expect(penny).At(ct.Deck)
	})
}
