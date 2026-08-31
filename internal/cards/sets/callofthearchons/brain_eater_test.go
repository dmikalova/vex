package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Brain Eater
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Cyborg • Beast
//
//	After a creature is destroyed fighting Brain Eater, draw a card.
func TestBrainEater(t *testing.T) {
	t.Run("draws a card after a creature is destroyed fighting it", func(t *testing.T) {
		var prey, drawn ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(BrainEater),
				Deck: ct.Cards(
					ct.Bind(&drawn, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(1))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&prey, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
			},
		})

		h.P1.Fight(BrainEater, prey)

		h.Expect(prey).At(ct.Discard) // the 1-power prey is destroyed
		h.Expect(drawn).At(ct.Hand)   // Brain Eater draws a card
	})
}
