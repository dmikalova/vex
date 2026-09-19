package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Siren Horn
//
//	House:  Saurian
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "Before Fight: Move 1 Æmber from this creature to the creature it fights."
func TestSirenHorn(t *testing.T) {
	t.Run("moves 1 Æmber from its host to the creature it fights", func(t *testing.T) {
		var host, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3))),
						SirenHorn,
					),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(5))))},
		})
		h.Game().State.Cards[host.ID()].Amber = 1

		h.P1.Fight(host, enemy)

		h.Expect(enemy).AmberOn(1)
	})
}
