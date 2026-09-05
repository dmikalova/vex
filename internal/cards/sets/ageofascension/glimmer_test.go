package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Glimmer
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Faerie
//
//	Alpha.
//	Play: Put a card from your discard pile into your hand.
func TestGlimmer(t *testing.T) {
	t.Run("returns a card from the discard pile to your hand when played", func(t *testing.T) {
		var stashed ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Untamed,
				Hand:    ct.Cards(Glimmer),
				Discard: ct.Cards(ct.Bind(&stashed, ct.Creature())),
			},
		})

		h.P1.Play(Glimmer)

		h.Expect(stashed).At(ct.Hand)
	})
}
