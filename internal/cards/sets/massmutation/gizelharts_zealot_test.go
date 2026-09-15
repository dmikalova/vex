package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gizelhart's Zealot
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human • Knight
//
//	Gizelhart's Zealot enters play ready and enrage Gizelhart's Zealot.
func TestGizelhartsZealot(t *testing.T) {
	t.Run("enters play ready and enraged", func(t *testing.T) {
		var zealot ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(ct.Bind(&zealot, GizelhartsZealot)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Creature())},
		})

		h.P1.Play(zealot)

		h.Expect(zealot).Ready()
		if !h.Game().Enraged(zealot.ID()) {
			t.Errorf("%s should be enraged", zealot.Name())
		}
	})
}
