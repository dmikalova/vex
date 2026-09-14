package anomalyexpansion

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Timequake
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Shuffle each friendly card in play into your deck. Draw a card for each card shuffled into your deck this way.
func TestTimequake(t *testing.T) {
	t.Run(
		"shuffles every friendly card in play into the deck and redraws that many",
		func(t *testing.T) {
			var tq, creature, artifact ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Brobnar,
					Hand:  ct.Cards(ct.Bind(&tq, Timequake)),
					InPlay: ct.Cards(
						ct.Bind(&creature, ct.Creature()),
						ct.Bind(&artifact, ct.Artifact()),
					),
					Deck: ct.Cards(ct.Creature(), ct.Creature(), ct.Creature()),
				},
				P2: ct.Side{InPlay: ct.Cards(ct.Creature())},
			})

			deckBefore := h.Game().State.Deck[0].Count
			handBefore := h.Game().State.Hand[0].Count

			h.P1.Play(tq)

			// The friendly creature and artifact left play (shuffled into the deck,
			// though the redraw may pull either back into hand).
			if creature.Location() == ct.PlayArea {
				t.Error("creature should have left play")
			}
			if artifact.Location() == ct.PlayArea {
				t.Error("artifact should have left play")
			}

			// Two cards were shuffled in and two were drawn out, so the deck is net
			// unchanged and the hand holds two more than after the tactic left it.
			if got := h.Game().State.Deck[0].Count; got != deckBefore {
				t.Errorf("deck count = %d, want %d", got, deckBefore)
			}
			if got := h.Game().State.Hand[0].Count; got != handBefore-1+2 {
				t.Errorf("hand count = %d, want %d", got, handBefore-1+2)
			}
		},
	)
}
