package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Commander Chan
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human
//
//	Fight/Reap: Use another Creature.
func TestCommanderChan(t *testing.T) {
	t.Run("reaps and uses another friendly creature", func(t *testing.T) {
		var cmdr, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Bind(&cmdr, CommanderChan),
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(cmdr)

		h.Expect(ally).Exhausted()
		h.P1.ExpectAmber(2) // 1 from Chan's reap, 1 from the ally's reap
	})
}
