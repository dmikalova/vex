package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Imperial Scutum
//
//	House:  Saurian
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This creature gains +2 armor.
//	This creature gains, "Destroyed: Move each Æmber on this creature to the common supply."
func TestImperialScutum(t *testing.T) {
	t.Run("grants +2 armor", func(t *testing.T) {
		var host ct.Card
		ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Upgraded(ct.Bind(&host, ct.Creature()), ImperialScutum),
				),
			},
		})

		if got := host.Armor(); got != 2 {
			t.Errorf("armor = %d, want 2", got)
		}
	})

	t.Run("destroyed moves the host's Æmber to the common supply", func(t *testing.T) {
		var host, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Saurian))),
						ImperialScutum,
					),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6))))},
		})
		h.Game().State.Cards[host.ID()].Amber = 2

		h.P1.Fight(host, enemy)

		h.Expect(host).At(ct.Discard)
		h.P2.ExpectAmber(0) // Æmber went to the common supply, not the opponent
	})
}
