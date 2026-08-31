package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

func TestSigilOfBrotherhood(t *testing.T) {
	t.Run("destroys itself and grants use of friendly Sanctum creatures", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(SigilOfBrotherhood)},
		})

		h.P1.UseAction(SigilOfBrotherhood)

		h.Expect(SigilOfBrotherhood).At(ct.Discard)
		if h.Game().State.MayUseHouse[0] != engine.Sanctum {
			t.Error("Sigil should grant use of friendly Sanctum creatures this turn")
		}
	})
}

func TestRitualOfTheHunt(t *testing.T) {
	t.Run("destroys itself and grants use of friendly Untamed creatures", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(RitualOfTheHunt)},
		})

		h.P1.UseAction(RitualOfTheHunt)

		h.Expect(RitualOfTheHunt).At(ct.Discard)
		if h.Game().State.MayUseHouse[0] != engine.Untamed {
			t.Error("Ritual should grant use of friendly Untamed creatures this turn")
		}
	})
}
