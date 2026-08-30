package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Nocturnal Maneuver
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Exhaust up to 3 creatures.
func TestNocturnalManeuver(t *testing.T) {
	t.Run("exhausts up to 3 creatures", func(t *testing.T) {
		var a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Hand:   ct.Cards(NocturnalManeuver),
				InPlay: ct.Cards(ct.Bind(&a, ct.Creature()), ct.Bind(&b, ct.Creature())),
			},
		})

		h.P1.Play(NocturnalManeuver)
		h.P1.ClickOption(a.Name())
		h.P1.ClickOption(b.Name())

		h.Expect(a).Exhausted()
		h.Expect(b).Exhausted()
	})
}
