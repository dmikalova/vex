package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mugwump
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Giant
//
//	After a creature is destroyed fighting Mugwump, fully heal Mugwump, and give Mugwump a +1 power counter.
func TestMugwump(t *testing.T) {
	t.Run("fully heals and gains a power counter after a kill in combat", func(t *testing.T) {
		var mugwump, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(ct.Bind(&mugwump, Mugwump))},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
			},
		})

		mugwump.Damaged(3)
		h.P1.Fight(mugwump, foe)

		h.Expect(foe).At(ct.Discard)
		h.Expect(mugwump).Damage(0).Power(7) // healed + a +1 power counter
	})
}
