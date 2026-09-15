package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Crystal Hive
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Location
//
//	Action: For the remainder of the turn, after a creature reaps, gain 1 Æmber.
func TestCrystalHive(t *testing.T) {
	t.Run("gains an extra Æmber each time a creature reaps this turn", func(t *testing.T) {
		var creature ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					CrystalHive,
					ct.Bind(&creature, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.UseAction(CrystalHive)
		h.P1.ExpectAmber(0)

		h.P1.Reap(creature) // 1 for the reap + 1 from Crystal Hive
		h.P1.ExpectAmber(2)
	})
}
