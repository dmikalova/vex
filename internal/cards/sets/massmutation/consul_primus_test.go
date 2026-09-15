package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Consul Primus
func TestConsulPrimus(t *testing.T) {
	t.Run("reaping moves 1 Æmber from one creature onto another", func(t *testing.T) {
		var from, onto ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Saurian, InPlay: ct.Cards(ConsulPrimus)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&from, ct.Creature(ct.Power(3))),
				ct.Bind(&onto, ct.Creature(ct.Power(3))),
			)},
		})
		h.Game().State.Cards[from.ID()].Amber = 1

		h.P1.Reap(ConsulPrimus)
		h.P1.ClickCard(onto)

		h.Expect(from).AmberOn(0)
		h.Expect(onto).AmberOn(1)
	})
}
