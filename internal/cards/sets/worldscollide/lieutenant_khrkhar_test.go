package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Lieutenant Khrkhar
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Alien • Handuhan
//
//	Taunt, Hazardous 3.
func TestLieutenantKhrkhar(t *testing.T) {
	t.Run("deals 3 hazardous damage to an attacker before combat", func(t *testing.T) {
		var attacker, khrkhar ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&attacker, ct.Creature(ct.Power(20)))),
			},
			P2: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&khrkhar, LieutenantKhrkhar)),
			},
		})
		attacker.Ready()

		h.P1.Fight(attacker, khrkhar)

		h.Expect(khrkhar).At(ct.Discard)
		// 5 fight damage from Khrkhar's power plus 3 from Hazardous.
		h.Expect(attacker).At(ct.PlayArea).Damage(8)
	})
}
