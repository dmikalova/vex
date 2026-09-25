package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Bad Penny
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Human • Thief
//
//	Destroyed: Put Bad Penny into its owner's hand.
func TestBadPenny(t *testing.T) {
	t.Run("returns to its owner's hand when destroyed", func(t *testing.T) {
		var penny ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&penny, BadPenny)),
			},
		})

		h.Game().DestroyEach(0, []engine.LocalID{penny.ID()})

		h.Expect(penny).At(ct.Hand) // returned to hand instead of the discard
	})
}
