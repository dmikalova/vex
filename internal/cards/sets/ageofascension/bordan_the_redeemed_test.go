package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bordan the Redeemed
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Thief
//
//	Elusive.
//	Action: Bordan the Redeemed captures 2 Æmber from your opponent.
func TestBordanTheRedeemed(t *testing.T) {
	t.Run("captures 2 aember from the opponent", func(t *testing.T) {
		var bordan ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&bordan, BordanTheRedeemed)),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(bordan)

		h.Expect(bordan).AmberOn(2)
		h.P2.ExpectAmber(1)
	})
}
