package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Book of leQ
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Reveal the top card of your deck. If it is a non-Star Alliance card, its house becomes your active house. Otherwise, end your turn.
func TestBookOfLeQ(t *testing.T) {
	t.Run("a non-Star Alliance top card becomes the active house", func(t *testing.T) {
		var brute ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					BookOfLeQ,
					ct.Bind(&brute, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
				),
				Deck: ct.Cards(ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(2))),
			},
		})

		// While Star Alliance is the active house, the Brobnar creature cannot reap.
		h.P1.ExpectCannotUse(brute)

		h.P1.UseAction(BookOfLeQ)

		// The revealed Brobnar top card made Brobnar the active house, so the Brobnar
		// creature can now reap and gain Æmber.
		if got := h.Game().State.ActiveHouse; got != card.House.Brobnar {
			t.Errorf("active house = %v, want Brobnar", got)
		}
		h.P1.Reap(brute)
		h.P1.ExpectAmber(1)
	})

	t.Run("a Star Alliance top card ends the turn with no cleanup", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(BookOfLeQ),
				Hand:   ct.Cards(ct.Creature(ct.OfHouse(card.House.StarAlliance))),
				Deck: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(2)),
					ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(2)),
				),
			},
		})

		handBefore := len(h.Game().Hand(0))

		h.P1.UseAction(BookOfLeQ)

		// The revealed Star Alliance top card ends P1's turn in place: the turn stops
		// the moment it resolves, so no ready, draw, or end-of-turn step runs.
		if got := h.Game().Phase(); got != engine.PhaseEndOfTurn {
			t.Errorf("phase after use = %v, want PhaseEndOfTurn (turn ended)", got)
		}
		// No ready step ran: Book of leQ was exhausted by its own Action and stays
		// exhausted rather than readying.
		h.Expect(BookOfLeQ).Exhausted()
		// No draw step ran: P1's hand is not refilled toward a full hand of six.
		if got := len(h.Game().Hand(0)); got != handBefore {
			t.Errorf("hand size after use = %d, want %d (no draw)", got, handBefore)
		}
	})
}
