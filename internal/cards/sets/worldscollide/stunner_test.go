package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Stunner
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: You may stun a creature."
func TestStunner(t *testing.T) {
	t.Run("its host may stun a creature when it reaps", func(t *testing.T) {
		var host, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&host,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						Stunner,
					),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature()))},
		})

		h.P1.Reap(host)
		h.P1.ClickOption("Yes")
		h.P1.ClickCard(foe)

		h.Expect(foe).Stunned(true)
	})
}
