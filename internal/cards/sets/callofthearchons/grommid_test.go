package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Grommid
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  10
//	Traits: Beast
//
//	You cannot play creatures.
//	After a creature is destroyed fighting Grommid, your opponent loses 1 Æmber.
func TestGrommid(t *testing.T) {
	setup := func(t *testing.T) (h *ct.Harness, grommid, enemy ct.Card) {
		h = ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ct.Bind(&grommid, Grommid)),
				Hand:   ct.Cards(ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
				Amber: 3,
			},
		})
		return h, grommid, enemy
	}

	t.Run("stops its controller from playing creatures", func(t *testing.T) {
		h, _, _ := setup(t)

		if _, err := h.Game().PlayCreature(0, 0, false); err != engine.ErrCannotPlayCreature {
			t.Errorf("PlayCreature = %v, want ErrCannotPlayCreature", err)
		}
	})

	t.Run("drains 1 Æmber when an enemy is destroyed fighting it", func(t *testing.T) {
		h, grommid, enemy := setup(t)

		h.P1.Fight(grommid, enemy)

		h.Expect(enemy).At(ct.Discard)
		h.P2.ExpectAmber(2) // lost 1
	})
}
