package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Igon the Green
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Giant
//
//	Destroyed: Purge Igon the Green, and put an Igon the Terrible from your discard pile into your hand.
func TestIgonTheGreen(t *testing.T) {
	t.Run("purges itself and returns Igon the Terrible from discard to hand", func(t *testing.T) {
		igonTerrible := engine.NewCard(
			"Igon the Terrible",
			engine.Brobnar, engine.Creature, engine.Rare,
			engine.WithPower(4),
		)
		var terrible, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Brobnar,
				InPlay:  ct.Cards(IgonTheGreen),
				Discard: ct.Cards(ct.Bind(&terrible, igonTerrible)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(4))))},
		})

		// Igon the Green (power 4) trades with a power-4 enemy and is destroyed.
		h.P1.Fight(IgonTheGreen, foe)

		h.Expect(IgonTheGreen).At(ct.Purge)
		h.Expect(terrible).At(ct.Hand)
	})
}
