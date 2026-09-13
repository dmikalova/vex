package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Fila the Researcher
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Human • Scientist
//
//	Elusive.
//	After a Creature is played adjacent to Fila the Researcher, draw a card.
func TestFilaTheResearcher(t *testing.T) {
	t.Run("draws a card when a creature is played next to it", func(t *testing.T) {
		var fila, newbie, drawn ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&fila, FilaTheResearcher)),
				Hand: ct.Cards(
					ct.Bind(&newbie, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
				Deck: ct.Cards(
					ct.Bind(&drawn, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(1))),
				),
			},
		})

		h.P1.Play(newbie)

		h.Expect(newbie).At(ct.PlayArea)
		h.Expect(drawn).At(ct.Hand) // Fila drew the top card
	})
}
