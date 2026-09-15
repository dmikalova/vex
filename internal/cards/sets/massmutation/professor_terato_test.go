package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Professor Terato
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Mutant • Scientist
//
//	Each Mutant creature gains, "Reap: Draw a card."
func TestProfessorTerato(t *testing.T) {
	t.Run("a Mutant creature draws a card when it reaps", func(t *testing.T) {
		var mutant, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ProfessorTerato,
					ct.Bind(&mutant, ct.Creature(
						ct.OfHouse(card.House.Logos),
						ct.Traits(card.Traits.Mutant),
					)),
				),
				Deck: ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
		})

		h.P1.Reap(mutant)

		h.Expect(top).At(ct.Hand)
	})

	t.Run("a non-Mutant creature draws nothing when it reaps", func(t *testing.T) {
		var plain, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ProfessorTerato,
					ct.Bind(&plain, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
				Deck: ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
		})

		h.P1.Reap(plain)

		h.Expect(top).At(ct.Deck)
	})
}
