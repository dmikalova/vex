package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hidden Stash
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Archive a card from your hand.
func TestHiddenStash(t *testing.T) {
	t.Run("archives a card from your hand when played", func(t *testing.T) {
		var other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(HiddenStash, ct.Bind(&other, ct.Creature())),
			},
		})

		h.P1.Play(HiddenStash)

		h.Expect(other).At(ct.Archives)
	})
}
