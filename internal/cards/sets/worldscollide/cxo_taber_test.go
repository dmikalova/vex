package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// CXO Taber
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Alien • Krxix
//
//	Fight/Reap: You may play or use one non-Star Alliance card this turn.
func TestCXOTaber(t *testing.T) {
	t.Run("reaping frees one off-house play this turn", func(t *testing.T) {
		var taber ct.Card
		marsCreature := ct.Creature(ct.OfHouse(card.House.Mars))
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&taber, CXOTaber)),
				Hand:   ct.Cards(marsCreature),
			},
		})

		h.P1.ExpectCannotPlay(marsCreature)

		h.P1.Reap(taber)

		h.P1.Play(marsCreature)
	})

	t.Run("reaping frees one off-house use this turn", func(t *testing.T) {
		var taber ct.Card
		marsCreature := ct.Creature(ct.OfHouse(card.House.Mars))
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&taber, CXOTaber), marsCreature),
			},
			P2: ct.Side{Amber: 2},
		})

		h.P1.ExpectCannotUse(marsCreature)

		h.P1.Reap(taber)

		// The freed Mars creature may now reap.
		h.P1.Reap(marsCreature)
		h.P1.ExpectAmber(2)
	})
}
