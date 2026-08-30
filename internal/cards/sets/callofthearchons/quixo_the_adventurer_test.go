package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Quixo the "Adventurer"
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Scientist
//
//	Skirmish.
//	Fight: Draw a card.
func TestQuixoTheAdventurer(t *testing.T) {
	t.Run("draws a card when it fights", func(t *testing.T) {
		var foe, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(QuixoTheAdventurer),
				Deck:   ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))))},
		})

		h.P1.Fight(QuixoTheAdventurer, foe)

		h.Expect(top).At(ct.Hand)
		h.Expect(QuixoTheAdventurer).At(ct.PlayArea).Damage(0) // skirmish
	})
}
