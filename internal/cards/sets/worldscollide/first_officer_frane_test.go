package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// First Officer Frane
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human
//
//	Play/Fight/Reap: A friendly creature captures 1 Æmber from your opponent.
func TestFirstOfficerFrane(t *testing.T) {
	t.Run("captures 1 Æmber from the opponent when played", func(t *testing.T) {
		var frane ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(ct.Bind(&frane, FirstOfficerFrane)),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(frane)

		h.Expect(frane).AmberOn(1)
		h.P2.ExpectAmber(2)
	})
}
