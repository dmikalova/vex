package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sergeant Zakiel
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human • Knight
//
//	Play: You may ready and fight with a neighboring creature.
func TestSergeantZakiel(t *testing.T) {
	t.Run("may ready and fight with a neighboring creature", func(t *testing.T) {
		var neighbor, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(SergeantZakiel),
				InPlay: ct.Cards(
					ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
			},
		})

		h.P1.Play(SergeantZakiel)
		h.P1.ClickCard(neighbor)

		h.Expect(foe).At(ct.Discard)
		h.Expect(neighbor).Exhausted()
	})
}
