package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Grumpus Tamer
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Giant
//
//	Reap: Search your deck and discard pile for a War Grumpus, reveal it, and put it into your hand.
func TestGrumpusTamer(t *testing.T) {
	t.Run("tutors a War Grumpus from the deck into hand", func(t *testing.T) {
		var grumpus ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(GrumpusTamer),
				Deck:   ct.Cards(ct.Bind(&grumpus, WarGrumpus)),
			},
		})

		h.P1.Reap(GrumpusTamer)
		h.Expect(grumpus).At(ct.Hand)
	})

	t.Run("tutors a War Grumpus from the discard pile into hand", func(t *testing.T) {
		var grumpus ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Brobnar,
				InPlay:  ct.Cards(GrumpusTamer),
				Deck:    ct.Cards(ct.Creature(ct.OfHouse(card.House.Brobnar))),
				Discard: ct.Cards(ct.Bind(&grumpus, WarGrumpus)),
			},
		})

		h.P1.Reap(GrumpusTamer)
		h.Expect(grumpus).At(ct.Hand)
	})
}
