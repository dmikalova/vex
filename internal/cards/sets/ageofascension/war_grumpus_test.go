package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// War Grumpus
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Beast
//
//	Fight/Reap: Ready and fight with a neighboring Giant trait creature.
func TestWarGrumpus(t *testing.T) {
	t.Run("reap readies and fights with a neighboring Giant", func(t *testing.T) {
		var giant, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Bind(&giant, ct.Creature(ct.Traits(card.Traits.Giant), ct.Power(6))),
					WarGrumpus,
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20)))),
			},
		})

		giant.Exhaust()

		h.P1.Reap(WarGrumpus)

		h.Expect(foe).Damage(6)
	})
}
