package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Opportunist
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Common
//	Bonus:  Æmber
//
//	This creature gains elusive.
//	Play: this creature captures 1 Æmber from your opponent.
func TestOpportunist(t *testing.T) {
	t.Run("grants elusive and captures 1 Æmber from the opponent when played", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				Hand:   ct.Cards(Opportunist),
				InPlay: ct.Cards(ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Shadows)))),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(Opportunist)

		h.Expect(host).AmberOn(1)
		h.P2.ExpectAmber(2)
		if !h.Game().HasKeyword(host.ID(), engine.Elusive) {
			t.Error("host should gain elusive from Opportunist")
		}
	})
}
