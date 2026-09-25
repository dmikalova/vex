package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Sinestra
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Demon
//
//	After your opponent plays a creature on their left flank, your opponent loses 1 Æmber.
func TestSinestra(t *testing.T) {
	t.Run("drains the opponent when they play onto their left flank", func(t *testing.T) {
		var beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(Sinestra),
			},
			P2: ct.Side{
				House: card.House.Brobnar,
				Amber: 3,
				Hand: ct.Cards(
					ct.Bind(&beast, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		// A lone creature enters on both flanks, so it counts as the left flank.
		h.P2.Play(beast)

		h.P2.ExpectAmber(2)
	})

	t.Run("ignores a creature played onto the right flank", func(t *testing.T) {
		var beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(Sinestra),
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
		// The harness plays new creatures onto the right flank, past the existing one.
		h.P2.Play(beast)

		h.P2.ExpectAmber(3)
	})
}
