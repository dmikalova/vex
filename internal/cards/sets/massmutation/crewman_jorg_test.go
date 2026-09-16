package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Crewman Jorg
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Thief
//
//	Action: If Crewman Jorg has no Star Alliance neighbor, steal 1 Æmber.
//	Enhance Capture.
func TestCrewmanJorg(t *testing.T) {
	t.Run("steals with no Star Alliance neighbor", func(t *testing.T) {
		var jorg ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&jorg, CrewmanJorg)),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(jorg)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("does not steal beside a Star Alliance creature", func(t *testing.T) {
		var jorg ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.StarAlliance)),
					ct.Bind(&jorg, CrewmanJorg),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(jorg)

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(3)
	})
}
