package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ballcano
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 4 damage to each Creature. Gain 2 chains.
func TestBallcano(t *testing.T) {
	t.Run("deals 4 damage to each creature and gains 2 chains", func(t *testing.T) {
		var mine, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				Hand:   ct.Cards(Ballcano),
				InPlay: ct.Cards(ct.Bind(&mine, ct.Creature(ct.Power(6)))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(6))))},
		})

		h.P1.Play(Ballcano)

		h.Expect(mine).Damage(4)
		h.Expect(foe).Damage(4)
		if got := h.Game().State.Chains[0]; got != 2 {
			t.Errorf("P1 chains = %d, want 2", got)
		}
	})
}
