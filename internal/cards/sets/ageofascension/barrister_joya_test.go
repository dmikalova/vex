package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Barrister Joya
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Armor:  1
//	Traits: Human • Knight
//
//	Enemy creatures cannot reap.
func TestBarristerJoya(t *testing.T) {
	t.Run("enemy creatures cannot reap", func(t *testing.T) {
		var joya, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&joya, BarristerJoya)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)

		h.P2.ExpectCannotUseTo(enemy, card.UseKind.Reap)
	})

	t.Run("friendly creatures still reap", func(t *testing.T) {
		var joya, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Bind(&joya, BarristerJoya),
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(ally)

		h.P1.ExpectAmber(1)
	})
}
