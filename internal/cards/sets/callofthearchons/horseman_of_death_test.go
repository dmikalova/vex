package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Horseman of Death
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Connected
//	Power:  5
//	Traits: Horseman • Spirit
//
//	Play: Put each Horseman trait creature from your discard pile into your hand.
func TestHorsemanOfDeath(t *testing.T) {
	t.Run("returns each Horseman creature from your discard pile to hand", func(t *testing.T) {
		var rider, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(HorsemanOfDeath),
				Discard: ct.Cards(
					ct.Bind(
						&rider,
						ct.Creature(
							ct.OfHouse(card.House.Sanctum),
							ct.Traits(card.Traits.Horseman),
						),
					),
					ct.Bind(
						&other,
						ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Traits(card.Traits.Human)),
					),
				),
			},
		})

		h.P1.Play(HorsemanOfDeath)

		h.Expect(rider).At(ct.Hand)
		h.Expect(other).At(ct.Discard)
	})
}
