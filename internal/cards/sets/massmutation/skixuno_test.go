package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Skixuno
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Demon
//
//	Omega.
//	Play: Destroy each other creature. For each creature destroyed this way, give Skixuno a +1 power counter.
func TestSkixuno(t *testing.T) {
	t.Run("destroys each other creature and grows per kill", func(t *testing.T) {
		var enemy1, enemy2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Skixuno),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&enemy1, ct.Creature(ct.Power(5))),
				ct.Bind(&enemy2, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Play(Skixuno)

		h.Expect(enemy1).At(ct.Discard)
		h.Expect(enemy2).At(ct.Discard)
		h.Expect(Skixuno).Power(3) // 1 base + one +1 counter per creature destroyed
	})

	t.Run("gains no counters when nothing else is in play", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Skixuno),
			},
		})

		h.P1.Play(Skixuno)

		h.Expect(Skixuno).Power(1)
	})
}
