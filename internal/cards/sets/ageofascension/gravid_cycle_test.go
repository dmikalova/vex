package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gravid Cycle
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Omega.
//	Play: Put a card from your discard pile into your hand.
func TestGravidCycle(t *testing.T) {
	t.Run("returns a card from the discard pile to your hand when played", func(t *testing.T) {
		var stashed ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Untamed,
				Hand:    ct.Cards(GravidCycle),
				Discard: ct.Cards(ct.Bind(&stashed, ct.Creature())),
				// Omega ends the turn, whose draw step recycles the discard when
				// the deck is empty; stock the deck so the returned card stays put.
				Deck: ct.Cards(
					ct.Creature(),
					ct.Creature(),
					ct.Creature(),
					ct.Creature(),
					ct.Creature(),
					ct.Creature(),
				),
			},
			P2: ct.Side{},
		})

		h.P1.Play(GravidCycle)

		h.Expect(stashed).At(ct.Hand)
	})
}
