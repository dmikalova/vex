package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Horseman of Famine
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Connected
//	Power:  5
//	Traits: Horseman • Spirit
//
//	Play/Fight/Reap: Destroy the least powerful creature.
func TestHorsemanOfFamine(t *testing.T) {
	t.Run("destroys the least powerful creature when played", func(t *testing.T) {
		var weak, strong ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(HorsemanOfFamine)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&weak, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
					ct.Bind(&strong, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.P1.Play(HorsemanOfFamine)

		h.Expect(weak).At(ct.Discard)
		h.Expect(strong).At(ct.PlayArea)
	})
}
