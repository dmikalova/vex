package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Barehanded
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Put each artifact on top of its owner's deck.
func TestBarehanded(t *testing.T) {
	t.Run("returns each artifact to the top of its owner's deck", func(t *testing.T) {
		var mine, theirs ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				Hand:   ct.Cards(Barehanded),
				InPlay: ct.Cards(ct.Bind(&mine, ct.Artifact(ct.OfHouse(card.House.Brobnar)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&theirs, ct.Artifact(ct.OfHouse(card.House.Untamed)))),
			},
		})

		h.P1.Play(Barehanded)

		h.Expect(mine).At(ct.Deck)
		h.Expect(theirs).At(ct.Deck)
		h.P1.ExpectAmber(1) // the Æmber bonus pip resolves on play
		if top := h.Game().State.Deck[0].IDs[0]; top != mine.ID() {
			t.Errorf("player 0 deck top = %d, want mine (%d)", top, mine.ID())
		}
		if top := h.Game().State.Deck[1].IDs[0]; top != theirs.ID() {
			t.Errorf("player 1 deck top = %d, want theirs (%d)", top, theirs.ID())
		}
	})

	t.Run("orders several artifacts without pausing", func(t *testing.T) {
		var artA, artB ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(Barehanded),
				InPlay: ct.Cards(
					ct.Bind(&artA, ct.Artifact(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&artB, ct.Artifact(ct.OfHouse(card.House.Brobnar))),
				),
			},
		})

		h.P1.Play(Barehanded) // two artifacts to order, but no click is needed
		h.Expect(artA).At(ct.Deck)
		h.Expect(artB).At(ct.Deck)
	})

	t.Run("Player.Order controls the stacking order", func(t *testing.T) {
		var artA, artB ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(Barehanded),
				InPlay: ct.Cards(
					ct.Bind(&artA, ct.Artifact(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&artB, ct.Artifact(ct.OfHouse(card.House.Brobnar))),
				),
			},
		})

		// Returning artB first stacks artA last, so artA ends up on top.
		h.P1.Order(artB, artA)
		h.P1.Play(Barehanded)

		if top := h.Game().State.Deck[0].IDs[0]; top != artA.ID() {
			t.Errorf("deck top = %d, want artA (%d)", top, artA.ID())
		}
	})
}
