package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Chaos Portal
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Choose a house. Reveal the top card of your deck. If it is of the chosen house, play it.
func TestChaosPortal(t *testing.T) {
	t.Run("plays the top card when it is of the chosen house", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ChaosPortal),
				Deck: ct.Cards(
					ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
		})

		h.P1.UseAction(ChaosPortal)
		h.P1.ExpectPrompt("Choose a house").Source("Chaos Portal")
		h.P1.ClickOption("Mars")

		h.Expect(top).At(ct.PlayArea).Exhausted()
	})

	t.Run("leaves the revealed card on top when it is not of the chosen house", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ChaosPortal),
				Deck: ct.Cards(
					ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
		})

		h.P1.UseAction(ChaosPortal)
		h.P1.ExpectPrompt("Choose a house").Source("Chaos Portal")
		h.P1.ClickOption("Logos")

		h.Expect(top).At(ct.Deck)
		if deck := h.Game().Deck(0); len(deck) != 1 || deck[0] != top.ID() {
			t.Errorf("deck = %v, want the revealed card still on top", deck)
		}
	})
}
