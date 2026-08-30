package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Terms of Redress
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: A friendly creature captures 2 Æmber from your opponent.
func TestTermsOfRedress(t *testing.T) {
	t.Run("a friendly creature captures 2 Æmber from the opponent", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				Hand:   ct.Cards(TermsOfRedress),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum)))),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(TermsOfRedress)

		h.Expect(ally).AmberOn(2)
		h.P2.ExpectAmber(3)
	})
}
