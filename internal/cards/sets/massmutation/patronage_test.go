package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Patronage
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Move half the Æmber from a creature to your pool, rounding up. Move all Æmber from the chosen creature to your opponent's pool.
func TestPatronage(t *testing.T) {
	t.Run("splits an odd pool, rounding your share up", func(t *testing.T) {
		var patronage, laden ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(ct.Bind(&patronage, Patronage)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&laden, ct.Creature()))},
		})
		h.Game().AddAmberOn(laden.ID(), 3)

		h.P1.Play(patronage) // sole creature auto-selected

		h.Expect(laden).AmberOn(0)
		h.P1.ExpectAmber(2) // half of 3, rounded up
		h.P2.ExpectAmber(1) // the remainder
	})

	t.Run("splits an even pool evenly", func(t *testing.T) {
		var patronage, laden ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(ct.Bind(&patronage, Patronage)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&laden, ct.Creature()))},
		})
		h.Game().AddAmberOn(laden.ID(), 4)

		h.P1.Play(patronage)

		h.Expect(laden).AmberOn(0)
		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(2)
	})
}
