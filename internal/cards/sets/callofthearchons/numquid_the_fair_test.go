package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Numquid the Fair
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human
//
//	Play: Destroy an enemy creature -> if you are overwhelmed, repeat this effect.
func TestNumquidTheFair(t *testing.T) {
	t.Run("destroys enemy creatures while overwhelmed", func(t *testing.T) {
		var foe1, foe2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(NumquidTheFair)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe1, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				ct.Bind(&foe2, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			)},
		})

		// Numquid enters: P1 controls 1, P2 controls 2 -> overwhelmed. One destroy
		// brings it to parity, so the loop stops after a single destruction.
		h.P1.Play(NumquidTheFair)
		h.P1.ClickCard(foe1)

		h.Expect(foe1).At(ct.Discard)
		h.Expect(foe2).At(ct.PlayArea)
	})
}
