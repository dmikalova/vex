package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Memory Chip
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Archive a card from your hand.
func TestMemoryChip(t *testing.T) {
	t.Run("archives a card from your hand after you choose a house", func(t *testing.T) {
		var other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(MemoryChip),
				Hand:   ct.Cards(ct.Bind(&other, ct.Creature())),
			},
			P2: ct.Side{},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Logos)
		h.P1.ClickOption("No")

		h.Expect(other).At(ct.Archives)
	})
}
