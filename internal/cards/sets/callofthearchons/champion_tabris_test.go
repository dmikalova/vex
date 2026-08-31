package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Champion Tabris
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Armor:  2
//	Traits: Human • Knight
//
//	Fight: Champion Tabris captures 1 Æmber from your opponent.
func TestChampionTabris(t *testing.T) {
	t.Run("captures 1 Æmber when it fights", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(ChampionTabris)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
				Amber: 2,
			},
		})

		h.P1.Fight(ChampionTabris, enemy)

		h.Expect(ChampionTabris).AmberOn(1)
		h.P2.ExpectAmber(1)
	})
}
