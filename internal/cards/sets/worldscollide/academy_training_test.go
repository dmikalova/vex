package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Academy Training
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Rare
//
//	This Creature belongs to Logos and this Creature gains "Reap: Draw a card."
func TestAcademyTraining(t *testing.T) {
	t.Run("host belongs to Logos and reaps to draw a card", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar))),
						AcademyTraining,
					),
				),
				Deck: ct.Cards(ct.Creature()),
			},
		})
		host.Ready()

		if got := h.Game().House(host.ID()); got != card.House.Logos {
			t.Errorf("house = %v, want Logos", got)
		}

		h.P1.Reap(host)

		h.P1.ExpectAmber(1)
		if got := len(h.Game().Hand(0)); got != 1 {
			t.Errorf("hand = %d cards, want 1", got)
		}
	})
}
