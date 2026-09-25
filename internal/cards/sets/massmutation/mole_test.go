package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Mole
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Your opponent may spend Æmber on this creature as if it were in their pool."
func TestMole(t *testing.T) {
	t.Run("lets the creature's controller's opponent spend its Æmber to forge", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Upgraded(ct.Bind(&host, ct.Creature(ct.Power(4))), Mole)),
			},
			P2: ct.Side{Amber: engine.KeyCost - 3},
		})
		g := h.Game()
		g.AddAmberOn(host.ID(), 3)

		// P2 is the opponent of the host's controller, so the host's 3 Æmber counts
		// toward P2's key at the start of their turn.
		g.EndPlayPhase(0)
		g.StartTurn(1)

		if g.Keys(1) != 1 {
			t.Errorf("opponent keys = %d, want 1 (forged using the creature's Æmber)", g.Keys(1))
		}
		if g.AmberOn(host.ID()) != 0 {
			t.Errorf("host Æmber = %d, want 0 (spent by the opponent)", g.AmberOn(host.ID()))
		}
	})
}
