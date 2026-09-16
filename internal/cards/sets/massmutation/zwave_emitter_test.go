package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Z-Wave Emitter
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Connected
//	Bonus:  Æmber
//
//	This creature gains, "At the start of your turn, ward this creature."
func TestZWaveEmitter(t *testing.T) {
	t.Run("wards its host at the start of its controller's turn", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
						ZWaveEmitter,
					),
				),
			},
		})

		h.P1.EndTurn() // to the opponent
		h.P2.EndTurn() // back to P1: the ward fires at the start of the turn

		if !h.Game().Warded(host.ID()) {
			t.Error("the host should be warded at the start of its controller's turn")
		}
	})
}
