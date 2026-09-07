package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// The Shadow Council
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Elf • Leader • Thief
//
//	Elusive.
//	While The Shadow Council is in the center of your battleline, it gains, "Action: Steal 2 Æmber."
func TestTheShadowCouncil(t *testing.T) {
	t.Run("centered, may use the granted action to steal 2 Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(
				ct.Creature(),
				TheShadowCouncil,
				ct.Creature(),
			)},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(TheShadowCouncil)

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(1)
	})

	t.Run("off-center, the granted action is unavailable", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(
				TheShadowCouncil,
				ct.Creature(),
			)},
			P2: ct.Side{Amber: 3},
		})

		council := councilID(t, h)
		if h.Game().HasTrigger(council, engine.TriggerAction) {
			t.Error("off-center The Shadow Council should not have the granted action")
		}
	})
}

// councilID finds The Shadow Council's id on player 1's battleline.
func councilID(t *testing.T, h *ct.Harness) engine.LocalID {
	t.Helper()
	for _, id := range h.Game().Battleline(0) {
		if h.Game().Name(id) == "The Shadow Council" {
			return id
		}
	}
	t.Fatal("The Shadow Council not found in battleline")
	return 0
}
