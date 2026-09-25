package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Charybdis
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Connected
//	Power:  7
//	Traits: Beast
//
//	Each enemy creature gains, "Before Fight: Lose 1 Æmber."
func TestCharybdis(t *testing.T) {
	t.Run("an enemy creature's controller loses 1 Æmber before it fights", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Amber: 2,
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(8))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(Charybdis)},
		})

		h.P1.Fight(enemy, Charybdis)

		h.P1.ExpectAmber(1)
	})
}
