package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Zorg
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  7
//	Traits: Beast
//
//	Zorg enters play stunned.
//	Before Fight: Stun the creature Zorg fights and each of its neighbors.
func TestZorg(t *testing.T) {
	t.Run("enters play stunned", func(t *testing.T) {
		var zorg ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(ct.Bind(&zorg, Zorg)),
			},
		})

		h.P1.Play(Zorg)
		h.Expect(zorg).Stunned(true)
	})

	t.Run("stuns the creature it fights and that creature's neighbors", func(t *testing.T) {
		var zorg, left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ct.Bind(&zorg, Zorg)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(9))),
					ct.Bind(&mid, ct.Creature(ct.Power(9))),
					ct.Bind(&right, ct.Creature(ct.Power(9))),
				),
			},
		})
		zorg.Ready()

		h.P1.Fight(zorg, mid)

		h.Expect(mid).Stunned(true)
		h.Expect(left).Stunned(true)
		h.Expect(right).Stunned(true)
	})
}
