package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Kangaphant
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Beast
//
//	Each Creature gains, "Reap: Destroy this Creature."
func TestKangaphant(t *testing.T) {
	t.Run("each creature gains Reap: Destroy this creature", func(t *testing.T) {
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					Kangaphant,
					ct.Bind(&reaper, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(reaper)

		h.Expect(reaper).At(ct.Discard)
	})
}
