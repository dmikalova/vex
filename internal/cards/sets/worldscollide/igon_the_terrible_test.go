package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Igon the Terrible
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  8
//	Traits: Giant
//
//	Play: If Igon the Green has not been purged, destroy Igon the Terrible.
//	Fight: Steal 1 Æmber.
func TestIgonTheTerrible(t *testing.T) {
	t.Run("destroys itself when Igon the Green has not been purged", func(t *testing.T) {
		var terrible ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(ct.Bind(&terrible, IgonTheTerrible)),
			},
		})

		h.P1.Play(terrible)

		h.Expect(terrible).At(ct.Discard)
	})

	t.Run("survives when Igon the Green has been purged", func(t *testing.T) {
		var terrible, green ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Brobnar,
				Hand:    ct.Cards(ct.Bind(&terrible, IgonTheTerrible)),
				Discard: ct.Cards(ct.Bind(&green, IgonTheGreen)),
			},
		})
		// Move Igon the Green from the discard pile into the purge pile so the play
		// condition sees it as already purged.
		h.Game().PurgeFromDiscard(0, green.ID())

		h.P1.Play(terrible)

		h.Expect(terrible).At(ct.PlayArea)
	})

	t.Run("steals 1 Æmber when it fights", func(t *testing.T) {
		var terrible, target, green ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Brobnar,
				InPlay:  ct.Cards(ct.Bind(&terrible, IgonTheTerrible)),
				Discard: ct.Cards(ct.Bind(&green, IgonTheGreen)),
			},
			P2: ct.Side{
				Amber:  2,
				InPlay: ct.Cards(ct.Bind(&target, ct.Creature(ct.Power(1)))),
			},
		})
		h.Game().PurgeFromDiscard(0, green.ID())

		h.P1.Fight(terrible, target)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})
}
