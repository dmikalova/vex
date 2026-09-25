package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Tachyon Pulse
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Destroy each artifact. Exhaust each creature with an upgrade.
func TestTachyonPulse(t *testing.T) {
	t.Run("destroys each artifact and exhausts each upgraded creature", func(t *testing.T) {
		var artifact, upgraded, plain ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(TachyonPulse),
				InPlay: ct.Cards(
					ct.Bind(&artifact, ct.Artifact()),
					ct.Upgraded(
						ct.Bind(&upgraded, ct.Creature(ct.Power(4))),
						ct.Upgrade(),
					),
					ct.Bind(&plain, ct.Creature(ct.Power(4))),
				),
			},
		})

		h.P1.Play(TachyonPulse)

		h.Expect(artifact).At(ct.Discard)
		h.Expect(upgraded).Exhausted()
		h.Expect(plain).Ready()
	})
}
