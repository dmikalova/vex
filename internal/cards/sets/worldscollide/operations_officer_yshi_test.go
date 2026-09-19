package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Operations Officer Yshi
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  2
//	Traits: Spirit
//
//	Taunt.
//	Each neighboring creature gains, "Reap: This creature captures 1 Æmber from your opponent."
//	Each neighboring creature gains, "Fight: This creature captures 1 Æmber from your opponent."
func TestOperationsOfficerYshi(t *testing.T) {
	t.Run("a neighbor captures 1 Æmber when it reaps", func(t *testing.T) {
		var neighbor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
					OperationsOfficerYshi,
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Reap(neighbor)

		h.Expect(neighbor).AmberOn(1)
		h.P1.ExpectAmber(1) // the reap itself
		h.P2.ExpectAmber(2) // 1 captured off the opponent's pool
	})

	t.Run("a distant friendly creature does not capture when it reaps", func(t *testing.T) {
		var far ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					OperationsOfficerYshi,
					ct.Creature(ct.OfHouse(card.House.StarAlliance)),
					ct.Bind(&far, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Reap(far)

		h.Expect(far).AmberOn(0)
		h.P2.ExpectAmber(3) // untouched
	})
}
