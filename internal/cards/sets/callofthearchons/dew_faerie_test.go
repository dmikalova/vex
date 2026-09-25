package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Dew Faerie
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Faerie
//
//	Elusive.
//	Reap: Gain 1 Æmber.
func TestDewFaerie(t *testing.T) {
	t.Run("reaps for the usual Æmber plus 1 more", func(t *testing.T) {
		var faerie ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&faerie, DewFaerie)),
			},
		})

		h.Expect(faerie).Power(2)
		h.P1.Reap(faerie)
		h.P1.ExpectAmber(2) // 1 reap + 1 ability
	})
}
