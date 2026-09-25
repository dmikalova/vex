package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Resurgence
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Put a creature from your discard pile into your hand. If it is a Mutant creature, put another creature from your discard pile into your hand.
//	Enhance Draw.
func TestResurgence(t *testing.T) {
	t.Run("a Mutant first return brings back a second creature", func(t *testing.T) {
		var mutant, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(Resurgence),
				Discard: ct.Cards(
					ct.Bind(&mutant, ct.Creature(
						ct.OfHouse(card.House.Untamed), ct.Traits(card.Traits.Mutant))),
					ct.Bind(&other, ct.Creature(
						ct.OfHouse(card.House.Untamed), ct.Traits(card.Traits.Beast))),
				),
			},
		})

		h.P1.Play(Resurgence)
		h.P1.ClickCard(mutant) // sole remaining creature auto-returns second

		h.Expect(mutant).At(ct.Hand)
		h.Expect(other).At(ct.Hand)
	})

	t.Run("a non-Mutant first return brings back nothing else", func(t *testing.T) {
		var beast, mutant ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(Resurgence),
				Discard: ct.Cards(
					ct.Bind(&beast, ct.Creature(
						ct.OfHouse(card.House.Untamed), ct.Traits(card.Traits.Beast))),
					ct.Bind(&mutant, ct.Creature(
						ct.OfHouse(card.House.Untamed), ct.Traits(card.Traits.Mutant))),
				),
			},
		})

		h.P1.Play(Resurgence)
		h.P1.ClickCard(beast)

		h.Expect(beast).At(ct.Hand)
		h.Expect(mutant).At(ct.Discard)
	})
}
