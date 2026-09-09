package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dexus
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Traits: Demon
//
//	After your opponent plays a card, if it is on the right flank, your opponent loses 1 Æmber.
func TestDexus(t *testing.T) {
	t.Run("drains the opponent when they play onto their right flank", func(t *testing.T) {
		var beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(Dexus),
			},
			P2: ct.Side{
				House: card.House.Brobnar,
				Amber: 3,
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3)),
				),
				Hand: ct.Cards(
					ct.Bind(&beast, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		// The harness plays new creatures onto the right flank.
		h.P2.Play(beast)

		h.P2.ExpectAmber(2)
	})
}
